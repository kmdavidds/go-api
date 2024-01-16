package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kmdavidds/go-api/initializers"
	"github.com/kmdavidds/go-api/models"
)

func CreateAgenda(c *gin.Context) {
	// Get the email, kecamatan, date off req body
	var body struct {
		CreatorEmail string
		Kecamatan    string
		Date         string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, body.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed parsing date",
			"a":     err,
		})

		return
	}

	// Create the agenda
	agenda := models.Agenda{CreatorEmail: body.CreatorEmail, Kecamatan: body.Kecamatan, Date: parsedDate}

	result := initializers.DB.Create(&agenda)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create agenda",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func ShowAgendas(c *gin.Context) {
	// Get the email, kecamatan, date off req body
	var body struct {
		Kecamatan string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Get the agendas with kecamatan <> and date >
	var agendas []models.Agenda
	initializers.DB.Find(&agendas, "kecamatan = ? AND date > ?", body.Kecamatan, time.Now())

	if len(agendas) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No available agendas",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"agendas": agendas,
	})
}

func TakeAgenda(c *gin.Context) {
	// Get id from url
	id := c.Param("id")
	agendaId, err := strconv.Atoi(c.Param("agendaid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad id formatting",
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

	// Look up requested agenda
	var agenda models.Agenda
	initializers.DB.First(&agenda, "id = ?", agendaId)
	if agenda.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}
	// Create a taken agenda object
	takenAgenda := models.TakenAgenda{AgendaId: uint(agendaId), CreatorEmail: agenda.CreatorEmail, TakerEmail: user.Email, Kecamatan: user.Kecamatan, Address: user.Address}

	result := initializers.DB.Create(&takenAgenda)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to take agenda",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func ShowTasks(c *gin.Context) {
	// Get creator email from req body
	var body struct {
		CreatorEmail string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Get the taken agendas that are not done and with email of scavenger
	var takenAgendas []models.TakenAgenda
	initializers.DB.Find(&takenAgendas, "Is_Done = ? AND Creator_Email = ?", false, body.CreatorEmail)

	if len(takenAgendas) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No available tasks",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"agendas": takenAgendas,
	})
}

func InputDetails(c *gin.Context) {
	// Get task id from url
	taskId := c.Param("taskid")

	// Get details from req body
	var body struct {
		OrganicKilo    uint
		AnorganicKilo  uint
		MetalKilo      uint
		ElectronicKilo uint
		OtherKilo      uint
		Stars          uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up task using task id
	var takenAgenda models.TakenAgenda
	initializers.DB.First(&takenAgenda, "id = ?", taskId)
	if takenAgenda.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task not found",
		})

		return
	}

	// Update task
	initializers.DB.Model(&takenAgenda).Update("OrganicKilo", body.OrganicKilo)
	initializers.DB.Model(&takenAgenda).Update("AnorganicKilo", body.AnorganicKilo)
	initializers.DB.Model(&takenAgenda).Update("MetalKilo", body.MetalKilo)
	initializers.DB.Model(&takenAgenda).Update("ElectronicKilo", body.ElectronicKilo)
	initializers.DB.Model(&takenAgenda).Update("OtherKilo", body.OtherKilo)
	initializers.DB.Model(&takenAgenda).Update("IsDone", true)

	// Look up user using email
	var user models.User
	initializers.DB.First(&user, "email = ?", takenAgenda.TakerEmail)
	if takenAgenda.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	// Update taker points
	points := (body.OrganicKilo + body.AnorganicKilo + body.MetalKilo + body.ElectronicKilo + body.OtherKilo) * 1000 * body.Stars
	initializers.DB.Model(&user).Update("points", points)

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func ShowHistory(c *gin.Context) {
	// Get email from req body
	var body struct {
		TakerEmail string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up taken agendas with taker email and is done
	var takenAgendas []models.TakenAgenda
	initializers.DB.Find(&takenAgendas, "Taker_Email = ? AND Is_Done = ?", body.TakerEmail, true)

	if len(takenAgendas) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No history",
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"takenAgendas": takenAgendas,
	})
}
