package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kmdavidds/go-api/initializers"
	"github.com/kmdavidds/go-api/models"
)

func ExchangePoints(c *gin.Context) {
	// Get id from url
	id := c.Param("id")

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "id = ?", id)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	// Check if points are enough
	if user.Points < 5000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not enough points",
		})

		return
	}

	// Check if already subscribed and start or add duration of member
	if user.IsMember {
		initializers.DB.Model(&user).Update("MemberUntil", user.MemberUntil.AddDate(0, 1, 0))
	} else {
		initializers.DB.Model(&user).Update("MemberUntil", time.Now().AddDate(0, 1, 0))
		initializers.DB.Model(&user).Update("IsMember", true)
	}

	// Increment vouchers by 1 and decrement points by 5000
	initializers.DB.Model(&user).Update("Vouchers", user.Vouchers+1)
	initializers.DB.Model(&user).Update("Points", user.Points-5000)

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

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}
