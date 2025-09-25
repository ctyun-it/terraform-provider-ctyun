package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type HttpHook interface {
	BeforeRequest(context.Context, *http.Request)
	AfterResponse(context.Context, *http.Response)
}

type PrintLogHttpHook struct{}

func (d PrintLogHttpHook) BeforeRequest(_ context.Context, request *http.Request) {
	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		return
	}
	requestContent := string(dumpRequest)
	fmt.Printf("request content: \n%s\n", requestContent)
}

func (d PrintLogHttpHook) AfterResponse(_ context.Context, response *http.Response) {
	if response == nil {
		return
	}
	dumpResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		return
	}
	responseContent := string(dumpResponse)
	fmt.Printf("response content: \n%s\n", responseContent)
}
