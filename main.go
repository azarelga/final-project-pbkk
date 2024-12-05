package main

import (
	"gorm.io/gorm"
    "html/template"
    "log"

    "snipetty.com/main/services"
    "snipetty.com/main/database"
    "snipetty.com/main/middleware"
    "snipetty.com/main/handlers"
    "snipetty.com/main/repositories"
    "github.com/gin-gonic/gin"
)

var db *gorm.DB
var err error

func init() {
    database.LoadEnvs()
    database.InitializeDatabaseLayer()
    
    // Check if tables exist first
    if !database.TablesExist() {
        log.Println("Tables do not exist. Running migrations...")
        if err := database.Migrate(); err != nil {
            log.Fatalf("Failed to migrate database: %v", err)
        }
        log.Println("Migrations completed successfully")
    } else {
        log.Println("Tables already exist. Skipping migrations")
    }
}

func main() {
    // Create repository
    snippetRepo := repositories.NewSnippetRepository(db)

    // Create service
    snippetService := services.NewSnippetService(snippetRepo)

    // Create handler
    snippetHandler := handlers.NewSnippetHandler(snippetService)

    // setup gin router
    router := gin.Default()
    router.Use(gin.Logger())

    // Setup func map for templates
    router.SetFuncMap(template.FuncMap{
        "isAuthenticated": func(c *gin.Context) bool {
            auth, exists := c.Get("isAuthenticated")
            if !exists {
                return false
            }
            return auth.(bool)
        },
        "getUser": func(c *gin.Context) interface{} {
            user, _ := c.Get("user")
            return user
        },
    })
    // Load HTML templates
    router.LoadHTMLGlob("templates/*")

    // Auth routes
    auth := router.Group("/")
    {
        auth.GET("", middleware.CheckAuth, handlers.Home)
        auth.GET("/login", handlers.Login)
        auth.GET("/register", handlers.CreateUser)
        auth.POST("/login", handlers.Login)
        auth.POST("/register", handlers.CreateUser)
    }

    // Snippet routes with AuthMiddleware
    v1 := router.Group("/snippets", middleware.CheckAuth)
    {
        v1.GET("", snippetHandler.GetAllSnippets)
        v1.POST("/new", snippetHandler.CreateSnippet)
        v1.GET("/:id", snippetHandler.GetSnippetByID)
        v1.PUT("/:id/edit", snippetHandler.UpdateSnippet)
        v1.DELETE("/:id/delete", snippetHandler.DeleteSnippet)
    }

    // start server
    log.Println("starting server on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}