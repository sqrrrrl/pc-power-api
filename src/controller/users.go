package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/infra/entity"
	"github.com/pc-power-api/src/infra/repo"
	"github.com/pc-power-api/src/util"
)

const DeviceCodeLength = 6
const DeviceSecretLength = 12

type UsersHandler struct {
	deviceRepo *repo.DeviceRepository
	userRepo   *repo.UserRepository
}

func NewUsersHandler(e *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware, deviceRepo *repo.DeviceRepository, userRepo *repo.UserRepository) {
	handler := &UsersHandler{
		deviceRepo: deviceRepo,
		userRepo:   userRepo,
	}

	group := e.Group("/user", jwtMiddleware.MiddlewareFunc())
	{
		group.POST("/create-device", handler.createDevice)
	}
}

func (h *UsersHandler) createDevice(c *gin.Context) {
	var deviceInfo *api.DeviceInfo
	err := c.ShouldBind(&deviceInfo)
	if err != nil {
		c.Error(errors.New(err))
		return
	}

	ownerId := middleware.GetUserIdFromContext(c)
	deviceUuid := uuid.New()
	deviceCode := util.GenerateRandomString(DeviceCodeLength)
	deviceSecret := util.GenerateRandomString(DeviceSecretLength)

	var device = entity.Device{
		ID:     deviceUuid.String(),
		Name:   deviceInfo.Name,
		Code:   deviceCode,
		Secret: deviceSecret,
		Status: 0,
		UserID: ownerId,
	}

	aerr := h.deviceRepo.Create(&device)
	if aerr != nil {
		c.Error(aerr)
		return
	}

	c.JSON(200, api.DeviceInfo{
		Name:   device.Name,
		Code:   device.Code,
		Secret: device.Secret,
		Status: device.Status,
	})
}
