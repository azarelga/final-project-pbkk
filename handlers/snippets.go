package handlers

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "snipetty.com/main/services"
    "snipetty.com/main/middleware"
    "snipetty.com/main/repositories"
)

type SnippetHandler struct {
    service *services.SnippetService
}

type LanguageSnippets struct {
    Language string                 // The language name (e.g., "Python", "Go").
    Snippets []repositories.Snippet // The list of snippets for this language.
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
        return
    }

    claims := middleware.JwtClaims(c)
    if claims["id"] == nil {
        c.HTML(http.StatusUnauthorized, "create.html", gin.H{
            "Error": "Unauthorized",
        })
        return
    }
    idFloat, ok := claims["id"].(float64)
    if !ok {
        c.HTML(http.StatusUnauthorized, "create.html", gin.H{
            "Error": "Unauthorized",
        })
        return
    }
    snippet.UID = fmt.Sprintf("%d", uint(idFloat))
    snippetID, err := h.service.CreateSnippet(&snippet)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "create.html", gin.H{
            "Error": err.Error(),
        })
        return
    }
    
    c.Redirect(http.StatusSeeOther, fmt.Sprintf("/snippets/%s", snippetID))
}

func (h *SnippetHandler) GetSnippetsByUsername(c *gin.Context) {
    username := c.Param("username")

    // If username empty, get username from jwt claims
    if username == "" {
        claims := middleware.JwtClaims(c)
        if claims == nil {
            c.HTML(http.StatusUnauthorized, "mylist.html", gin.H{
                "Error": "Unauthorized",
            })
            return
        }
        if usernameClaim, ok := claims["username"].(string); ok {
            username = usernameClaim
        } else {
            c.HTML(http.StatusUnauthorized, "mylist.html", gin.H{
                "Error": "Unauthorized",
            })
            return
        }
    }

    snippets, err := h.service.GetSnippetsByUsername(username)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
            "Error": err.Error(),
        })
        return
    }
    c.HTML(http.StatusOK, "mylist.html", gin.H{
        "snippets": snippets,
    })
}

func (h *SnippetHandler) GetSnippetsByLanguage(c *gin.Context) {
    languages := []string{"Python", "Javascript", "Go", "Rust", "Typescript"}

    // Call the service to get the snippets grouped by language
    groupedSnippets, err := h.service.GetSnippetsByLanguage(languages)
    if err != nil {
        // Handle error by showing it on the page
        c.HTML(http.StatusInternalServerError, "list.html", gin.H{"error": err.Error()})
        return
    }

    // Pass the grouped snippets to the template
    c.HTML(http.StatusOK, "list.html", gin.H{
        "groupedSnippets": groupedSnippets,
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
    var currentUserID uint
    claims := middleware.JwtClaims(c)
    if claims != nil {
        if idFloat, ok := claims["id"].(float64); ok {
            currentUserID = uint(idFloat)
        }
    }
    c.HTML(http.StatusOK, "viewsnippet.html", gin.H{
        "Title": snippet.Title,
        "Username": snippet.User.Username,
        "Language": snippet.Language,
        "Description": snippet.Description,
        "Code": snippet.Content,
        "CreatedAt": snippet.CreatedAt,
        "ID": snippet.ID,
        "IsOwner": currentUserID == snippet.UserID,
    })
}

func (h *SnippetHandler) UpdateSnippet(c *gin.Context) {
    id := c.Param("id")

    // Show edit form for GET requests
    if c.Request.Method == http.MethodGet {
        snippet, err := h.service.GetSnippetByID(id)
        if err != nil {
            c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
                "Error": err.Error(),
            })
            return
        }

        // Check if user owns this snippet
        claims := middleware.JwtClaims(c)
        if claims == nil {
            c.Redirect(http.StatusSeeOther, "/login")
            return
        }
        
        if idFloat, ok := claims["id"].(float64); ok {
            currentUserID := uint(idFloat)
            if currentUserID != snippet.UserID {
                c.HTML(http.StatusForbidden, "edit.html", gin.H{
                    "Error": "Not authorized to edit this snippet",
                })
                return
            }
        }

        c.HTML(http.StatusOK, "edit.html", gin.H{
            "ID": snippet.ID,
            "Title": snippet.Title,
            "Description": snippet.Description,
            "Language": snippet.Language,
            "Content": snippet.Content,
        })
        return
    }

    // Handle PUT request to update snippet
    var updatedSnippet repositories.CreateSnippetRequest
    if err := c.ShouldBind(&updatedSnippet); err != nil {
        c.HTML(http.StatusBadRequest, "edit.html", gin.H{
            "Error": err.Error(),
            "Title": updatedSnippet.Title,
            "Description": updatedSnippet.Description,
            "Language": updatedSnippet.Language,
            "Content": updatedSnippet.Content,
        })
        return
    }

    if err := h.service.UpdateSnippet(id, updatedSnippet); err != nil {
        c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
            "Error": err.Error(),
            "ID": id,
            "Title": updatedSnippet.Title,
            "Description": updatedSnippet.Description,
            "Language": updatedSnippet.Language,
            "Content": updatedSnippet.Content,
        })
        return
    }

    c.Redirect(http.StatusSeeOther, fmt.Sprintf("/snippets/%s", id))
}

func (h *SnippetHandler) DeleteSnippet(c *gin.Context) {
    id := c.Param("id")

    // Show delete confirmation for GET requests
    if c.Request.Method == http.MethodGet {
        snippet, err := h.service.GetSnippetByID(id)
        if err != nil {
            c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
                "Error": err.Error(),
            })
            return
        }

        // Check if user owns this snippet
        claims := middleware.JwtClaims(c)
        if claims == nil {
            c.Redirect(http.StatusSeeOther, "/login")
            return
        }

        if idFloat, ok := claims["id"].(float64); ok {
            currentUserID := uint(idFloat)
            if currentUserID != snippet.UserID {
                c.HTML(http.StatusForbidden, "mylist.html", gin.H{
                    "Error": "Not authorized to delete this snippet",
                })
                return
            }
        }

        c.Redirect(http.StatusOK, "/snippets/listed")
        return
    }

    // Handle DELETE request
    if err := h.service.DeleteSnippet(id); err != nil {
        c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
            "Error": err.Error(),
        })
        return
    }

    c.Redirect(http.StatusSeeOther, "/snippets/listed")
}