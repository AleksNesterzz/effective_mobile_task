package controllers

import (
	"future_today/models"
	services "future_today/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PersonController struct {
	service *services.PersonService
	logger  *logrus.Logger
}

func NewPersonController(service *services.PersonService, logger *logrus.Logger) *PersonController {
	return &PersonController{
		service: service,
		logger:  logger,
	}
}

// @Summary Create a new person
// @Description Create a new person with the input payload
// @Tags persons
// @Accept  json
// @Produce  json
// @Param person body models.CreatePersonRequest true "Create person"
// @Success 200 {object} models.PersonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons [post]
func (c *PersonController) CreatePerson(ctx *gin.Context) {
	var req models.CreatePersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	person, err := c.service.CreatePerson(&req)
	if err != nil {
		c.logger.Errorf("Error creating person: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := models.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
		IsActive:    person.IsActive,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get all persons
// @Description Get all persons with optional filters
// @Tags persons
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param name query string false "Name filter"
// @Param surname query string false "Surname filter"
// @Param min_age query int false "Minimum age filter"
// @Param max_age query int false "Maximum age filter"
// @Param gender query string false "Gender filter"
// @Param nationality query string false "Nationality filter"
// @Success 200 {array} models.PersonResponse
// @Failure 500 {object} map[string]string
// @Router /persons [get]
func (c *PersonController) GetAllPersons(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	name := ctx.Query("name")
	surname := ctx.Query("surname")
	patronymic := ctx.Query("patronymic")
	minAgeStr := ctx.Query("min_age")
	maxAgeStr := ctx.Query("max_age")
	gender := ctx.Query("gender")
	nation := ctx.Query("nation")

	var minAge, maxAge *int

	if minAgeStr != "" {
		age, err := strconv.Atoi(minAgeStr)
		if err == nil {
			minAge = &age
		}
	}
	if maxAgeStr != "" {
		age, err := strconv.Atoi(maxAgeStr)
		if err == nil {
			maxAge = &age
		}
	}

	var namePtr, surnamePtr, genderPtr, nationPtr, patronymicPtr *string
	if name != "" {
		namePtr = &name
	}
	if surname != "" {
		surnamePtr = &surname
	}
	if patronymic != "" {
		patronymicPtr = &patronymic
	}
	if gender != "" {
		genderPtr = &gender
	}
	if nation != "" {
		nationPtr = &nation
	}

	persons, err := c.service.GetAllPersons(
		limit,
		offset,
		namePtr,
		surnamePtr,
		patronymicPtr,
		minAge,
		maxAge,
		genderPtr,
		nationPtr,
	)
	if err != nil {
		c.logger.Errorf("Error getting persons: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]models.PersonResponse, len(persons))
	for i, person := range persons {
		response[i] = models.PersonResponse{
			ID:          person.ID,
			Name:        person.Name,
			Surname:     person.Surname,
			Patronymic:  person.Patronymic,
			Age:         person.Age,
			Gender:      person.Gender,
			Nationality: person.Nationality,
			IsActive:    person.IsActive,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get a person by ID
// @Description Get a person by ID
// @Tags persons
// @Accept  json
// @Produce  json
// @Param id path int true "Person ID"
// @Success 200 {object} models.PersonResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons/{id} [get]
func (c *PersonController) GetPerson(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.logger.Errorf("Error parsing ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	person, err := c.service.GetPerson(uint(id))
	if err != nil {
		c.logger.Errorf("Error getting person: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	response := models.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
		IsActive:    person.IsActive,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Update a person
// @Description Update a person by ID
// @Tags persons
// @Accept  json
// @Produce  json
// @Param id path int true "Person ID"
// @Param person body models.UpdatePersonRequest true "Update person"
// @Success 200 {object} models.PersonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons/{id} [put]
func (c *PersonController) UpdatePerson(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.logger.Errorf("Error parsing ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var dto models.UpdatePersonRequest
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		c.logger.Errorf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := c.service.UpdatePerson(uint(id), &dto)
	if err != nil {
		c.logger.Errorf("Error updating person: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.PersonResponse{
		ID:          person.ID,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
		IsActive:    person.IsActive,
	}

	ctx.JSON(http.StatusOK, response)
}

// мб сделать действительное удаление ?

// @Summary Delete a person
// @Description Delete a person by ID (soft delete)
// @Tags persons
// @Accept  json
// @Produce  json
// @Param id path int true "Person ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons/{id} [delete]
func (c *PersonController) DeletePerson(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.logger.Errorf("Error parsing ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = c.service.DeletePerson(uint(id))
	if err != nil {
		c.logger.Errorf("Error deleting person: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "person deleted successfully"})
}
