package controller

import (
	"github.com/ZdravkoGyurov/myx-image-api/pkg/config"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/storage/postgresql"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/storage/s3"
)

type Controller struct {
	Config      config.Config
	storage     *postgresql.Storage
	fileStorage *s3.Storage
}

func New(cfg config.Config, psql *postgresql.Storage, s3 *s3.Storage) Controller {
	return Controller{
		Config:      cfg,
		storage:     psql,
		fileStorage: s3,
	}
}
