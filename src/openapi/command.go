// Package openapi adapts typed SDK operations to command-line commands.
package openapi

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/Lixsp11/inspire-go-sdk/src/cli"
	"github.com/Lixsp11/inspire-go-sdk/src/config"
	sdk "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	train "github.com/Lixsp11/inspire-go-sdk/src/trainv2/client"
	"github.com/alibabacloud-go/tea/dara"
)

type Api struct {
	Name, Product string
	NewRequest    func() requestModel
	Call          func(context.Context, *sdk.Client, requestModel, *dara.RuntimeOptions) (interface{}, error)
}

func NewCommand() *cli.Command {
	command := &cli.Command{Name: "inspire", Short: "Inspire platform command-line interface.", Usage: "inspire <product> <Action> [--Parameter value] [options]"}
	config.AddFlags(command.Flags())
	command.Flags().Add(&cli.Flag{Name: "timeout", Short: "Total command budget in seconds (default 600).", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "caller"})
	command.Flags().Add(NewCliOutputFlag())

	command.AddSubCommand(&cli.Command{Name: "version", Short: "Print the CLI version.", Usage: "inspire version", Run: func(ctx *cli.Context, _ []string) error {
		return cli.WriteJSON(ctx.Stdout(), map[string]string{"version": sdk.Version})
	}})
	command.AddSubCommand(newConfigureCommand())
	products := map[string]*cli.Command{}
	for _, api := range apiRegistry() {
		product := products[api.Product]
		if product == nil {
			product = &cli.Command{Name: api.Product, Short: "Native " + api.Product + " operations.", Usage: "inspire " + api.Product + " <Action> [options]"}
			command.AddSubCommand(product)
			products[api.Product] = product
		}
		product.AddSubCommand(newApiCommand(api))
	}
	return command
}

func requestSchema(t reflect.Type) map[string]interface{} {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	result := map[string]interface{}{}
	switch t.Kind() {
	case reflect.Struct:
		properties := map[string]interface{}{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := strings.Split(field.Tag.Get("json"), ",")[0]
			if field.PkgPath != "" || name == "" || name == "-" {
				continue
			}
			property := requestSchema(field.Type)
			property["cli_name"] = field.Name
			properties[name] = property
		}
		result["type"] = "object"
		result["properties"] = properties
		result["additionalProperties"] = false
	case reflect.Slice:
		result["type"] = "array"
		result["items"] = requestSchema(t.Elem())
	case reflect.Map:
		result["type"] = "object"
		result["additionalProperties"] = requestSchema(t.Elem())
	case reflect.String:
		result["type"] = "string"
	case reflect.Bool:
		result["type"] = "boolean"
	case reflect.Int, reflect.Int32, reflect.Int64:
		result["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		result["type"] = "number"
	}
	return result
}

func newApiCommand(api Api) *cli.Command {
	request := api.NewRequest()
	schema := requestSchema(reflect.TypeOf(request))
	required := map[string][]string{"ListJobs": {"workspace_id"}, "GetJob": {"job_id"}, "StopJob": {"job_id"}, "GetNotebook": {"notebook_id"}, "ListNotebooks": {"workspace_id", "filter_by"}, "ListLogicComputeGroups": {"filter", "page_size"}, "CreateJobConsole": {"name", "workspace_id", "project_id", "logic_compute_group_id", "command", "framework", "task_priority", "framework_config"}}
	if fields := required[api.Name]; len(fields) > 0 {
		schema["required"] = fields
	}
	command := &cli.Command{Name: api.Name, Short: api.Name + " through the native SDK.", Usage: "inspire " + api.Product + " " + api.Name + " [--Parameter value] [options]", EnableUnknownFlag: true, RequestSchema: schema}
	t := reflect.TypeOf(request).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" || name == "" || field.PkgPath != "" {
			continue
		}
		command.Flags().Add(&cli.Flag{Name: field.Name, Short: name + "; see JSON help for nested types.", AssignedMode: cli.AssignedOnce, Category: "parameter"})
	}
	AddFlags(command.Flags())

	if api.Product == "train" && api.Name == "ListJobs" {
		command.Flags().Add(&cli.Flag{Name: "pager", Short: "Collect all jobs from page 1 using the SDK paginator.", AssignedMode: cli.AssignedNone, Category: "helper"})
	}
	if api.Product == "train" && api.Name == "GetJob" {
		command.Flags().Add(&cli.Flag{Name: "waiter", Short: "expr=Result.status to=STATE timeout=180 interval=5 (seconds).", AssignedMode: cli.AssignedRepeatable, Category: "helper"})
	}
	command.Run = func(ctx *cli.Context, _ []string) error {
		request := api.NewRequest()
		if err := loadRequest(ctx, request); err != nil {
			return err
		}
		if err := validateCallOptions(ctx, request); err != nil {
			return err
		}
		if ctx.Flags().Get("cli-dry-run-json").IsAssigned() {
			return cli.WriteJSON(ctx.Stdout(), map[string]interface{}{"Product": api.Product, "Action": api.Name, "Parameters": request, "DryRun": true})
		}
		cfg, runtime, err := loadConfig(ctx)
		if err != nil {
			return err
		}
		timeout, err := integerFlag(ctx, "timeout", 600)
		if err != nil {
			return err
		}
		cfg.OperationTimeout = dara.Int(timeout * 1000)
		client, err := sdk.NewClient(cfg)
		if err != nil {
			return err
		}
		defer client.Close()
		callContext, cancel := context.WithTimeout(ctx.Context, time.Duration(timeout)*time.Second)
		defer cancel()
		var result interface{}
		if f := ctx.Flags().Get("pager"); f.IsAssigned() {
			result, err = runPager(ctx, callContext, client, request.(*train.ListJobsRequest), runtime)
		} else if f := ctx.Flags().Get("waiter"); f.IsAssigned() {
			result, err = runWaiter(callContext, client, request.(*train.GetJobRequest), f.GetValues())
		} else {
			if _, set := ctx.Flags().GetValue("max-pages"); set {
				return cli.UsageError("InvalidParameter", "--max-pages requires --pager.")
			}
			result, err = api.Call(callContext, client, request, runtime)
		}
		if err != nil {
			return err
		}
		return cli.WriteJSON(ctx.Stdout(), result)
	}
	return command
}

