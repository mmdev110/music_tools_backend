package handlers

import (
	"context"

	"example.com/app/models"
	"example.com/app/utils"
)

func getUserFromContext(ctx context.Context) *models.User {
	userId := utils.GetUidFromContext(ctx)
	user := models.GetUserByID(userId)
	return user
}
