// package internal
package internal

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/etiamsani/go-rest-api-postgresl-template/api/model"
	"github.com/etiamsani/go-rest-api-postgresl-template/api/store"
	"github.com/mailgun/mailgun-go/v4"

	"crypto/rand"
	"encoding/base64"
)


type MailService struct {
    Domain       string
    APIKey       string
    SMTPLogin    string
    SMTPPassword string
}





// SendMail envoie un e-mail à une ou plusieurs adresses spécifiées.
func (ms *MailService) SendMail(from, to, subject, body string) error {
    mg := mailgun.NewMailgun(ms.Domain, ms.APIKey)

	fromAddr, err := mail.ParseAddress(from)
	if err != nil {
		
		return err
	}

    message := mg.NewMessage(fromAddr.Address, subject, body, to)


    ctx := context.Background()
    _, _, err = mg.Send(ctx, message)
    if err != nil {
        return err
    }

    return nil
}

// SendVerificationEmail envoie un e-mail de vérification à l'utilisateur.
func SendVerificationEmail(c *gin.Context, user *model.User) {
	verificationToken := GenerateRandomString(32)

	

	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"verificationToken": verificationToken,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	
	// Signer le token avec une clé secrète
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate JWT"})
		return
	}

	verificationLink := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", tokenString)
	// mailBody := fmt.Sprintf("Bonjour,<br><br>Cliquez sur le lien suivant pour vérifier votre adresse e-mail : <a href='%s'>Vérifier mon e-mail</a>", verificationLink)

	data := struct {
		VerificationLink string
	}{
		VerificationLink: verificationLink,
	}

	// Définition du modèle HTML
	const emailTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Verify Email</title>
	</head>
	<body>
		<p>Bonjour,</p>
		<p>Cliquez sur le lien suivant pour vérifier votre adresse e-mail :</p>
		<p><a href="{{.VerificationLink}}">Vérifier mon e-mail</a></p>
	</body>
	</html>
	`

	// Analyse du modèle
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatalf("Erreur lors de l'analyse du modèle: %v", err)
	}

	// Application des données au modèle
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		log.Fatalf("Erreur lors de l'application des données au modèle: %v", err)
	}

	// Récupération du corps de l'e-mail au format HTML
	emailBody := tpl.String()

	var (
		mailgunEmailDomain = os.Getenv("MAILGUN_DOMAIN")
		mailgunAPIKey      = os.Getenv("MAILGUN_API_KEY")
	)
	
	

	from := "sender_email@example.com" 
	to := user.Email
	subject := "Email Verification for Bookite"

	mg := mailgun.NewMailgun(mailgunEmailDomain, mailgunAPIKey)

	message := mg.NewMessage(from, subject, "", to)
	message.SetHtml(emailBody)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to send mail"})
		return
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	
	// TODO: Enregistrer les données dans la base de données
	expirationTime := time.Now().Add(24 * time.Hour)
	fmt.Println(verificationToken, "verificationToken")
	verificationData := model.VerificationData{
		Email:     user.Email,
		VerificationToken:     verificationToken,
		ExpiresAt: expirationTime,
	}

	if err := store.DB.Create(&verificationData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store mail verification data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

// GenerateRandomString génère une chaîne aléatoire de longueur n.
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func VerifyEmail(c *gin.Context) {
	// intercepte le token dans le lien 
	paresedURL, err := url.Parse(c.Request.RequestURI)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse URL"})
		return
	}

	token := paresedURL.Query().Get("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}


	// decode le token jwt 

	secret := []byte(os.Getenv("JWT_SECRET"))

	claims, err := decodeJWT(token, secret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
	
	// extraction du email, verifitokenhashé

	email := claims["email"].(string)
	fmt.Println(email)
	verificationToken := claims["verificationToken"].(string)
	fmt.Println(verificationToken)
	exp := claims["exp"].(float64)

	var verificationData model.VerificationData
	if err := store.DB.Find(&verificationData, "email = ?", email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find verification data"})
		return
	}
	hashFromDB := verificationData.VerificationToken
	fmt.Println(hashFromDB, "hashFromDB")
	if hashFromDB != verificationToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification token"})
		return
	} else {
		var user model.User
		if err := store.DB.Find(&user, "email = ?", email).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find user"})
			return
		}
		user.IsVerified = true
		if err := store.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save user"})
			return
		}
		c.Redirect(http.StatusFound, "http://votre-front-end.com/chemin-de-la-page")
	}
	

	if exp < float64(time.Now().Unix()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}

}


func decodeJWT(tokenString string, secret []byte) (jwt.MapClaims, error) {
    // Décoder le JWT
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Vérification de l'algorithme utilisé pour signer le token
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        // Retourner la clé secrète pour vérifier le token
        return secret, nil
    })

    if err != nil {
        return nil, fmt.Errorf("Failed to parse token: %v", err)
    }

    // Extraire les revendications (claims) si le token est valide
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("Invalid token")
}

func ResendVerificationEmail(c *gin.Context, email string) {
	fmt.Println(email, "email")
	var verificationData model.VerificationData
	if err := store.DB.Find(&verificationData, "email = ?", email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find verification data"})
		return
	}
	verificationToken := verificationData.VerificationToken

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": verificationData.Email,
		"verificationToken": verificationToken,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// Signer le token avec une clé secrète
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate JWT"})
		return
	}

	verificationLink := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", tokenString)
	// mailBody := fmt.Sprintf("Bonjour,<br><br>Cliquez sur le lien suivant pour vérifier votre adresse e-mail : <a href='%s'>Vérifier mon e-mail</a>", verificationLink)

	data := struct {
		VerificationLink string
	}{
		VerificationLink: verificationLink,
	}

	// Définition du modèle HTML
	const emailTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Verify Email</title>
	</head>
	<body>
		<p>Bonjour,</p>
		<p>Cliquez sur le lien suivant pour vérifier votre adresse e-mail :</p>
		<p><a href="{{.VerificationLink}}">Vérifier mon e-mail</a></p>
	</body>
	</html>
	`

	// Analyse du modèle
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatalf("Erreur lors de l'analyse du modèle: %v", err)
	}

	// Application des données au modèle
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		log.Fatalf("Erreur lors de l'application des données au modèle: %v", err)
	}

	// Récupération du corps de l'e-mail au format HTML
	emailBody := tpl.String()

	var (
		mailgunEmailDomain = os.Getenv("MAILGUN_DOMAIN")
		mailgunAPIKey      = os.Getenv("MAILGUN_API_KEY")
	)
	
	

	from := "sender_email@example.com" 
	to := verificationData.Email
	subject := "Email Verification for Bookite"

	mg := mailgun.NewMailgun(mailgunEmailDomain, mailgunAPIKey)

	message := mg.NewMessage(from, subject, "", to)
	message.SetHtml(emailBody)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to send mail"})
		return
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

}