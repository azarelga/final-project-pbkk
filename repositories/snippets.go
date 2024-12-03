package repositories

import (
    "gorm.io/gorm"
	"time"
)

type Snippet struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Language    string    `json:"language"`
    Description    string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSnippetRequest struct {
    Title       string `json:"title" binding:"required"`
    Content     string `json:"content" binding:"required"`
    Description string `json:"description" binding:"required"` 
    Language    string `json:"language" binding:"required"`
}

type SnippetRepository struct {
    db *gorm.DB
}

func NewSnippetRepository(db *gorm.DB) *SnippetRepository {
    return &SnippetRepository{db: db}
}

func (r *SnippetRepository) Create(snippet *CreateSnippetRequest) (error) {
    return r.db.Create(snippet).Error
}

func (r *SnippetRepository) FindAll() ([]Snippet, error) {
    var snippets []Snippet
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