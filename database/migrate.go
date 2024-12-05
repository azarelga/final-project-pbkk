package database

import (
    "log"
	"snipetty.com/main/repositories"	
)

func TablesExist() bool {
    db := GetDB()
    
    // Check if tables exist
    hasUser := db.Migrator().HasTable(&repositories.User{})
    hasSnippet := db.Migrator().HasTable(&repositories.Snippet{})
    
    return hasUser && hasSnippet
}

func Migrate() error {
    db := GetDB()
    
    // Run migrations
    err := db.AutoMigrate(
        &repositories.User{},
        &repositories.Snippet{}, // Note: Changed from Snippets to Snippet to match model name
    )
    
    if err != nil {
        log.Printf("Failed to migrate database: %v", err)
        return err
    }
    
    log.Println("Database migration completed successfully")
    return nil
}