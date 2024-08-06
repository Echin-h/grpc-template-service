package logx

import (
	"go.uber.org/zap/zapcore"
	"grpc-template-service/conf"
	"testing"
)

func TestInit(t *testing.T) {
	conf.LoadConfig()

	Init(zapcore.DebugLevel)

	suger := NameSpace("test")

	suger.Info("test")
}
