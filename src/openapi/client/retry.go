package client

import (
	"context"
	"fmt"
	"github.com/alibabacloud-go/tea/dara"
	"math"
	"net/http"
	"strconv"
	"time"
)

func cloneRetry(value *dara.RetryOptions) *dara.RetryOptions {
	if value == nil {
		return nil
	}
	out := *value
	copyConditions := func(values []*dara.RetryCondition) []*dara.RetryCondition {
		result := make([]*dara.RetryCondition, len(values))
		for i, v := range values {
			if v == nil {
				continue
			}
			c := *v
			c.Exception = append([]string(nil), v.Exception...)
			c.ErrorCode = append([]string(nil), v.ErrorCode...)
			switch b := v.Backoff.(type) {
			case *dara.FixedBackoffPolicy:
				if b != nil {
					x := *b
					c.Backoff = &x
				}
			case *dara.RandomBackoffPolicy:
				if b != nil {
					x := *b
					c.Backoff = &x
				}
			case *dara.ExponentialBackoffPolicy:
				if b != nil {
					x := *b
					c.Backoff = &x
				}
			case *dara.EqualJitterBackoffPolicy:
				if b != nil {
					x := *b
					c.Backoff = &x
				}
			case *dara.FullJitterBackoffPolicy:
				if b != nil {
					x := *b
					c.Backoff = &x
				}
			}
			result[i] = &c
		}
		return result
	}
	out.RetryCondition = copyConditions(value.RetryCondition)
	out.NoRetryCondition = copyConditions(value.NoRetryCondition)
	return &out
}
func resolveRetry(config *dara.RetryOptions, r *dara.RuntimeOptions) (*dara.RetryOptions, error) {
	legacy := r.Autoretry != nil || r.MaxAttempts != nil || r.BackoffPeriod != nil || r.BackoffPolicy != nil
	if config != nil && legacy {
		return nil, configError("Choose Config.RetryOptions or legacy RuntimeOptions retry fields.")
	}
	value := cloneRetry(config)
	if config == nil && legacy {
		if r.Autoretry == nil {
			return nil, configError("Legacy retry fields require explicit Autoretry.")
		}
		if !*r.Autoretry {
			return nil, nil
		}
		n, period := 3, 100
		if r.MaxAttempts != nil {
			n = *r.MaxAttempts
		}
		if r.BackoffPeriod != nil {
			period = *r.BackoffPeriod
		}
		policy := dara.StringValue(r.BackoffPolicy)
		if policy == "no" {
			period = 0
		}
		if n < 1 || period < 0 || (policy != "" && policy != "fix" && policy != "Fixed" && policy != "no") {
			return nil, configError("Invalid legacy retry limit or fixed backoff.")
		}
		value = &dara.RetryOptions{Retryable: true, MaxAttempts: n, RetryCondition: []*dara.RetryCondition{{MaxAttempts: n, Exception: []string{"TransportError", "ServerError", "ThrottlingError", "AuthenticationError"}, MaxDelay: 120000, Backoff: &dara.FixedBackoffPolicy{Period: period}}}}
	}
	if value == nil {
		return nil, nil
	}
	if value.MaxAttempts < 0 {
		return nil, configError("Retry max attempts cannot be negative.")
	}
	for _, conditions := range [][]*dara.RetryCondition{value.RetryCondition, value.NoRetryCondition} {
		for _, c := range conditions {
			if c == nil || c.MaxAttempts < 0 || c.MaxDelay < 0 {
				return nil, configError("Invalid retry condition.")
			}
			switch b := c.Backoff.(type) {
			case *dara.FixedBackoffPolicy:
				if b == nil || b.Period < 0 {
					return nil, configError("Invalid fixed backoff.")
				}
			case *dara.RandomBackoffPolicy:
				if b == nil || b.Period <= 0 || b.Cap < 1 {
					return nil, configError("Invalid random backoff.")
				}
			case *dara.ExponentialBackoffPolicy:
				if b == nil || b.Period < 1 || b.Cap < 1 {
					return nil, configError("Invalid exponential backoff.")
				}
			case *dara.EqualJitterBackoffPolicy:
				if b == nil || b.Period < 1 || b.Cap < 1 {
					return nil, configError("Invalid jitter backoff.")
				}
			case *dara.FullJitterBackoffPolicy:
				if b == nil || b.Period < 1 || b.Cap < 1 {
					return nil, configError("Invalid jitter backoff.")
				}
			}
		}
	}
	return value, nil
}
func retryDecision(options *dara.RetryOptions, err *SDKError, attempt int) (retry bool, delay time.Duration, failure error) {
	defer func() {
		if recover() != nil {
			retry = false
			failure = configError("Retry policy failed to evaluate.")
		}
	}()
	ctx := &dara.RetryPolicyContext{RetriesAttempted: attempt, Exception: err}
	if !dara.ShouldRetry(options, ctx) {
		return false, 0, nil
	}
	ms := dara.GetBackoffDelay(options, ctx)
	if ms < 0 || int64(ms) > math.MaxInt64/int64(time.Millisecond) {
		return false, 0, configError("Retry delay is outside the supported range.")
	}
	return true, time.Duration(ms) * time.Millisecond, nil
}
func retryAfter(value string) time.Duration {
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		if seconds <= 0 {
			return 0
		}
		if seconds > math.MaxInt64/int64(time.Second) {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(seconds) * time.Second
	}
	if date, err := http.ParseTime(value); err == nil {
		if delay := time.Until(date); delay > 0 {
			return delay
		}
	}
	return 0
}
func wait(ctx context.Context, delay time.Duration) error {
	if end, ok := ctx.Deadline(); ok && delay >= time.Until(end) {
		return context.DeadlineExceeded
	}
	if delay == 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func timeError(ctx context.Context) *SDKError {
	err := NewError("OperationTimeoutError", "OperationTimeout", "Operation deadline expired.")
	if ctx.Err() == context.Canceled {
		err = NewError("CanceledError", "Canceled", "Operation was canceled.")
	}
	err.cause = ctx.Err()
	if err.cause == nil {
		err.cause = context.DeadlineExceeded
	}
	return err
}

// NewContextError retains errors.Is semantics while providing native SDK diagnostics.
func NewContextError(ctx context.Context) *SDKError {
	if ctx == nil || ctx.Err() == nil {
		return configError("Context has not finished.")
	}
	return timeError(ctx)
}
func policyError(err error) *SDKError {
	if sdk, ok := err.(*SDKError); ok {
		return sdk
	}
	return configError(fmt.Sprint("Invalid retry policy."))
}
