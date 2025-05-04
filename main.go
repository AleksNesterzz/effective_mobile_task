package main

import (
	"future_today/internal/addition"
	"future_today/internal/config"
	"future_today/internal/controllers"
	"future_today/internal/storage"
	person_service "future_today/services"
	"future_today/utils"
	"log"

	_ "future_today/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files

// @title Person Addition Service API
// @version 1.0
// @description This is a service for adding most often age, gender and nationality to person's name and surname
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {

	//config init
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Can't get config")
	}
	//logger init
	logger := utils.NewLogger()
	//db init
	db, err := storage.InitDb(cfg)
	if err != nil {
		log.Fatal("Can't connect to db")
	}
	//reqs
	ormReq := storage.NewOrmRequestManager(db)
	//services
	add := addition.NewAddition(cfg)
	personService := person_service.NewPersonService(add, ormReq)
	//controllers
	personCtrl := controllers.NewPersonController(personService, logger)
	//router
	router := gin.Default()

	api := router.Group("/personApi/v1")
	{
		api.GET("/persons", personCtrl.GetAllPersons)
		api.GET("/persons/:id", personCtrl.GetPerson)
		api.POST("/persons", personCtrl.CreatePerson)
		api.PUT("/persons/:id", personCtrl.UpdatePerson)
		api.DELETE("/persons/:id", personCtrl.DeletePerson)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	//init server
	logger.Infof("Listening server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
