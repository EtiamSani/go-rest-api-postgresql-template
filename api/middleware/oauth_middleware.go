package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("randomString"))

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        session, err := store.Get(c.Request, "session-name")
        // if session != nil {
        //     fmt.Println("no session")
        // } else {
        //     // Handle the error case where sessionStore is nil
        // }
        if err != nil || session.IsNew {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}