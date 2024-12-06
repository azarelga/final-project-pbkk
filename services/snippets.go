package services

import (
    "errors"
    "snipetty.com/main/repositories"
)

type SnippetService struct {
    repo *repositories.SnippetRepository
}

func NewSnippetService(repo *repositories.SnippetRepository) *SnippetService {
    return &SnippetService{repo: repo}
}

func (s *SnippetService) CreateSnippet(input *repositories.CreateSnippetRequest) (string, error) {
    if s.repo == nil {
        return "", errors.New("repository is nil")
    }
    id, err := s.repo.Create(input)
    return id, err
}

func (s *SnippetService) GetAllSnippets() ([]repositories.Snippet,error) {
    if s.repo == nil {
        return nil, errors.New("repository is nil")
    }
    return s.repo.FindAll()
}

func (s *SnippetService) GetSnippetByID(id string) (*repositories.Snippet, error) {
    if s.repo == nil {
        return nil, errors.New("repository is nil")
    }
    return s.repo.FindByID(id)
}

func (s *SnippetService) UpdateSnippet(id string, input repositories.CreateSnippetRequest) (error) {
    if s.repo == nil {
        return errors.New("repository is nil")
    }
    return s.repo.Update(id, &input)
}

func (s *SnippetService) GetSnippetsByUserID(uid uint) ([]repositories.Snippet, error) {
    if s.repo == nil {
        return nil, errors.New("repository is nil")
    }
    return s.repo.FindByUserID(uid)
}

func (s *SnippetService) DeleteSnippet(id string) error {
    if s.repo == nil {
        return errors.New("repository is nil")
    }
    return s.repo.Delete(id)
}