package proxy

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type httpError struct {
	StatusCode int
	GrpcCode   string
	StatusText string
}

func writeHTTPError(ctx context.Context, w http.ResponseWriter, err error) {
	code := status.Code(err)
	httpCode := getHTTPStatus(code)

	httpErr := httpError{
		StatusCode: httpCode,
		GrpcCode:   code.String(),
		StatusText: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	jsonStr, err := json.Marshal(httpErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"code\":500,\"error\":\"Internal Server Error\""))
		return
	}
	w.WriteHeader(httpCode)
	w.Write(jsonStr)
	return
}

func getHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
