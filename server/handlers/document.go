package handlers

import (
	"strconv"

	"pcms/services"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	Service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{Service: service}
}

func (h *DocumentHandler) getUserID(c *gin.Context) uint64 {
	return c.GetUint64("user_id")
}

func (h *DocumentHandler) getID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

// List 文档列表
func (h *DocumentHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)

	var categoryID *uint64
	if cidStr := c.Query("category_id"); cidStr != "" {
		cid, err := strconv.ParseUint(cidStr, 10, 64)
		if err == nil {
			categoryID = &cid
		}
	}

	status := c.DefaultQuery("status", "")

	result, err := h.Service.List(h.getUserID(c), page, pageSize, categoryID, status)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.Success(c, result)
}

// Search 搜索文档
func (h *DocumentHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.BadRequest(c, "搜索关键词不能为空")
		return
	}

	page, pageSize := utils.GetPagination(c)

	result, err := h.Service.Search(h.getUserID(c), keyword, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.Success(c, result)
}

// Get 文档详情
func (h *DocumentHandler) Get(c *gin.Context) {
	id, err := h.getID(c)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	document, err := h.Service.GetByID(h.getUserID(c), id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, document)
}

// Create 创建文档
func (h *DocumentHandler) Create(c *gin.Context) {
	var input services.CreateDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	document, err := h.Service.Create(h.getUserID(c), input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, document)
}

// Update 更新文档
func (h *DocumentHandler) Update(c *gin.Context) {
	id, err := h.getID(c)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var input services.UpdateDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	document, err := h.Service.Update(h.getUserID(c), id, input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, document)
}

// Delete 删除文档
func (h *DocumentHandler) Delete(c *gin.Context) {
	id, err := h.getID(c)
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
