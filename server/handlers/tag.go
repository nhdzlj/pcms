package handlers

import (
	"strconv"

	"pcms/services"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	Service *services.TagService
}

func NewTagHandler(service *services.TagService) *TagHandler {
	return &TagHandler{Service: service}
}

// List 获取标签列表
func (h *TagHandler) List(c *gin.Context) {
	userID := c.GetUint64("user_id")
	tags, err := h.Service.List(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.Success(c, tags)
}

// Create 创建标签
func (h *TagHandler) Create(c *gin.Context) {
	var input services.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tag, err := h.Service.Create(c.GetUint64("user_id"), input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, tag)
}

// Delete 删除标签
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.Service.Delete(c.GetUint64("user_id"), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, nil)
}
