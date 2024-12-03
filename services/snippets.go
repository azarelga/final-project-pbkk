package services

import (
    "snipetty.com/main/repositories"
)

type SnippetService struct {
    repo *repositories.SnippetRepository
}

func NewSnippetService(repo *repositories.SnippetRepository) *SnippetService {
    return &SnippetService{repo: repo}
}

func (s *SnippetService) CreateSnippet(input repositories.CreateSnippetRequest) error {
    return s.repo.Create(&input)
}

func (s *SnippetService) GetAllSnippets() ([]repositories.Snippet,error) {
    return s.repo.FindAll()
}

func (s *SnippetService) GetSnippetByID(id string) (*repositories.Snippet, error) {
    return s.repo.FindByID(id)
}

func (s *SnippetService) UpdateSnippet(id string, input repositories.CreateSnippetRequest) (error) {
    return s.repo.Update(id, &input)
}

func (s *SnippetService) DeleteSnippet(id string) error {
    return s.repo.Delete(id)
}