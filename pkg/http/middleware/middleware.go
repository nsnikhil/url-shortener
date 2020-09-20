package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"urlshortner/pkg/http/liberr"
	"urlshortner/pkg/http/util"
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

func WithStatsD(statsd reporters.StatsDClient, api string, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO CHANGE THIS
		hasError := func(code int) bool {
			return code >= 400 && code <= 600
		}

		statsd.ReportAttempt(api)

		cr := util.NewCopyWriter(resp)

		handler(cr, req)
		if hasError(cr.Code()) {
			statsd.ReportFailure(api)
			return
		}

		statsd.ReportSuccess(api)
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
