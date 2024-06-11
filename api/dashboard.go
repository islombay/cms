package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/islombay/cms/model"
	"gorm.io/gorm"
	"net/http"
)

var userList = []model.User{
	{ID: "hi", Username: "admin", Role: "admin"},
}

func (h *Handler) Dashboard(c *gin.Context) {
	cu, exists := GetContextValue[model.User](c, "user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	if cu.Role == string(model.UserGenPod) {
		c.Redirect(http.StatusFound, "/dashboard/gen")
	} else if cu.Role == string(model.UserSubPod) {
		c.Redirect(http.StatusFound, "/dashboard/sub")
	} else if cu.Role == string(model.UserBrigad) {
		c.Redirect(http.StatusFound, "/dashboard/brigad")
	} else if cu.Role == string(model.UserQA) {
		c.Redirect(http.StatusFound, "/dashboard/qa")
	}
}

func (h *Handler) DashboardGenPOD(c *gin.Context) {
	cu, exists := GetContextValue[model.User](c, "user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	var subPodList []model.User
	if err := h.db.Where("role = ?", "sub_pod").Find(&subPodList).Error; err != nil {
		fmt.Println("could not scan subpod_list", err)
	}

	c.HTML(http.StatusOK, "dashboard_gen_pod.html", gin.H{
		"user":  cu,
		"users": subPodList,
	})
}

func (h *Handler) DashboardSubPOD(c *gin.Context) {
	cu, exists := GetContextValue[model.User](c, "user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	var brigadList []model.User
	if err := h.db.Where("role = ?", string(model.UserBrigad)).Find(&brigadList).Error; err != nil {
		fmt.Println("could not scan brigad_list", err)
	}

	var taskList []model.Task
	if err := h.db.Preload("User").Find(&taskList).Error; err != nil {
		fmt.Println("could not scan task list", err)
	}

	c.HTML(http.StatusOK, "dashboard_sub_pod.html", gin.H{
		"user":  cu,
		"users": brigadList,
		"tasks": taskList,
	})
}
func (h *Handler) DashboardBrigad(c *gin.Context) {
	cu, exists := GetContextValue[model.User](c, "user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	var taskList []model.Task
	if err := h.db.Where(&model.Task{
		DeletedAt: gorm.DeletedAt{
			Valid: false,
		},
	}).Where("responsible is null").Find(&taskList).Error; err != nil {
		fmt.Println("could not scan task list", err)
	}

	var mytasks []model.Task
	if err := h.db.Where(&model.Task{
		Responsible: &cu.ID,
	}).Find(&mytasks).Error; err != nil {
		fmt.Println("could not scan my tasks", err)
	}

	c.HTML(http.StatusOK, "dashboard_brigad.html", gin.H{
		"user":    cu,
		"tasks":   taskList,
		"mytasks": mytasks,
	})
}

func (h *Handler) DashboardQA(c *gin.Context) {
	cu, exists := GetContextValue[model.User](c, "user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	var taskList []model.Task
	if err := h.db.Where(&model.Task{
		DeletedAt: gorm.DeletedAt{
			Valid: false,
		},
	}).Preload("User").Find(&taskList).Error; err != nil {
		fmt.Println("could not scan task list", err)
	}

	c.HTML(http.StatusOK, "dashboard_qa.html", gin.H{
		"user":  cu,
		"tasks": taskList,
	})
}
