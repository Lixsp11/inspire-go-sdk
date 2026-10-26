package openapi

import "github.com/Lixsp11/inspire-go-sdk/src/cli"

const (
	BodyFlagName          = "body"
	BodyFileFlagName      = "body-file"
	CliDryRunJsonFlagName = "cli-dry-run-json"
	CliOutputFlagName     = "cli-output"
)

func AddFlags(fs *cli.FlagSet) {
	fs.Add(NewBodyFlag())
	fs.Add(NewBodyFileFlag())
	fs.Add(NewCliDryRunJsonFlag())
	fs.Add(&cli.Flag{Name: "max-pages", Short: "Page limit for --pager (default 100).", AssignedMode: cli.AssignedOnce, Category: "helper"})
}
func NewBodyFlag() *cli.Flag {
	return &cli.Flag{Name: BodyFlagName, Short: "Request JSON (SDK or native field names).", AssignedMode: cli.AssignedOnce, Category: "caller"}
}
func NewBodyFileFlag() *cli.Flag {
	return &cli.Flag{Name: BodyFileFlagName, Short: "Request JSON file; - reads stdin.", AssignedMode: cli.AssignedOnce, Category: "caller"}
}
func NewCliDryRunJsonFlag() *cli.Flag {
	return &cli.Flag{Name: CliDryRunJsonFlagName, Short: "Validate request and explicit caller options without platform access.", AssignedMode: cli.AssignedNone, Category: "helper"}
}
func NewCliOutputFlag() *cli.Flag {
	return &cli.Flag{Name: CliOutputFlagName, Short: "Help and diagnostics format: json.", AssignedMode: cli.AssignedOnce, Persistent: true, Category: "help"}
}
