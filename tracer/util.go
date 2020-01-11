package tracer

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func GetRequestID(ctx *gin.Context) string {
	if ctx != nil {
		requestID, exist := ctx.Get(RequestKey)
		if exist {
			return requestID.(string)
		}
	}

	return ""
}

func ChildOfContext(ctx *gin.Context) (childCtx *gin.Context) {
	if ctx == nil {
		return nil
	}

	childCtx = ctx.Copy()

	parent, exist := ctx.Get(TraceContextKey)
	if exist {
		parentSpan, ok := parent.(*jaeger.Span)
		if ok {
			child := opentracing.StartSpan(parentSpan.OperationName(), opentracing.ChildOf(parentSpan.Context()))
			for k, v := range parentSpan.Tags() {
				child.SetTag(k, v)
			}
			childCtx.Set(TraceContextKey, child)
			childCtx.Set(RequestKey, getRequestIDFromTrace(childCtx))
		}
	}

	return
}
