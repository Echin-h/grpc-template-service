package dao

import (
	"context"
	"gorm.io/gorm"
	"grpc-template-service/internal/mod/hello/model"
)

var helloOrm = &HOrm{}

type HOrm struct {
	*gorm.DB
}

func Init(db *gorm.DB) error {
	ctx := context.TODO()
	helloOrm = &HOrm{DB: db.WithContext(ctx)}
	if helloOrm == nil {
		return gorm.ErrRegistered
	}
	return helloOrm.autoMigrate()
}

func (h *HOrm) autoMigrate() error {
	return h.AutoMigrate(&model.Hello{})
}

func Get() *HOrm {
	return helloOrm
}
