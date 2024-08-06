package examples

import (
	"context"
	"grpc-template-service/core/kernel"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string { return "examples" }

func (m *Mod) PreInit(hub *kernel.Hub) error { return nil }

func (m *Mod) Init(hub *kernel.Hub) error { return nil }

func (m *Mod) Load(hub *kernel.Hub) error {
	//var dbs *redis.Client
	//if hub.Load(&dbs) != nil {
	//	return errors.New("can't load redis client from kernel")
	//}
	//var db *redis.Client
	//if hub.Load(&db) != nil {
	//	return errors.New("can't load redis client from kernel")
	//}
	//if err := dao.Init(db, rdb); err != nil {
	//	return err
	//}
}

func (m *Mod) Start(hub *kernel.Hub) error { return nil }

func (m *Mod) PostInit(hub *kernel.Hub) error { return nil }

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error { return nil }
