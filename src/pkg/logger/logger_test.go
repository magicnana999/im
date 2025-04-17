package logger

import (
	"github.com/magicnana999/im/define"
	"go.uber.org/zap"
	"testing"
)

func TestDefault(t *testing.T) {

	Init(nil)
	defer Close()

	log := Named("test")
	log.Info("haha", zap.String(define.OP, define.OpInit))
	log.Info("hehe")

}
