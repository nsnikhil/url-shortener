package util

import (
	"encoding/json"
	"net/http"
	"urlshortner/pkg/http/contract"
	"urlshortner/pkg/http/liberr"
)

func writeResponse(code int, data []byte, resp http.ResponseWriter) {
	resp.WriteHeader(code)
	_, _ = resp.Write(data)
}

func writeAPIResponse(code int, ar contract.APIResponse, resp http.ResponseWriter) {
	b, err := json.Marshal(&ar)
	if err != nil {
		// TODO
		writeResponse(http.StatusInternalServerError, []byte("internal server error"), resp)
		return
	}

	writeResponse(code, b, resp)
}

func WriteSuccessResponse(statusCode int, data interface{}, resp http.ResponseWriter) {
	writeAPIResponse(statusCode, contract.NewSuccessResponse(data), resp)
}

func WriteFailureResponse(gr liberr.ResponseError, resp http.ResponseWriter) {
	writeAPIResponse(gr.StatusCode(), contract.NewFailureResponse(gr.ErrorCode(), gr.Error()), resp)
}
