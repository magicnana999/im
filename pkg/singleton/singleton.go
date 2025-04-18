package singleton

import (
	"github.com/magicnana999/im/pkg/logger"
	"go.uber.org/zap"
	"sync"
)

type Singleton[T any] struct {
	instance T
	once     sync.Once
}

func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{}
}

func (s *Singleton[T]) Get(initFunc func() T) T {
	s.once.Do(func() {
		s.instance = initFunc()
	})
	return s.instance
}

func (s *Singleton[T]) GetWithCloser(initFunc func() (T, func(), error)) (T, func()) {
	var closer func()
	s.once.Do(func() {
		var err error
		s.instance, closer, err = initFunc()
		if err != nil {
			logger.Fatal("Singleton init failed", zap.Error(err))
		}
	})
	return s.instance, closer
}
