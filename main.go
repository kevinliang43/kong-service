package main

import (
	"database/sql"
	"dev/kong-service/config"
	"dev/kong-service/service_module"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
)

func main() {

	configuration := getConfig()

	db := setupDatabase(configuration.DatabaseConfig)
	setupServer(configuration.ServerConfig, db)

}

func setupServer(configuration config.ServerConfig, database *sql.DB) {

	serverPath := fmt.Sprintf("%s:%d", configuration.Host, configuration.Port)
	ss := service_module.NewServiceService(database)

	router := gin.Default()

	// Services
	router.GET("services", func(c *gin.Context) { ss.GetAllServices(c) })
	router.GET("services/:serviceId", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		ss.GetService(c, serviceId)
	})

	// Service Records
	router.GET("service-records/:serviceId", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		ss.GetServiceAllRecords(c, serviceId)
	})
	router.GET("service-records/:serviceId/versions/:version", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		version, _ := strconv.ParseFloat(c.Param("version"), 64)
		ss.GetServiceVersion(c, serviceId, version)
	})

	// Service Versions
	router.GET("service-versions/:serviceId", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		ss.GetServiceAllVersions(c, serviceId)
	})

	// create
	router.POST("/services", func(c *gin.Context) { ss.CreateService(c) })

	router.Run(serverPath)

}

func setupDatabase(configuration config.DatabaseConfig) *sql.DB {

	databasePath := fmt.Sprintf("postgresql://%s:%d/%s", configuration.Host, configuration.Port, configuration.Name)

	db, err := sql.Open("pgx", databasePath)
	if err != nil {
		log.Fatalf("Error connecting to database at '%s': %v", databasePath, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping database at '%s': %v", databasePath, err)
	}

	fmt.Printf("Connected to Database at %s\n", databasePath)

	return db

}

func getConfig() config.Config {
	var configuration config.Config

	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return configuration
}
