package handler

import "github.com/godpepe7/chirpy/internal/db"

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{db: db}
}
