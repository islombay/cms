package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/islombay/cms/model"
	"gorm.io/gorm"
	"net/http"
)

func (h *Handler) CreateSub(c *gin.Context) {
	var data struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cu, _ := GetContextValue[model.User](c, "user")
	if cu.Role != string(model.UserGenPod) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.db.Create(&model.User{
		ID:       uuid.NewString(),
		Username: data.Username,
		Password: data.Password,
		Role:     string(model.UserSubPod),
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CreateBrigadir(c *gin.Context) {
	var data struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cu, _ := GetContextValue[model.User](c, "user")
	if cu.Role != string(model.UserSubPod) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.db.Create(&model.User{
		ID:       uuid.NewString(),
		Username: data.Username,
		Password: data.Password,
		Role:     string(model.UserBrigad),
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var data struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cu, _ := GetContextValue[model.User](c, "user")
	if cu.Role != string(model.UserSubPod) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.db.Create(&model.Task{
		ID:          uuid.NewString(),
		Name:        data.Name,
		Responsible: nil,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
