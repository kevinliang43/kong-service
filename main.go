package main

import (
	"database/sql"
	"dev/kong-service/config"
	"dev/kong-service/service_manager"
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
	ss := service_manager.NewServiceService(database)

	router := gin.Default()
	router.GET("/services", func(c *gin.Context) { ss.GetServices(c) })
	router.GET("/services/:serviceId", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		ss.GetServiceAllRecords(c, serviceId)
	})
	router.GET("/services/:serviceId/versions", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		ss.GetServiceAllVersions(c, serviceId)
	})
	router.GET("/services/:serviceId/versions/:version", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		version, _ := strconv.ParseFloat(c.Param("version"), 64)
		ss.GetServiceVersion(c, serviceId, version)
	})
	router.GET("/services/:serviceId/versions/latest", func(c *gin.Context) {
		serviceId, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
		ss.GetServiceLatestVersion(c, serviceId)
	})

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
