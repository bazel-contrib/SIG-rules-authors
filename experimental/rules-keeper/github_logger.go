package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ernesto-jimenez/httplogger"
	"github.com/golang/glog"
)

func newLoggedHTTPClient() *http.Client {
	return &http.Client{
		Transport: httplogger.NewLoggedTransport(http.DefaultTransport, httpLogger{}),
	}
}

type httpLogger struct{}

const logDepth = 2

func (httpLogger) LogRequest(req *http.Request) {
	glog.InfoDepth(logDepth, fmt.Sprintf("Request %s %s", req.Method, req.URL))
}

func (httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	if err != nil {
		glog.ErrorDepth(logDepth, err)
		return
	}
	glog.InfoDepth(logDepth, fmt.Sprintf("Response method=%s status=%d durationMs=%d %s",
		req.Method,
		res.StatusCode,
		duration.Milliseconds(),
		req.URL,
	))
}
