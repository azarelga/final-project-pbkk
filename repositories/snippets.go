package repositories

import (
    "fmt"
    "gorm.io/gorm"
	"time"
)

type Snippet struct {
    ID          string    `json:"id"`
    UserID      uint      `json:"user_id"`              // Foreign key field
    User        User      `gorm:"foreignKey:UserID"`    // Association
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Language    string    `json:"language"`
    Description    string  `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSnippetRequest struct {
    UID    string `form:"username"`
    Title       string `form:"title" binding:"required"`
    Content     string `form:"content" binding:"required"`
    Description string `form:"description" binding:"required"` 
    Language    string `form:"language" binding:"required"`
}

type SnippetRepository struct {
    db *gorm.DB
}

func NewSnippetRepository(db *gorm.DB) *SnippetRepository {
    return &SnippetRepository{db: db}
}

func (r *SnippetRepository) Create(snippet *CreateSnippetRequest) (string, error) {
    // Get count of user's snippets to generate ID
    var count int64
    if err := r.db.Model(&Snippet{}).Where("user_id = ?", snippet.UID).Count(&count).Error; err != nil {
        return "", err
    }

    // Fetch the user's username from the User model
    var user User
    if err := r.db.Where("id = ?", snippet.UID).First(&user).Error; err != nil {
        return "", err
    }
    id :=  fmt.Sprintf("%s-%d", user.Username, count+1)

    // Create a new snippet instance with ID format (username-snippetcount)
    newSnippet := Snippet{
        ID:          id,
        UserID:      user.ID,   // Set the UserID foreign key
        Title:       snippet.Title,
        Content:     snippet.Content,
        Description: snippet.Description,
        Language:    snippet.Language,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    return id, r.db.Create(&newSnippet).Error
}

func (r *SnippetRepository) FindByLanguage(language string) ([]Snippet, error) {
    var snippets []Snippet
    err := r.db.Where("language = ?", language).Preload("User").Find(&snippets).Error
    return snippets, err
}

func (r *SnippetRepository) FindByUserID(uid uint) ([]Snippet, error) {
    var snippets []Snippet
    err := r.db.Preload("User").Where("user_id = ?", uid).Find(&snippets).Error
    return snippets, err
}

func (r *SnippetRepository) FindByID(id string) (*Snippet, error) {
    var snippet Snippet
    err := r.db.Where("id = ?", id).First(&snippet).Error
    return &snippet, err
}
func (r *SnippetRepository) Update(id string, snippet *CreateSnippetRequest) error {
    var existingSnippet Snippet
    if err := r.db.Where("id = ?", id).First(&existingSnippet).Error; err != nil {
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