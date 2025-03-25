package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	ApplicationJson = "application/json"
)

func (h *Handler) writeProtoResponse(ctx *gin.Context, httpStatusCode int, result proto.Message) {
	buf, err := marshalProtoJson(result)
	if err != nil {
		ctx.PureJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error marshalling response: %w", err)})
		return
	}
	ctx.Data(httpStatusCode, ApplicationJson, buf)
}

func marshalProtoJson(source proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}.Marshal(source)
}
