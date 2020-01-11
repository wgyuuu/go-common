package tracer

import (
	"fmt"
	"testing"

	"git.afpai.com/jxzt/go-common/endless"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHttpTraer(t *testing.T) {
	InitTracer("test")

	engine := gin.New()
	engine.Use(TracerHttp())
	engine.GET("/test", testPrint)
	server := endless.NewServer(":8080", engine)
	assert.Nil(t, server.ListenAndServe(0))
}

func testPrint(ctx *gin.Context) {
	fmt.Println("ori: ", GetRequestID(ctx))
	childCtx := ChildOfContext(ctx)
	fmt.Println("child: ", GetRequestID(childCtx))
}
