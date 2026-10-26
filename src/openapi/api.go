package openapi

import (
	"context"
	notebook "github.com/Lixsp11/inspire-go-sdk/src/notebookv2/client"
	sdk "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	project "github.com/Lixsp11/inspire-go-sdk/src/projectv2/client"
	train "github.com/Lixsp11/inspire-go-sdk/src/trainv2/client"
	user "github.com/Lixsp11/inspire-go-sdk/src/userv2/client"
	workspace "github.com/Lixsp11/inspire-go-sdk/src/workspacev2/client"
	"github.com/alibabacloud-go/tea/dara"
)

func apiRegistry() []Api {
	return []Api{
		{Product: "openapi", Name: "GetDiscovery", NewRequest: func() requestModel { return &emptyRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := client.GetDiscoveryWithContext(ctx, runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "train", Name: "ListJobs", NewRequest: func() requestModel { return &train.ListJobsRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&train.Client{Client: client}).ListJobsWithContext(ctx, request.(*train.ListJobsRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "train", Name: "GetJob", NewRequest: func() requestModel { return &train.GetJobRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&train.Client{Client: client}).GetJobWithContext(ctx, request.(*train.GetJobRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "train", Name: "CreateJobConsole", NewRequest: func() requestModel { return &train.CreateJobConsoleRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&train.Client{Client: client}).CreateJobConsoleWithContext(ctx, request.(*train.CreateJobConsoleRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "train", Name: "StopJob", NewRequest: func() requestModel { return &train.StopJobRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&train.Client{Client: client}).StopJobWithContext(ctx, request.(*train.StopJobRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "notebook", Name: "ListNotebooks", NewRequest: func() requestModel { return &notebook.ListNotebooksRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&notebook.Client{Client: client}).ListNotebooksWithContext(ctx, request.(*notebook.ListNotebooksRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "notebook", Name: "GetNotebook", NewRequest: func() requestModel { return &notebook.GetNotebookRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&notebook.Client{Client: client}).GetNotebookWithContext(ctx, request.(*notebook.GetNotebookRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "project", Name: "ListProjects", NewRequest: func() requestModel { return &project.ListProjectsRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&project.Client{Client: client}).ListProjectsWithContext(ctx, request.(*project.ListProjectsRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "workspace", Name: "ListLogicComputeGroups", NewRequest: func() requestModel { return &workspace.ListLogicComputeGroupsRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&workspace.Client{Client: client}).ListLogicComputeGroupsWithContext(ctx, request.(*workspace.ListLogicComputeGroupsRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
		{Product: "user", Name: "GetUserDetail", NewRequest: func() requestModel { return &user.GetUserDetailRequest{} }, Call: func(ctx context.Context, client *sdk.Client, request requestModel, runtime *dara.RuntimeOptions) (interface{}, error) {
			response, err := (&user.Client{Client: client}).GetUserDetailWithContext(ctx, request.(*user.GetUserDetailRequest), runtime)
			if err != nil {
				return nil, err
			}
			return response.Body, nil
		}},
	}
}
