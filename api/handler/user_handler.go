package handler

import (
	"fmt"
	"net/http"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UseHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func GetAllUsers(c *gin.Context) {
	users := []model.User{}
	store.DB.Find(&users)

	c.JSON(200, gin.H{
		"user": users,
	})
}

func CreateUser(email, name string) {
	user := model.User{Name: name, Email: email}

	result := store.DB.Create(&user) 

	if result.Error != nil {
		return 
	}
	
	fmt.Println(result)
}

func FindUser(email string) bool {
	user := model.User{}
    result := store.DB.Where("Email = ?", email).First(&user)
    if result.Error == gorm.ErrRecordNotFound {
        return false
    } else if result.Error != nil {
        fmt.Println("Erreur lors de la requête de base de données:", result.Error)
        return false
    }
    return true
}

func GetUserData(c *gin.Context) {

	sessionStore := store.GetSessionStore()

    session, err := sessionStore.Get(c.Request, "session-name")
    if err != nil || session.IsNew {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID := session.Values["user_id"]
    username := session.Values["username"]
	userEmail := session.Values["user_email"]
	userPicture := session.Values["user_picture"]

    c.JSON(http.StatusOK, gin.H{"user_id": userID, "username": username, "user_email": userEmail, "user_picture": userPicture})
}