package ginPprof

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"grpc-template-service/core/kernel"
	"net/http"
	"net/http/pprof"
)

// this module is used to add pprof to gin
// pprof is a tool for golang to profile the performance of the program
// https://pkg.go.dev/runtime/pprof

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string {
	return "jinPprof"
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var ginE *gin.Engine
	err := hub.Load(&ginE)
	if err != nil {
		return errors.New("can't load gin.Engine from kernel")
	}

	// 定义 Basic 认证的凭据
	authStr := "Basic " + base64.StdEncoding.EncodeToString([]byte("pprof:nemertes"))

	// 创建一个路由组用于 pprof，并添加中间件
	pprofGroup := ginE.Group("/debug/pprof")
	pprofGroup.Use(func(c *gin.Context) {
		// 从请求中获取 Authorization 头部
		auth := c.Request.Header.Get("Authorization")
		// 检查 Authorization 是否正确
		if auth != authStr {
			// 如果不正确，返回 401 状态码，并要求认证
			c.Writer.Header().Set("WWW-Authenticate", "Basic")
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// 如果正确，调用下一个 Handler
		c.Next()
	})

	// 注册 pprof 的路由处理器
	pprofGroup.GET("/", func(c *gin.Context) { pprof.Index(c.Writer, c.Request) })
	pprofGroup.GET("/cmdline", func(c *gin.Context) { pprof.Cmdline(c.Writer, c.Request) })
	pprofGroup.GET("/profile", func(c *gin.Context) { pprof.Profile(c.Writer, c.Request) })
	pprofGroup.POST("/symbol", func(c *gin.Context) { pprof.Symbol(c.Writer, c.Request) })
	pprofGroup.GET("/trace", func(c *gin.Context) { pprof.Trace(c.Writer, c.Request) })
	pprofGroup.GET("/allocs", func(c *gin.Context) { pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request) })
	pprofGroup.GET("/block", func(c *gin.Context) { pprof.Handler("block").ServeHTTP(c.Writer, c.Request) })
	pprofGroup.GET("/goroutine", func(c *gin.Context) { pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request) })
	pprofGroup.GET("/heap", func(c *gin.Context) { pprof.Handler("heap").ServeHTTP(c.Writer, c.Request) })
	pprofGroup.GET("/mutex", func(c *gin.Context) { pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request) })
	pprofGroup.GET("/threadcreate", func(c *gin.Context) { pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request) })

	return nil
}
