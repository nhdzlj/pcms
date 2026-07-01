package handlers

import (
	"strconv"

	"pcms/services"
	"pcms/utils"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	Service *services.AttachmentService
}

func NewAttachmentHandler(service *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{Service: service}
}

// List 获取附件列表
func (h *AttachmentHandler) List(c *gin.Context) {
	userID := c.GetUint64("user_id")
	var query services.ListAttachmentQuery
	if docIDStr := c.Query("document_id"); docIDStr != "" {
		docID, err := strconv.ParseUint(docIDStr, 10, 64)
		if err == nil {
			query.DocumentID = &docID
		}
	}

	attachments, err := h.Service.List(userID, query)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.Success(c, attachments)
}

// Create 创建附件记录
func (h *AttachmentHandler) Create(c *gin.Context) {
	var input services.CreateAttachmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint64("user_id")
	var documentID *uint64
	if docIDStr := c.Query("document_id"); docIDStr != "" {
		docID, err := strconv.ParseUint(docIDStr, 10, 64)
		if err == nil {
			documentID = &docID
		}
	}

	attachment, err := h.Service.Create(userID, input, documentID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, attachment)
}

// Delete 删除附件
func (h *AttachmentHandler) Delete(c *gin.Context) {
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

// BindDocument 关联附件到文档
func (h *AttachmentHandler) BindDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var input struct {
		DocumentID uint64 `json:"document_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.Service.BindDocument(c.GetUint64("user_id"), id, input.DocumentID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, nil)
}