func runPager(ctx *cli.Context, callContext context.Context, client *sdk.Client, request *train.ListJobsRequest, runtime *dara.RuntimeOptions) (interface{}, error) {
	if request.PageNum != nil && *request.PageNum != 1 {
		return nil, cli.UsageError("InvalidParameter", "--pager starts at PageNum 1.")
	}
	paginator, err := train.NewListJobsPaginator(&train.Client{Client: client}, request, runtime)
	if err != nil {
		return nil, err
	}
	limit, err := integerFlag(ctx, "max-pages", 100)
	if err != nil {
		return nil, err
	}
	if err = paginator.SetMaxPages(limit); err != nil {
		return nil, err
	}
	jobs := []*train.Job{}
	var body *train.ListJobsResponseBody
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(callContext)
		if err != nil {
			return nil, err
		}
		body = page.Body
		jobs = append(jobs, page.Body.Result.Jobs...)
	}
	if body == nil {
		return nil, cli.UsageError("InvalidResponse", "Paginator did not return a page.")
	}
	body.Result.Jobs = jobs
	return body, nil
}

func validateCallOptions(ctx *cli.Context, request requestModel) error {
	if _, err := integerFlag(ctx, "timeout", 600); err != nil {
		return err
	}
	for _, name := range []string{"profile", "config-path"} {
		if value, set := ctx.Flags().GetValue(name); set && value == "" {
			return cli.UsageError("InvalidParameter", "--"+name+" cannot be empty.")
		}
	}
	profile := config.Profile{}
	if err := profile.OverwriteWithFlags(ctx); err != nil {
		return err
	}
	transport := (&sdk.Config{}).SetEndpoint(profile.Endpoint).SetSocks5Proxy(profile.Socks5Proxy).SetHttpProxy(profile.HttpProxy).SetHttpsProxy(profile.HttpsProxy).SetNoProxy(profile.NoProxy)
	client, err := sdk.NewClient(transport)
	if err != nil {
		return err
	}
	client.Close()
	if f := ctx.Flags().Get("pager"); f.IsAssigned() {
		if p, ok := request.(*train.ListJobsRequest); !ok || (p.PageNum != nil && *p.PageNum != 1) {
			return cli.UsageError("InvalidParameter", "--pager starts at PageNum 1.")
		}
		if _, err := integerFlag(ctx, "max-pages", 100); err != nil {
			return err
		}
	} else if _, set := ctx.Flags().GetValue("max-pages"); set {
		return cli.UsageError("InvalidParameter", "--max-pages requires --pager.")
	}
	if f := ctx.Flags().Get("waiter"); f.IsAssigned() {
		if _, err := newWaiterOptions(f.GetValues()); err != nil {
			return err
		}
	}
	return nil
}

