package mysql

import (
	"context"
	"grpc-template-service/core/kernel"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string { return "mysql" }

func (m *Mod) PreInit(hub *kernel.Hub) error {
	return nil
}

func (m *Mod) Init(hub *kernel.Hub) error {
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	return nil
}
