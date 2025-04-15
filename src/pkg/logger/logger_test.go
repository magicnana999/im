package logger

import (
	"go.uber.org/zap"
	"testing"
)

func Test(t *testing.T) {

	Init(nil)
	defer Close()

	log := Named("test")
	log.Info("haha", zap.String(Op, OpInit))

}
