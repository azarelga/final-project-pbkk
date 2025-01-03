package main

import (
	"gorm.io/gorm"
    "log"

    "snipetty.com/main/services"
    "snipetty.com/main/database"
    "snipetty.com/main/middleware"
    "snipetty.com/main/handlers"
    "snipetty.com/main/repositories"
    "github.com/gin-gonic/gin"
)

var db *gorm.DB

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
    db = database.GetDB()
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

    // Load HTML templates
    router.LoadHTMLGlob("templates/*")

    // Auth routes
    auth := router.Group("/")
    {
        auth.GET("", handlers.Home)
        auth.GET("/login", handlers.Login)
        auth.GET("/logout", handlers.Logout)
        auth.GET("/register", handlers.CreateUser)
        auth.POST("/login", handlers.Login)
        auth.POST("/register", handlers.CreateUser)
    }

    // Snippet routes
    snip := router.Group("/snippets")
    {
        // Guest routes
        snip.GET("", snippetHandler.GetSnippetsByLanguage)
        snip.GET("/user/:username", snippetHandler.GetSnippetsByUsername)
        snip.GET("/:id", snippetHandler.GetSnippetByID)
        
        // Authenticated routes
        snip.GET("/my", middleware.CheckAuth, snippetHandler.GetSnippetsByUsername)
        snip.GET("/new", middleware.CheckAuth,snippetHandler.CreateSnippet)
        snip.POST("/new", middleware.CheckAuth,snippetHandler.CreateSnippet)
        snip.GET("/:id/edit",middleware.CheckAuth, snippetHandler.UpdateSnippet)
        snip.POST("/:id/edit",middleware.CheckAuth, snippetHandler.UpdateSnippet)
        snip.POST("/:id/delete", middleware.CheckAuth, snippetHandler.DeleteSnippet)
        snip.GET("/:id/delete", middleware.CheckAuth, snippetHandler.DeleteSnippet)
    }

    // start server
    log.Println("starting server on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}