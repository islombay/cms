package api

import (
	"github.com/gin-gonic/gin"
	"github.com/islombay/cms/model"
	"gorm.io/gorm"
	"net/http"
)

func Init(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	h := Handler{db: db}

	api := r.Group("/api")
	{
		api.POST("/login", h.LoginAPI)

		api.POST("/create/sub", h.authRequiredAPI(), h.CreateSub)
		api.POST("/create/brigad", h.authRequiredAPI(), h.CreateBrigadir)
		api.POST("/create/task", h.authRequiredAPI(), h.CreateTask)

		api.POST("/assign/task", h.authRequiredAPI(), h.IsRole(string(model.UserBrigad)), h.AssignTask)

		api.POST("/task/status", h.authRequiredAPI(), h.IsRole(string(model.UserBrigad)), h.ChangeTaskStatus)
		api.POST("/task/status/qa", h.authRequiredAPI(), h.IsRole(string(model.UserQA)), h.ChangeTaskStatusQA)

		api.POST("/task/checked", h.authRequiredAPI(), h.IsRole(string(model.UserQA)), h.TaskCheck)
	}

	req := r.Group("/")
	{
		req.Use(h.authRequired())

		req.GET("/login", h.Login)
		req.GET("/logout", h.Logout)

		dashboard := req.Group("/dashboard")
		{
			dashboard.GET("", h.Dashboard)

			dashboard.GET("/gen", h.IsRole(string(model.UserGenPod)), h.DashboardGenPOD)
			dashboard.GET("/sub", h.IsRole(string(model.UserSubPod)), h.DashboardSubPOD)
			dashboard.GET("/brigad", h.IsRole(string(model.UserBrigad)), h.DashboardBrigad)
			dashboard.GET("/qa", h.IsRole(string(model.UserQA)), h.DashboardQA)
		}
	}

	r.GET("/ping", h.Ping)

	r.GET("/:any", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})

	return r
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
	})
}

type Handler struct {
	db *gorm.DB
}
