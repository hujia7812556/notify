package server

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"notify/internal/config"
	"notify/internal/dispatcher"
	"notify/internal/parser"
	"notify/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	config     config.ServerConfig
	dispatcher *dispatcher.Dispatcher
	engine     *gin.Engine
}

func New(config config.ServerConfig, dispatcher *dispatcher.Dispatcher) *Server {
	// 设置 gin 模式
	if config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// 使用 gin 的恢复中间件
	engine.Use(gin.Recovery())

	// 添加日志中间件
	engine.Use(func(c *gin.Context) {
		// 请求开始前
		path := c.Request.URL.Path
		start := time.Now()

		c.Next()

		// 请求结束后
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP Request",
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)
	})

	return &Server{
		config:     config,
		dispatcher: dispatcher,
		engine:     engine,
	}
}

// authMiddleware 验证请求的 token
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-API-Token")
		if token == "" {
			logger.Warn("Missing API token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing API token",
			})
			return
		}

		// 使用 subtle.ConstantTimeCompare 进行安全的字符串比较
		if subtle.ConstantTimeCompare([]byte(token), []byte(s.config.Token)) != 1 {
			logger.Warn("Invalid API token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API token",
			})
			return
		}

		c.Next()
	}
}

func (s *Server) registerRoutes() {
	// API 版本分组
	v1 := s.engine.Group("/api/v1")
	{
		// 健康检查接口不需要验证
		v1.GET("/health", s.handleHealth)

		// notify 接口需要验证
		v1.POST("/notify", s.authMiddleware(), s.handleNotify)
	}
}

func (s *Server) handleNotify(c *gin.Context) {
	var msg parser.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 验证消息
	if err := msg.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 分发消息
	s.dispatcher.Dispatch(&msg)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Message accepted",
	})
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) Start() error {
	// 注册路由
	s.registerRoutes()

	// 启动服务器
	addr := fmt.Sprintf(":%d", s.config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	return srv.ListenAndServe()
}
