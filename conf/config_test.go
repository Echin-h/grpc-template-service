package conf

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	_, err := os.Stat("config.yaml")
	if err != nil {
		t.Error("os.Stat() = ", err)
	}
	//if stat == nil {
	LoadConfig()
	conf := Get()
	if conf == nil {
		t.Error("Get() = nil")
	}
}
