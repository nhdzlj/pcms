package handlers

import (
	"strconv"

	"pcms/services"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: service}
}

func (h *CategoryHandler) getUint64Param(c *gin.Context, name string) (uint64, error) {
	return strconv.ParseUint(c.Param(name), 10, 64)
}

func (h *CategoryHandler) getUserID(c *gin.Context) uint64 {
	return c.GetUint64("user_id")
}

// GetTree 获取分类树
func (h *CategoryHandler) GetTree(c *gin.Context) {
	tree, err := h.Service.GetTree(h.getUserID(c))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.Success(c, tree)
}

// Create 创建分类
func (h *CategoryHandler) Create(c *gin.Context) {
	var input services.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.Service.Create(h.getUserID(c), input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, category)
}

// Update 更新分类
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := h.getUint64Param(c, "id")
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var input services.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.Service.Update(h.getUserID(c), id, input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, category)
}

// Move 移动分类
func (h *CategoryHandler) Move(c *gin.Context) {
	id, err := h.getUint64Param(c, "id")
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var input services.MoveCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.Service.Move(h.getUserID(c), id, input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, category)
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := h.getUint64Param(c, "id")
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.Service.Delete(h.getUserID(c), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
