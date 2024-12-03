package handlers

import (
	"snipetty.com/main/database"
	"snipetty.com/main/repositories"
	"net/http"
	"os"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Home(c *gin.Context) {
        c.HTML(http.StatusOK, "base.html", gin.H{
            "content": "home.html",
        })
}

func CreateUser(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "base.html", gin.H{
            "content": "register.html",
        })
        return
    }

	var authInput repositories.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error reading json": err.Error()})
		return
	}

	var userFound repositories.User
	database.GetDB().Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error generating password": err.Error()})
		return
	}

	user := repositories.User{
		Username: authInput.Username,
		Password: string(passwordHash),
	}

	database.GetDB().Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func Login(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "base.html", gin.H{
            "content": "login.html",
        })
        return
    }
	var authInput repositories.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error reading json": err.Error()})
		return
	}

	var userFound repositories.User
	database.GetDB().Where("username=?", authInput.Username).Find(&userFound)
	log.Println(userFound)
	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password test"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}