package handlers

import (
	"context"

	"example.com/app/models"
	"example.com/app/utils"
)

func (h *HandlersConf) getUserFromContext(ctx context.Context) *models.User {
	//userId := utils.GetUidFromContext(ctx)
	uuid, _ := utils.GetParamsFromContext(ctx)
	user := models.GetUserByUUID(h.DB, uuid)
	return user
}
