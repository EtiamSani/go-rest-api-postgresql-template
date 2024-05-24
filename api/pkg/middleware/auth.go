package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)



func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        
        tokenString, err := c.Cookie("Authorization")
       
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)   
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                return []byte(os.Getenv("JWT_SECRET")), nil
            
            })

            if err != nil || token == nil {
                c.Set("authFailed", true)
                c.Next()
                return
            }
          

            if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

                if float64(time.Now().Unix()) > claims["exp"].(float64) {
                    c.AbortWithStatus(http.StatusUnauthorized)
                    return
                }
                
                var user model.User
                store.DB.First(&user, claims["sub"])

                if user.ID == 0 {
                    c.AbortWithStatus(http.StatusUnauthorized)
                    return
                }

                c.Set("user", user)

                c.Next()
                return 

            } 
        
            c.Set("authFailed", true)
            c.Next()
            
    }
}

var session = sessions.NewCookieStore([]byte("randomString"))
func OAuthAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authFailed, exists := c.Get("authFailed")
        if exists && authFailed.(bool) {
            session, err := session.Get(c.Request, "Authorization")
            if err != nil || session.IsNew {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
                c.Abort()
                return
            }
            c.Next()
            return
        }
        c.Next()
    }
}