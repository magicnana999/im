package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestDefault(t *testing.T) {

	Init(nil)
	defer Close()

	log := NameWithOptions("test", zap.AddCaller())
	log.Info("haha")
	log.Info("hehe")

}
