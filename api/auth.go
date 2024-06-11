package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/islombay/cms/model"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func (h *Handler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *Handler) Logout(c *gin.Context) {
	domain := firstOrDefault(os.Getenv("domain"), "localhost")
	c.SetCookie("login", "", -1, "/", domain, false, true)
	c.SetCookie("pwd", "", -1, "/", domain, false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (h *Handler) LoginAPI(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Example authentication (replace with your actual authentication logic)
	var user model.User
	if err := h.db.Where("username = ?", loginData.Username).Where("password = ?", loginData.Password).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
			c.Abort()
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	domain := firstOrDefault(os.Getenv("domain"), "localhost")
	c.SetCookie("login", user.Username, 3600, "/", domain, false, true)
	c.SetCookie("pwd", user.Password, 3600, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func firstOrDefault(f, s string) string {
	if f == "" {
		return s
	}
	return f
}
