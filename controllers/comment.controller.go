package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/models"
	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController(DB *gorm.DB) CommentController {
	return CommentController{DB}
}

func (cc *CommentController) CreateComment(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateCommentRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newComment := models.Comment{
		Content:   payload.Content,
		PostID:    payload.PostID,
		UserID:    currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := cc.DB.Create(&newComment)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newComment})
}

func (cc *CommentController) GetComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")

	var comment models.Comment
	result := cc.DB.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Comment not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": comment})
}

func (cc *CommentController) GetAllComments(ctx *gin.Context) {
	var comments []models.Comment
	result := cc.DB.Find(&comments)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch comments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": comments})
}

func (cc *CommentController) UpdateComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdateComment
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var comment models.Comment
	result := cc.DB.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Comment not found"})
		return
	}

	if comment.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not authorized to update this comment"})
		return
	}

	now := time.Now()
	updateData := models.Comment{
		Content:   payload.Content,
		UpdatedAt: now,
	}

	cc.DB.Model(&comment).Updates(updateData)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": comment})
}

func (cc *CommentController) DeleteComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var comment models.Comment
	result := cc.DB.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Comment not found"})
		return
	}

	if comment.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "You are not authorized to delete this comment"})
		return
	}

	cc.DB.Delete(&comment)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Comment deleted successfully"})
}
