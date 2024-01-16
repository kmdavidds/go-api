package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kmdavidds/go-api/initializers"
	"github.com/kmdavidds/go-api/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
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

	// Create the user
	user := models.User{Email: body.Email, Password: string(hash), Name: body.Name}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
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

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Compare sent in password with saved user password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

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
		"sub": user.ID,
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
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// Update MemberDays
	var memberDays int
	if user.MemberUntil.After(time.Now()) {
		duration := time.Until(user.MemberUntil)

		memberDays = int(duration.Hours() / 24)
		initializers.DB.Model(&user).Update("MemberDays", memberDays)
		// Update IsMember
		if memberDays > 0 {
			initializers.DB.Model(&user).Update("IsMember", true)
		} else {
			initializers.DB.Model(&user).Update("IsMember", false)
		}
	} else {
		initializers.DB.Model(&user).Update("MemberDays", 0)
		initializers.DB.Model(&user).Update("IsMember", false)
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	// // user.(models.User).ID ----- the method to access the object

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func Feedback(c *gin.Context) {
	// Get id from url
	id := c.Param("id")

	// Get feedback from req body
	var body struct {
		Feedback string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "id = ?", id)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	// Update feedback
	initializers.DB.Model(&user).Update("Feedback", body.Feedback)

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Referral(c *gin.Context) {
	// Get email from req body
	var body struct {
		Email string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	// Increase Points
	if user.HasReferral {
		// Respond
		c.JSON(http.StatusOK, gin.H{
			"message": "referral unavailable",
		})
	} else {
		initializers.DB.Model(&user).Update("Points", user.Points+500)
		initializers.DB.Model(&user).Update("HasReferral", true)
		// Respond
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}

func SetAddress(c *gin.Context) {
	// Get id from url
	id := c.Param("id")

	// Get address from req body
	var body struct {
		Address   string
		Kecamatan string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "id = ?", id)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	// Update address and kecamatan
	lowerCaseKecamatan := strings.ToLower(body.Kecamatan)
	initializers.DB.Model(&user).Update("Address", body.Address)
	initializers.DB.Model(&user).Update("Kecamatan", lowerCaseKecamatan)

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func OrderPetan(c *gin.Context) {
	// Get info off req body
	id := c.Param("id")
	vehicle := c.Param("vehicle")

	// Look up user
	var user models.User
	initializers.DB.First(&user, "id = ?", id)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}
	
	// Get vehicle from url and Create petan object
	var petan models.Petan

	if vehicle == "0" {
		petan = models.Petan{UserEmail: user.Email, Name: user.Name, Kecamatan: user.Kecamatan, Address: user.Address, IsTruk: false}
	} else {
		petan = models.Petan{UserEmail: user.Email, Name: user.Name, Kecamatan: user.Kecamatan, Address: user.Address, IsTruk: true}
	}

	result := initializers.DB.Create(&petan)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create petan object",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}
