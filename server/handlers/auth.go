package handlers

import (
	"pcms/services"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.Service.Register(input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, result)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.Service.Login(input)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, result)
}

// Me 获取当前用户信息
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint64("user_id")
	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, user)
}
