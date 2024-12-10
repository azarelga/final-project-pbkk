package services

import (
    "errors"
    "snipetty.com/main/repositories"
)

type LanguageSnippets struct {
    Language string                 // The language name (e.g., "Python", "Go").
    Snippets []repositories.Snippet // The list of snippets for this language.
}

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

func (s *SnippetService) GetSnippetsByLanguage(languages []string) ([]LanguageSnippets, error) {
    groupedSnippets := []LanguageSnippets{}

    for _, lang := range languages {
        snippets, err := s.repo.FindByLanguage(lang)
        if err != nil {
            return nil, err
        }
        groupedSnippets = append(groupedSnippets, LanguageSnippets{
            Language: lang,
            Snippets: snippets,
        })
    }

    return groupedSnippets, nil
}

func (s *SnippetService) GetSnippetsByUsername(username string) ([]repositories.Snippet, error) {
    if s.repo == nil {
        return nil, errors.New("repository is nil")
    }
    return s.repo.FindByUsername(username)
}

func (s *SnippetService) DeleteSnippet(id string) error {
    if s.repo == nil {
        return errors.New("repository is nil")
    }
    return s.repo.Delete(id)
}