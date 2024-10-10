package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type CommentRouteController struct {
	commentController controllers.CommentController
}

func NewRouteCommentController(commentController controllers.CommentController) CommentRouteController {
	return CommentRouteController{commentController}
}

func (cc *CommentRouteController) CommentRoute(rg *gin.RouterGroup) {

	router := rg.Group("comments")
	router.Use(middleware.DeserializeUser())
	router.POST("/", cc.commentController.CreateComment)
	router.GET("/", cc.commentController.GetAllComments)
	router.PUT("/:commentId", cc.commentController.UpdateComment)
	router.GET("/:commentId", cc.commentController.GetComment)
	router.DELETE("/:commentId", cc.commentController.DeleteComment)
}
