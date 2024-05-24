package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

const sessionName = "Authorization"

func GetAuthCallBackFunction(c *gin.Context) {
    provider := c.Param("provider")

    req := contextWithProviderName(c, provider)
    c.Request = req
    user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    findUser := FindUser(user.Email)
    if !findUser {
        CreateUser(user.Email, user.Name)
    }

    sessionStore := store.GetSessionStore()

    session, err := sessionStore.Get(c.Request, sessionName)
    if err != nil {
        fmt.Printf("Error retrieving session: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        c.Abort()
        return
    }

    session.Values["user_email"] = user.Email
    session.Values["user_name"] = user.Name
    session.Values["user_id"] = user.UserID
    session.Values["user_picture"] = user.AvatarURL
    session.Values["access_token"] = user.AccessToken

    if err := session.Save(c.Request, c.Writer); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
        return
    }

    frontendURL := os.Getenv("FRONTEND_URL")
    c.Redirect(http.StatusFound, frontendURL)
}




func LogoutHandler(c *gin.Context) {
    provider := c.Param("provider")

    if provider != "google" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
        return
    }

    req := contextWithProviderName(c, provider)
    c.Request = req
    gothic.Logout(c.Writer, c.Request)

    sessionStore := store.GetSessionStore()

    session, err := sessionStore.Get(c.Request, sessionName)
    if err != nil {
        HandleError(c, http.StatusInternalServerError, "Failed to get session")
        return
    }

    for key := range session.Values {
        delete(session.Values, key)
    }
    session.Options.MaxAge = -1 

    if err := session.Save(c.Request, c.Writer); err != nil {
        HandleError(c, http.StatusInternalServerError, "Failed to save session")
        return
    }


    c.Redirect(http.StatusTemporaryRedirect,  os.Getenv("FRONTEND_URL"))
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
        gothic.BeginAuthHandler(c.Writer, c.Request)
    }
}

func contextWithProviderName(c *gin.Context, provider string) *http.Request {
    req := c.Request
    
    ctx := context.WithValue(req.Context(), "provider", provider)
    
    req = req.WithContext(ctx)
    
    return req
}

