package handler

import (
	"github.com/Pineapple217/Corap-web/pkg/database"
)

type Handler struct {
	DB *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{
		DB: db,
	}
}