func newWaiterOptions(values []string) (*train.WaiterOptions, error) {
	fields := map[string]string{"timeout": "180", "interval": "5"}
	seen := map[string]bool{}
	for _, raw := range values {
		key, value, ok := strings.Cut(raw, "=")
		if !ok || value == "" || seen[key] || (key != "expr" && key != "to" && key != "timeout" && key != "interval") {
			return nil, cli.UsageError("InvalidParameter", "Invalid --waiter key=value arguments.")
		}
		seen[key] = true
		fields[key] = value
	}
	if fields["expr"] != "Result.status" || fields["to"] == "" {
		return nil, cli.UsageError("InvalidParameter", "Waiter requires expr=Result.status and to=STATE.")
	}
	duration := func(key string) (time.Duration, error) {
		raw := fields[key]
		value, err := time.ParseDuration(raw + "s")
		if err != nil || value <= 0 {
			return 0, cli.UsageError("InvalidParameter", "Waiter durations must be positive seconds.")
		}
		return value, nil
	}
	timeout, err := duration("timeout")
	if err != nil {
		return nil, err
	}
	interval, err := duration("interval")
	if err != nil {
		return nil, err
	}
	failures := []string{}
	for _, state := range []string{"FAILED", "CANCELLED", "STOPPED", "job_failed", "job_cancelled", "job_stopped"} {
		if state != fields["to"] {
			failures = append(failures, state)
		}
	}
	return &train.WaiterOptions{Timeout: timeout, Interval: interval, TargetStates: []string{fields["to"]}, FailureStates: failures}, nil
}

func runWaiter(ctx context.Context, client *sdk.Client, request *train.GetJobRequest, values []string) (interface{}, error) {
	options, err := newWaiterOptions(values)
	if err != nil {
		return nil, err
	}
	response, err := (&train.Client{Client: client}).WaitForJob(ctx, dara.StringValue(request.JobId), options)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func normalizeError(err error) error {
	var existing *cli.AgentError
	if errors.As(err, &existing) {
		return existing
	}
	var remote *sdk.SDKError
	if errors.As(err, &remote) {
		envelope := cli.AgentErrorEnvelope{Message: "Platform request failed.", ErrorCode: dara.StringValue(remote.Code), StatusCode: dara.IntValue(remote.StatusCode), RequestId: dara.StringValue(remote.RequestId), Recovery: cli.AgentErrorRecovery{Action: "inspect", Hint: "Inspect the error code and command; reconcile writes before retrying."}}
		code := 1
		if dara.StringValue(remote.Name) == "ValidationError" || dara.StringValue(remote.Name) == "ConfigurationError" {
			envelope.Message = dara.StringValue(remote.Message)
			envelope.Recovery.Action = "help"
			envelope.Recovery.Command = "inspire --help --cli-output json"
			code = 2
		}
		if dara.StringValue(remote.Name) == "CredentialError" || dara.StringValue(remote.Name) == "AuthenticationError" || dara.StringValue(remote.Name) == "PermissionDeniedError" || envelope.StatusCode == 401 || envelope.StatusCode == 403 {
			envelope.Recovery.Action = "configure"
			envelope.Recovery.Hint = "Check the existing management session and profile; stop retrying rejected credentials."
		}
		if dara.StringValue(remote.Name) == "MutationUncertainError" {
			envelope.Recovery.Action = "reconcile"
			envelope.Recovery.Hint = "Query platform state before another write; the operation may have been accepted."
		}
		return cli.NewAgentError(envelope, err).WithSchemaVersion().WithExitCode(code)
	}
	return err
}

func Run(ctx *cli.Context, args []string) int {
	err := NewCommand().Execute(ctx, args)
	if err != nil {
		err = normalizeError(err)
		_ = cli.WriteJSON(ctx.Stderr(), cli.ErrorEnvelope(err))
	}
	return cli.ExitCode(err)
}
