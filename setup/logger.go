package setup

import (
	"errors"
	"sync"

	onion_log "github.com/kms9/publicyc/pkg/onion-log"
)

var loggerOnce sync.Once
var Logger *onion_log.Log

// StartLogger 初始日志
func StartLogger() error {
	var err error
	loggerOnce.Do(func() {
		Logger = onion_log.UseConfig("logger").Build()
		if Logger == nil {
			err = errors.New("setup.logger start failed")
		}
	})
	return err
}
