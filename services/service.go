package services

import (
	"github.com/gin-gonic/gin"
	"ginblog/setup"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

// Service 调用总方法
type Service struct {
	BaseContent      *logger.BaseContentInfo
	Ctx              *gin.Context
	CtxLog           *logger.Logger
	Request          interface{}
}

// Get 写入扩展后的context,并返回去全部service信息
func Get(ctx *gin.Context) *Service {
	baseContent := onion_log.GetBaseByContext(ctx)
	s := &Service{
		BaseContent: baseContent,
		Ctx:         ctx,
		CtxLog:      setup.Logger.With(baseContent),
		Request:     nil,
	}

	ctx.Set("service", s)
	return s
}

// get 获取扩展后的context
func get(c *gin.Context) *Service {
	return c.MustGet("service").(*Service)
}

// ServiceInit 部分初始化,service部分,需要启动时,自启动的项目
func ServiceInit() error {
	return nil
}
