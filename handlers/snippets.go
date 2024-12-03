package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "snipetty.com/main/services"
    "snipetty.com/main/repositories"
)

type SnippetHandler struct {
    service *services.SnippetService
}

func NewSnippetHandler(service *services.SnippetService) *SnippetHandler {
    return &SnippetHandler{service: service}
}

func (h *SnippetHandler) CreateSnippet(c *gin.Context) {
    var snippet repositories.CreateSnippetRequest
    if err := c.ShouldBindJSON(&snippet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateSnippet(snippet); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, snippet)
}

func (h *SnippetHandler) GetAllSnippets(c *gin.Context) {
    snippets, err := h.service.GetAllSnippets()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, snippets)
}

func (h *SnippetHandler) GetSnippetByID(c *gin.Context) {
    id := c.Param("id")
    snippet, err := h.service.GetSnippetByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, snippet)
}

func (h *SnippetHandler) UpdateSnippet(c *gin.Context) {
    id := c.Param("id")
    var snippet repositories.CreateSnippetRequest
    if err := c.ShouldBindJSON(&snippet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.UpdateSnippet(id, snippet); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, snippet)
}

func (h *SnippetHandler) DeleteSnippet(c *gin.Context) {
    id := c.Param("id")
    if err := h.service.DeleteSnippet(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Snippet deleted"})
}