package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
	"urlshortner/pkg/http/internal/liberr"
	"urlshortner/pkg/http/internal/util"
	"urlshortner/pkg/reporters"
)

func WithError(handler func(resp http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		err := handler(resp, req)
		switch err.(type) {
		case nil:
			return
		case liberr.ResponseError:
			util.WriteFailureResponse(err.(liberr.ResponseError), resp)
			return
		default:
			util.WriteFailureResponse(liberr.InternalError(err.Error()), resp)
			return
		}

	}
}

func WithReqRespLog(lgr *zap.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		cr := util.NewCopyWriter(resp)

		handler(cr, req)

		b, _ := cr.Body()

		lgr.Sugar().Debug(req)
		lgr.Sugar().Debug(string(b))
	}
}

func WithResponseHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		handler(resp, req)
	}
}

func WithRequestContext(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO: CHANGE TEMP VALUE
		ctx := context.WithValue(req.Context(), "key", "val")
		handler(resp, req.WithContext(ctx))
	}
}

func WithPrometheus(prometheus reporters.Prometheus, api string, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO CHANGE THIS
		hasError := func(code int) bool {
			return code >= 400 && code <= 600
		}

		start := time.Now()
		prometheus.ReportAttempt(api)

		cr := util.NewCopyWriter(resp)

		handler(cr, req)
		if hasError(cr.Code()) {
			duration := time.Since(start)
			prometheus.Observe(api, duration.Seconds())
			prometheus.ReportFailure(api)
			return
		}

		duration := time.Since(start)
		prometheus.Observe(api, duration.Seconds())

		prometheus.ReportSuccess(api)
	}
}
