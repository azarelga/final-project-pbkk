package handlers

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "snipetty.com/main/services"
    "snipetty.com/main/middleware"
    "snipetty.com/main/repositories"
)

type SnippetHandler struct {
    service *services.SnippetService
}

func NewSnippetHandler(service *services.SnippetService) *SnippetHandler {
    return &SnippetHandler{service: service}
}

func (h *SnippetHandler) CreateSnippet(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "create.html", nil)
        return
    }

    var snippet repositories.CreateSnippetRequest
    if err := c.ShouldBind(&snippet); err != nil {
        c.HTML(http.StatusBadRequest, "create.html", gin.H{
            "Error": err.Error(),
        })
        log.Println(snippet)
        return
    }

    claims := middleware.JwtClaims(c)
    log.Println(claims)
    snippet.Username = claims["username"].(string)
    if err := h.service.CreateSnippet(snippet); err != nil {
        c.HTML(http.StatusInternalServerError, "create.html", gin.H{
            "Error": err.Error(),
        })
        return
    }
    
    c.Redirect(http.StatusSeeOther, "/snippets/:id")
}

func (h *SnippetHandler) GetAllSnippets(c *gin.Context) {
    snippets, err := h.service.GetAllSnippets()
    if err != nil {
        c.HTML(http.StatusInternalServerError, "list.html", gin.H{
            "Error": err.Error(),
        })
        return
    }
    c.HTML(http.StatusOK, "list.html", gin.H{
        "snippets": snippets,
    })
}

func (h *SnippetHandler) GetSnippetByID(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.HTML(http.StatusBadRequest, "home.html", gin.H{
            "Error": "Invalid snippet ID",
        })
        return
    }
    snippet, err := h.service.GetSnippetByID(id)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "home.html", gin.H{
            "Error": err.Error(),
        })
        return
    }
    c.HTML(http.StatusOK, "viewsnippet.html", gin.H{
        "snippet": snippet,
    })
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