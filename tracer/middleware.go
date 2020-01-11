package tracer

import (
	"github.com/gin-gonic/gin"
)

func TracerHttp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := TracerStartSpan(getSpanContextFromHttp(ctx.Request), "HTTP_"+ctx.Request.URL.Path, getSpanTagsByHttp(ctx.Request))
		defer span.Finish()

		ctx.Set(TraceContextKey, span)
		ctx.Set(RequestKey, getRequestIDFromTrace(ctx))
		// 执行
		ctx.Next()
	}
}
