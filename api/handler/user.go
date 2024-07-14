package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	internal "github.com/etiamsani/go-rest-api-postgresl-template/api/internal/user"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)



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
	// Vérifie si l'authentification JWT a échoué
    authFailed, exists := c.Get("authFailed")
    if exists && authFailed.(bool) {
        // Si l'authentification JWT a échoué, essayez l'authentification OAuth
        session, err := sessionStore.Get(c.Request, "Authorization")
        if err != nil || session.IsNew {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }
        // Vous pouvez extraire les informations utilisateur de la session OAuth ici
        // et renvoyer les données utilisateur au format JSON
        // Par exemple :
        userID := session.Values["user_id"]
        username := session.Values["username"]
        userEmail := session.Values["user_email"]
        userPicture := session.Values["user_picture"]
        accessToken := session.Values["access_token"]
        c.JSON(http.StatusOK, gin.H{"user_id": userID, "username": username, "user_email": userEmail, "user_picture": userPicture, "access_token": accessToken})
        return
    }

	// sessionStore := store.GetSessionStore()

    // session, err := sessionStore.Get(c.Request, "Authorization")
    // if err != nil || session.IsNew {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    //     return
    // }

    // userID := session.Values["user_id"]
    // username := session.Values["username"]
	// userEmail := session.Values["user_email"]
	// userPicture := session.Values["user_picture"]
	// accesToken := session.Values["acces_token"] 

    // c.JSON(http.StatusOK, gin.H{"user_id": userID, "username": username, "user_email": userEmail, "user_picture": userPicture, "acces_token": accesToken})

	user, exists  := c.Get("user")
	if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Signup(c *gin.Context) {
	var body struct {
		Email string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}


	hash,err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	user := model.User{Email: body.Email, Password: string(hash)}
	result := store.DB.Create(&user)
	
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	
	internal.SendVerificationEmail(c, &user)

	c.JSON(http.StatusOK, gin.H{})
}

func Login(c* gin.Context) {
	var body struct {
		Email string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user model.User
	store.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	err :=bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func VerifyEmail(c *gin.Context) {
	internal.VerifyEmail(c)
}

func ResendVerificationEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is missing"})
		return
	}
	internal.ResendVerificationEmail(c, email) 
}