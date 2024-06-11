package api

import (
	"github.com/gin-gonic/gin"
	"github.com/islombay/cms/model"
	"net/http"
)

func (h *Handler) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/logout" {
			c.Next()
			return
		}

		uid, err := c.Cookie("login")
		if err != nil || uid == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		pwd, err := c.Cookie("pwd")
		if err != nil || pwd == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var user model.User
		h.db.Where("username = ?", uid).Where("password = ?", pwd).First(&user)

		if user.ID == "" {
			c.Redirect(http.StatusNotModified, "/login")
			c.Abort()
			return
		}

		c.Set("user", user)
		// Continue to the next handler if authenticated
		c.Next()
	}
}

func (h *Handler) authRequiredAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := c.Cookie("login")
		if err != nil || uid == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		pwd, err := c.Cookie("pwd")
		if err != nil || pwd == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user model.User
		h.db.Where("username = ?", uid).Where("password = ?", pwd).First(&user)

		if user.ID == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		// Continue to the next handler if authenticated
		c.Next()
	}
}

func GetContextValue[T any](c *gin.Context, key string) (T, bool) {
	value, exists := c.Get(key)
	if !exists {
		var zero T
		return zero, false
	}

	castValue, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return castValue, true
}

func (h *Handler) IsRole(r string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cu, _ := GetContextValue[model.User](c, "user")

		if cu.Role == r {
			c.Next()
		} else {
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		}
	}
}
