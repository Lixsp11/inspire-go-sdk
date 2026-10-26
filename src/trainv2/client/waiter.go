package client

import (
	"context"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/alibabacloud-go/tea/dara"
	"time"
)

// WaiterOptions uses native Go time.Duration; Config and Dara RuntimeOptions use milliseconds.
type WaiterOptions struct {
	Timeout, Interval           time.Duration
	TargetStates, FailureStates []string
}

func (c *Client) WaitForJob(parent context.Context, jobId string, options *WaiterOptions) (*GetJobResponse, error) {
	if parent == nil || jobId == "" {
		return nil, openapi.ValidationError("Context and job ID are required.")
	}
	timeout, interval := 10*time.Minute, 5*time.Second
	targets, failures := []string{"SUCCEEDED", "job_succeeded"}, []string{"FAILED", "CANCELLED", "STOPPED", "job_failed", "job_cancelled", "job_stopped"}
	if options != nil {
		if options.Timeout < 0 || options.Interval < 0 {
			return nil, openapi.ValidationError("Waiter durations cannot be negative.")
		}
		if options.Timeout > 0 {
			timeout = options.Timeout
		}
		if options.Interval > 0 {
			interval = options.Interval
		}
		if options.TargetStates != nil {
			targets = options.TargetStates
		}
		if options.FailureStates != nil {
			failures = options.FailureStates
		}
	}
	target, failure := map[string]bool{}, map[string]bool{}
	for _, state := range targets {
		if state == "" {
			return nil, openapi.ValidationError("Waiter states must be nonempty.")
		}
		target[state] = true
	}
	if len(target) == 0 {
		return nil, openapi.ValidationError("Waiter needs a target state.")
	}
	for _, state := range failures {
		if state == "" || target[state] {
			return nil, openapi.ValidationError("Waiter states must be nonempty and disjoint.")
		}
		failure[state] = true
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	for {
		result, err := c.GetJobWithContext(ctx, &GetJobRequest{JobId: dara.String(jobId)}, &dara.RuntimeOptions{})
		if err != nil {
			if parent.Err() != nil {
				return nil, err
			}
			if ctx.Err() != nil {
				return nil, openapi.NewError("WaiterError", "WaiterTimeout", "Waiting timed out; the job was not stopped.")
			}
			return nil, err
		}
		state := dara.StringValue(result.Body.Result.Status)
		if state == "" {
			return nil, openapi.NewError("InvalidResponseError", "MissingJobStatus", "Job response omitted its status.")
		}
		if target[state] {
			return result, nil
		}
		if failure[state] {
			return nil, openapi.NewError("WaiterError", "JobFailed", "Job entered a failure terminal state.")
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			if parent.Err() != nil {
				return nil, openapi.NewContextError(parent)
			}
			return nil, openapi.NewError("WaiterError", "WaiterTimeout", "Waiting timed out; the job was not stopped.")
		case <-timer.C:
		}
	}
}
