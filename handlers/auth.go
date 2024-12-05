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
        c.HTML(http.StatusOK, "home.html", nil)
}

func CreateUser(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "register.html", nil)
        return
    }

	var authInput repositories.AuthInput

	if err := c.ShouldBind(&authInput); err != nil {
        c.HTML(http.StatusOK, "register.html", gin.H{"Error": err.Error()})
		return
	}

	var userFound repositories.User
	database.GetDB().Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID != 0 {
        c.HTML(http.StatusOK, "register.html", gin.H{"Error": "Username already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
        c.HTML(http.StatusOK, "register.html", gin.H{"Error": "Failed to hash password"})
		return
	}

	user := repositories.User{
		Username: authInput.Username,
		Password: string(passwordHash),
	}

	database.GetDB().Create(&user)

	c.HTML(http.StatusOK, "login.html", gin.H{"Success": "User created successfully"})
}

func Login(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "login.html", nil)
        return
    }
	var authInput repositories.AuthInput

	if err := c.ShouldBind(&authInput); err != nil {
        c.HTML(http.StatusOK, "login.html", gin.H{"Error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var userFound repositories.User
	database.GetDB().Where("username=?", authInput.Username).Find(&userFound)
	log.Println(userFound)
	if userFound.ID == 0 {
        c.HTML(http.StatusOK, "login.html", gin.H{"Error": "Invalid username or password"})
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
        c.HTML(http.StatusOK, "login.html", gin.H{"Error": "Invalid username or password"})
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid username or password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
        c.HTML(http.StatusOK, "login.html", gin.H{"Error": "Error generating token"})
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error generating token"})
	}
    // Set token in cookie
	c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie(
        "Authorization",           // name
        token,            // value
        3600 * 24,       // max age (24 hours)
        "",             // path
        "",              // domain
        false,           // secure
        false,            // httpOnly
    )

    c.Redirect(http.StatusSeeOther, "/")
}

func Logout(c *gin.Context) {
	c.SetCookie(
		"Authorization",
		"", // value
		-1, // max age
		"", // path
		"", // domain
		false, // secure
		false, // httpOnly
	)
	c.Redirect(http.StatusSeeOther, "/")
}