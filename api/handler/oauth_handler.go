package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)



func GetAuthCallBackFunction(c *gin.Context) {
    provider := c.Param("provider")
    
    fmt.Println("Provider:", provider)
    
	req := contextWithProviderName(c, provider)
    c.Request = req
    user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Println(user)

    c.Redirect(http.StatusFound, "http://localhost:5173")
}

func LogoutHandler(c *gin.Context) {
    provider := c.Param("provider")

    req := contextWithProviderName(c, provider)
    c.Request = req
    gothic.Logout(c.Writer, c.Request)

    c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/")
}

func AuthHandler(c *gin.Context) {
    provider := c.Param("provider")

    req := contextWithProviderName(c, provider)
    c.Request = req
    gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err == nil {
        c.JSON(http.StatusOK, gin.H{
            "user": gothUser,
        })
    } else {
        fmt.Println("it is time to begin logging in!")
        gothic.BeginAuthHandler(c.Writer, c.Request)
    }
}

func contextWithProviderName(c *gin.Context, provider string) *http.Request {
    req := c.Request
    
    ctx := context.WithValue(req.Context(), "provider", provider)
    
    req = req.WithContext(ctx)
    
    return req
}