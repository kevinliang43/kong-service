package service_module

import (
	"database/sql"
	"dev/kong-service/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceService struct {
	ServiceManager *ServiceManager
}

func NewServiceService(database *sql.DB) *ServiceService {
	return &ServiceService{
		ServiceManager: NewServiceManager(database),
	}
}

// Fetch all Services
func (ss ServiceService) GetServices(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ss.ServiceManager.GetServices())
}

// Fetch all Service records for a Service Id
func (ss ServiceService) GetServiceAllRecords(c *gin.Context, serviceId int64) {
	c.IndentedJSON(http.StatusOK, ss.ServiceManager.GetServiceAllRecords(serviceId))
}

// Fetch all Service versions for a Service Id
func (ss ServiceService) GetServiceAllVersions(c *gin.Context, serviceId int64) {
	c.IndentedJSON(http.StatusOK, ss.ServiceManager.GetServiceAllVersions(serviceId))
}

// Fetch Service record (latest version) for a given Service Id
func (ss ServiceService) GetServiceLatestVersion(c *gin.Context, serviceId int64) {

	if result := ss.ServiceManager.GetServiceLatestVersion(serviceId); result != nil {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no records found for serviceId '%d'", serviceId)})
	}
}

// Fetch Service record for a given Service Id and Version
func (ss ServiceService) GetServiceVersion(c *gin.Context, serviceId int64, version float64) {

	if result := ss.ServiceManager.GetServiceVersion(serviceId, version); result != nil {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no records found for serviceId '%d' and version '%f'", serviceId, version)})
	}
}

// Create new Service record
func (ss ServiceService) CreateService(c *gin.Context) {
	var newService models.Service

	if err := c.BindJSON(&newService); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result, err := ss.ServiceManager.CreateService(&newService); err == nil {
		c.IndentedJSON(http.StatusCreated, result)
	} else {
		println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
