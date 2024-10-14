package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/suidevv/tableye-api/controllers"
	"github.com/suidevv/tableye-api/middleware"
)

type AdminRouteController struct {
	adminController controllers.AdminController
}

func NewRouteAdminController(adminController controllers.AdminController) AdminRouteController {
	return AdminRouteController{adminController}
}

func (rc *AdminRouteController) AdminRoute(rg *gin.RouterGroup) {
	router := rg.Group("admin")
	router.Use(middleware.DeserializeUser(), middleware.AuthorizeRoles("admin"))

	router.POST("/assign-admin", rc.adminController.AssignAdminRole)
}
