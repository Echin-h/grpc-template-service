package pgsql

import (
	"errors"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grpc-template-service/conf"
	"grpc-template-service/core/kernel"
	"grpc-template-service/pkg/colorful"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string { return "pgsql" }

func (m *Mod) PreInit(hub *kernel.Hub) error {
	c := conf.Get().Postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.Addr, c.PORT, c.USER, c.DATABASE, c.PASSWORD)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	hub.Log.Info("pgsql init success")
	hub.Map(&db)
	return nil
}

func (m *Mod) Init(hub *kernel.Hub) error {
	var db *gorm.DB
	if hub.Load(&db) != nil {
		return errors.New("can't load gorm from kernel")
	}
	err := db.Use(otelgorm.NewPlugin())
	if err != nil {
		return err
	}

	var tables []string
	result := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var db *gorm.DB
	if err := hub.Load(&db); err != nil {
		return errors.New("can't load gorm from kernel")
	}

	fmt.Println(colorful.Green("pgsql service Loaded successfully"))
	return nil
}
