package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth(c *gin.Context) {
    if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
        c.Next()
        return
    }

    // Get token from cookie instead of header
    token, err := c.Cookie("token")
    if err != nil {
        c.Redirect(http.StatusSeeOther, "/login")
        c.Abort()
        return
    }

    // Validate token
    claims := jwt.MapClaims{}
    _, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("SECRET")), nil
    })

    if err != nil {
        c.Redirect(http.StatusSeeOther, "/login")
        c.Abort()
        return
    }

    c.Next()
}