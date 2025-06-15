package user

import (
	"github.com/DMaryanskiy/bookshare-api/internal/auth"
	"github.com/DMaryanskiy/bookshare-api/internal/task/distributor"
)

type Handler struct {
	TaskDistributor *distributor.TaskDistributor
	TokenStore auth.RefreshTokenStore
}

func NewHandler(dist *distributor.TaskDistributor, ts auth.RefreshTokenStore) *Handler {
	return &Handler{TaskDistributor: dist, TokenStore: ts}
}
