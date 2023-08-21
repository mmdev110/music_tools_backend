package handlers

import (
	"context"

	"example.com/app/models"
	"example.com/app/utils"
)

func (h *Base) getUserFromContext(ctx context.Context) *models.User {
	userId := utils.GetUidFromContext(ctx)
	user := models.GetUserByID(h.DB, userId)
	return user
}
