package repositories

import (
    "log"
    "fmt"
    "gorm.io/gorm"
	"time"
)

type Snippet struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Language    string    `json:"language"`
    Description    string  `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSnippetRequest struct {
    Title       string `form:"title" binding:"required"`
    Content     string `form:"content" binding:"required"`
    Description string `form:"description" binding:"required"` 
    Language    string `form:"language" binding:"required"`
    Username    string `form:"username"`
}

type SnippetRepository struct {
    db *gorm.DB
}

func NewSnippetRepository(db *gorm.DB) *SnippetRepository {
    return &SnippetRepository{db: db}
}

func (r *SnippetRepository) Create(snippet *CreateSnippetRequest) error {
    // Get count of user's snippets to generate ID
    var count int64
    var username = snippet.Username
    r.db.Model(&Snippet{}).Where("id LIKE ?", username + "-%").Count(&count)
    
    var new Snippet 
    new.ID = fmt.Sprintf("%s-%d", username, count+1)
    new.Title = snippet.Title
    new.Content = snippet.Content
    new.Description = snippet.Description
    new.Language = snippet.Language
    new.CreatedAt = time.Now()
    new.UpdatedAt = time.Now()
    
    return r.db.Create(&new).Error
}

func (r *SnippetRepository) FindAll() ([]Snippet, error) {
    var snippets []Snippet
    log.Println(r.db)
    err := r.db.Find(&snippets).Error
    return snippets, err
}

func (r *SnippetRepository) FindByID(id string) (*Snippet, error) {
    var snippet Snippet
    err := r.db.First(&snippet, id).Error
    return &snippet, err
}
func (r *SnippetRepository) Update(id string, snippet *CreateSnippetRequest) error {
    var existingSnippet Snippet
    if err := r.db.First(&existingSnippet, id).Error; err != nil {
        return err
    }

    // Update the fields of the existing snippet
    existingSnippet.Title = snippet.Title
    existingSnippet.Language = snippet.Language
    existingSnippet.Content = snippet.Content
    existingSnippet.Description = snippet.Description

    return r.db.Save(&existingSnippet).Error
}

func (r *SnippetRepository) Delete(id string) error {
    return r.db.Delete(&Snippet{}, id).Error
}