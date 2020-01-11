package tracer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

const (
	RequestKey      = "requestid"
	TraceContextKey = "trace-context"
)

func InitTracer(serverName string) {
	tracer, _ := jaeger.NewTracer(serverName, jaeger.NewConstSampler(true), jaeger.NewNullReporter(), jaeger.TracerOptions.CustomHeaderKeys(&jaeger.HeadersConfig{
		JaegerDebugHeader:        jaeger.JaegerDebugHeader,
		JaegerBaggageHeader:      jaeger.JaegerBaggageHeader,
		TraceContextHeaderName:   "trace-id",
		TraceBaggageHeaderPrefix: "bag-",
	}))
	opentracing.SetGlobalTracer(tracer)
}

func TracerStartSpan(parent opentracing.SpanContext, name string, tags map[string]interface{}) opentracing.Span {
	var options []opentracing.StartSpanOption
	if parent != nil {
		options = append(options, opentracing.ChildOf(parent))
	}

	span := opentracing.StartSpan(name, options...)
	for k, v := range tags {
		span.SetTag(k, v)
	}
	return span
}

/// pprivate function
//
func getRequestIDFromTrace(ctx *gin.Context) string {
	if ctx != nil {
		tempSpan, exist := ctx.Get(TraceContextKey)
		if exist {
			spanContext := tempSpan.(opentracing.Span).Context().(jaeger.SpanContext)
			return spanContext.String()
		}
	}
	return ""
}

/// tags function
// http
func getSpanTagsByHttp(req *http.Request) map[string]interface{} {
	return map[string]interface{}{
		ext.SpanKindRPCServer.Key: ext.SpanKindRPCServer.Value,
		string(ext.HTTPMethod):    req.Method,
		string(ext.HTTPUrl):       req.URL.Path,
	}
}

func getSpanContextFromHttp(req *http.Request) opentracing.SpanContext {
	if req.Header != nil {
		spanContext, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		return spanContext
	}
	return nil
}
