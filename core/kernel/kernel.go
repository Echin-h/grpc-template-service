package kernel

import (
	"context"
	"github.com/juanjiTech/inject/v2"
	"grpc-template-service/conf"
	"grpc-template-service/core/logx"
	"net"
	"sync"
)

type Engine struct {
	config Config

	Ctx    context.Context
	Cancel context.CancelFunc

	ConfigListener []func(*conf.GlobalConfig)

	listener net.Listener

	inject.Injector

	modules   map[string]Module
	modulesMu sync.Mutex
}

type Config struct {
	Listener     net.Listener
	EnableSentry bool
}

func New(config ...Config) *Engine {
	if len(config) == 0 {
		panic("config can't be empty")
	}
	return &Engine{
		config:   config[0],
		listener: config[0].Listener,
		Injector: inject.New(),
		modules:  make(map[string]Module),
	}
}

func (e *Engine) Init() {
	e.Ctx, e.Cancel = context.WithCancel(context.Background())
}

// StartModule start len(modules) goroutines to start modules
func (e *Engine) StartModule() error {
	hub := Hub{
		Injector: e.Injector,
	}

	for _, mod := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + mod.Name())
		if err := mod.PreInit(&hub); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}

	for _, mod := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + mod.Name())
		if err := mod.Init(&hub); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}

	for _, mod := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + mod.Name())
		if err := mod.PostInit(&hub); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}

	for _, mod := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + mod.Name())
		if err := mod.Load(&hub); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}

	for _, mod := range e.modules {
		go func(mod Module) {
			h4m := hub
			h4m.Log = logx.NameSpace("module." + mod.Name())
			if err := mod.Start(&hub); err != nil {
				h4m.Log.Error(err)
				panic(err)
			}
		}(mod)
	}

	return nil
}

func (e *Engine) Serve() {}

func (e *Engine) Stop() error {
	wg := sync.WaitGroup{}
	wg.Add(len(e.modules))
	for _, mod := range e.modules {
		err := mod.Stop(&wg, e.Ctx)
		if err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

// RegMod register modules
func (e *Engine) RegMod(mods ...Module) {
	e.modulesMu.Lock()
	defer e.modulesMu.Unlock()
	for _, mod := range mods {
		if mod.Name() == "" {
			panic("name of module can't be empty")
		}
		if _, ok := e.modules[mod.Name()]; ok {
			panic("module " + mod.Name() + " already exists")
		}
		e.modules[mod.Name()] = mod
	}
}
