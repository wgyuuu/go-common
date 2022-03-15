package tracer

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHttpTraer(t *testing.T) {
	InitTracer("test")

	engine := gin.New()
	engine.Use(TracerHttp())
	engine.GET("/test", testPrint)
	engine.Run(":8080")
}

func testPrint(ctx *gin.Context) {
	fmt.Println("ori: ", GetRequestID(ctx))
	childCtx := ChildOfContext(ctx)
	fmt.Println("child: ", GetRequestID(childCtx))
}
