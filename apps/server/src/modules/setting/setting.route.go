package setting

import (
	"net/http"
	"peekaping/src/utils"

	"github.com/gin-gonic/gin"
)

type Route struct {
	controller *Controller
}

func NewRoute(
	controller *Controller,
) *Route {
	return &Route{
		controller,
	}
}

func (uc *Route) ConnectRoute(
	rg *gin.RouterGroup,
	controller *Controller,
) {
	router := rg.Group("/settings")

	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.NewSuccessResponse("success", ""))
	})

	router.GET("key/:key", uc.controller.GetByKey)
	router.PUT("key/:key", uc.controller.SetByKey)
	router.DELETE("key/:key", uc.controller.DeleteByKey)
}
