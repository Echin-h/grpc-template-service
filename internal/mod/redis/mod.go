package redis

import "grpc-template-service/core/kernel"

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string { return "redis" }

func (m *Mod) PreInit(hub *kernel.Hub) error {
	return nil
}

func (m *Mod) Init(hub *kernel.Hub) error {
	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error {
	return nil
}
