package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/islombay/cms/model"
	"gorm.io/gorm"
	"net/http"
)

func (h *Handler) AssignTask(c *gin.Context) {
	var data struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cu, _ := GetContextValue[model.User](c, "user")

	var task model.Task
	if err := h.db.Where(&model.Task{
		ID:        data.ID,
		DeletedAt: gorm.DeletedAt{Valid: false},
	}).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not assign task to user", err)
		return
	}

	if task.Responsible != nil || task.IsDone {
		c.JSON(http.StatusGone, gin.H{"error": "already assigned"})
		return
	}

	if err := h.db.Model(&model.Task{ID: data.ID}).Update("responsible", cu.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not assign task to user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ChangeTaskStatus(c *gin.Context) {
	var data struct {
		ID     string `json:"id" binding:"required"`
		Status string `json:"status" binding:"required"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cu, _ := GetContextValue[model.User](c, "user")

	var task model.Task
	if err := h.db.Where(&model.Task{
		ID:        data.ID,
		DeletedAt: gorm.DeletedAt{Valid: false},
	}).Where("responsible = ?", cu.ID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not assign task to user", err)
		return
	}

	if task.Responsible == nil {
		c.JSON(http.StatusGone, gin.H{"error": "not assigned"})
		return
	}

	// check whether previous task is completed
	if data.Status == "done" {
		var prTaskFinished bool = false

		var prTask model.Task
		if err := h.db.Where(&model.Task{
			OrderID:   task.OrderID - 1,
			DeletedAt: gorm.DeletedAt{Valid: false},
		}).First(&prTask).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				prTaskFinished = true
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				fmt.Println("could not find previous task", err)
				return
			}
		}

		if prTask.IsDone {
			prTaskFinished = true
		}

		if prTaskFinished == false {
			c.JSON(http.StatusConflict, gin.H{"error": "Previous task is not completed"})
			return
		}
	}
	var isDone = false
	if data.Status == "done" {
		isDone = true
	}

	if err := h.db.Model(&model.Task{ID: data.ID}).Update("is_done", isDone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not change task status", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ChangeTaskStatusQA(c *gin.Context) {
	var data struct {
		ID   string `json:"id" binding:"required"`
		Done bool   `json:"done"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Println(err)
		return
	}

	var task model.Task
	if err := h.db.Where(&model.Task{
		ID:        data.ID,
		DeletedAt: gorm.DeletedAt{Valid: false},
	}).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not assign task to user", err)
		return
	}

	if task.Responsible == nil {
		c.JSON(http.StatusGone, gin.H{"error": "not assigned"})
		return
	}

	if err := h.db.Model(&model.Task{ID: data.ID}).Update("is_done", data.Done).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not change task status", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) TaskCheck(c *gin.Context) {
	var data struct {
		ID      string `json:"id" binding:"required"`
		IsCheck bool   `json:"isChecked"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Println(err)
		return
	}

	var task model.Task
	if err := h.db.Where(&model.Task{
		ID:        data.ID,
		DeletedAt: gorm.DeletedAt{Valid: false},
	}).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not assign task to user", err)
		return
	}

	if task.Responsible == nil {
		c.JSON(http.StatusGone, gin.H{"error": "not assigned"})
		return
	}

	if !task.IsDone {
		c.JSON(http.StatusForbidden, gin.H{"error": "Task not done"})
		return
	}

	if err := h.db.Model(&model.Task{ID: data.ID}).Update("is_checked", data.IsCheck).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println("could not change task status", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
