package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kmdavidds/go-api/initializers"
	"github.com/kmdavidds/go-api/models"
	"golang.org/x/crypto/bcrypt"
)

func SignupSc(c *gin.Context) {
	// Get the email, password, name off req body
	var body struct {
		Email    string
		Password string
		Name     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	// Create the scavenger
	scavenger := models.Scavenger{Email: body.Email, Password: string(hash), Name: body.Name}

	result := initializers.DB.Create(&scavenger)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func LoginSc(c *gin.Context) {
	// Get the email and password off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up requested scavenger
	var scavenger models.Scavenger
	initializers.DB.First(&scavenger, "email = ?", body.Email)

	if scavenger.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Compare sent in password with saved user password hash
	err := bcrypt.CompareHashAndPassword([]byte(scavenger.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Generate a jwt token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": scavenger.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create JWT",
		})

		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	// time is in seconds, first bool should be true if not in localhost
	c.SetCookie("AuthorizationSc", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func ValidateSc(c *gin.Context) {
	scavenger, _ := c.Get("scavenger")

	// user.(models.User).ID ----- the method to access the object

	c.JSON(http.StatusOK, gin.H{
		"message": scavenger,
	})
}

func ShowPetans(c *gin.Context) {
	// Look up petans today
	var petans []models.Petan
	initializers.DB.Find(&petans, "Created_at > ? AND Created_at < ? AND is_done = ?", time.Now().Add(time.Hour*24*-1), time.Now().Add(time.Hour*24), false)

	if len(petans) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No petans",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"petans": petans,
	})
}

func TakePetan(c *gin.Context) {
	// Get id from url
	id := c.Param("id")

	// Look up petan object by id
	var petan models.Petan
	initializers.DB.Find(&petan, "id = ?", id)

	if petan.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Petan not found",
		})

		return
	}

	// Update IsDone
	initializers.DB.Model(&petan).Update("is_done", true)

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}