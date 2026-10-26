package client

import (
	"context"
	"encoding/json"
	openapi "github.com/Lixsp11/inspire-go-sdk/src/openapi/client"
	"github.com/Lixsp11/inspire-go-sdk/src/openapi/model"
	"github.com/alibabacloud-go/tea/dara"
	"math"
)

// ListJobsPaginator is sequential and not safe for concurrent NextPage calls.
// Each returned response retains the actual platform total, page request and diagnostics.
type ListJobsPaginator struct {
	client                 *Client
	request                *ListJobsRequest
	runtime                *dara.RuntimeOptions
	next, last, size       int32
	pages, maxPages        int
	consumed, total        int64
	hasTotal, done, failed bool
	signatures             map[string]bool
	ids                    map[string]bool
}

func NewListJobsPaginator(client *Client, request *ListJobsRequest, runtime *dara.RuntimeOptions) (*ListJobsPaginator, error) {
	if client == nil || client.Client == nil || request == nil {
		return nil, openapi.ValidationError("Client and ListJobsRequest are required.")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, openapi.ValidationError("Request is not finite JSON.")
	}
	query := &ListJobsRequest{}
	if err = json.Unmarshal(data, query); err != nil {
		return nil, err
	}
	first, size := int32(1), int32(50)
	if query.PageNum != nil {
		first = *query.PageNum
	}
	if query.PageSize != nil {
		size = *query.PageSize
	}
	return &ListJobsPaginator{client: client, request: query, runtime: runtime, next: first, size: size, maxPages: 100, consumed: int64(first-1) * int64(size), signatures: map[string]bool{}, ids: map[string]bool{}}, nil
}
func (p *ListJobsPaginator) SetMaxPages(value int) error {
	if value < 1 || p.pages != 0 {
		return openapi.ValidationError("Set positive max pages before the first page.")
	}
	p.maxPages = value
	return nil
}
func (p *ListJobsPaginator) HasMorePages() bool           { return p != nil && !p.done && !p.failed }
func (p *ListJobsPaginator) LastPageNum() int32           { return p.last }
func (p *ListJobsPaginator) PagesRead() int               { return p.pages }
func (p *ListJobsPaginator) ReportedTotal() (int64, bool) { return p.total, p.hasTotal }
func (p *ListJobsPaginator) NextPage(ctx context.Context) (result *ListJobsResponse, err error) {
	if p == nil || !p.HasMorePages() {
		return nil, openapi.NewError("PaginationError", "NoMorePages", "No unread page remains.")
	}
	defer func() {
		if err != nil {
			p.failed = true
		}
	}()
	if p.pages >= p.maxPages {
		return nil, openapi.NewError("PaginationError", "PageLimitReached", "Pagination reached its explicit page limit.")
	}
	p.request.PageNum, p.request.PageSize = dara.Int32(p.next), dara.Int32(p.size)
	result, err = p.client.ListJobsWithContext(ctx, p.request, p.runtime)
	if err != nil {
		return nil, err
	}
	value := result.Body.Result
	total, present, err := model.Total(value.Total)
	if err != nil {
		return nil, openapi.NewError("InvalidResponseError", "InvalidTotal", "Pagination total is invalid.")
	}
	if present && p.hasTotal && total != p.total {
		return nil, openapi.NewError("PaginationError", "ChangedTotal", "Reported pagination total changed.")
	}
	if present {
		p.total, p.hasTotal = total, true
	}
	raw, err := json.Marshal(value.Jobs)
	if err != nil {
		return nil, openapi.NewError("InvalidResponseError", "InvalidCollection", "Page cannot be encoded.")
	}
	signature := string(raw)
	if len(value.Jobs) > 0 && p.signatures[signature] {
		return nil, openapi.NewError("PaginationError", "RepeatedPage", "Platform returned a repeated page.")
	}
	p.signatures[signature] = true
	if len(value.Jobs) == 0 && p.hasTotal && p.consumed < p.total {
		return nil, openapi.NewError("PaginationError", "TruncatedPages", "Pagination ended before the reported total.")
	}
	for _, job := range value.Jobs {
		if job == nil {
			return nil, openapi.NewError("InvalidResponseError", "InvalidCollection", "Page contains no valid job.")
		}
		if job.JobId != nil {
			if p.ids[*job.JobId] {
				return nil, openapi.NewError("PaginationError", "DuplicateItem", "Duplicate job across pages.")
			}
			p.ids[*job.JobId] = true
		}
	}
	p.consumed += int64(len(value.Jobs))
	p.pages++
	p.last = p.next
	if p.hasTotal {
		p.done = p.consumed >= p.total
	} else {
		p.done = len(value.Jobs) < int(p.size)
	}
	if !p.done && p.next == math.MaxInt32 {
		return nil, openapi.NewError("PaginationError", "PageLimitReached", "Page number cannot be increased.")
	}
	p.next++
	return result, nil
}
