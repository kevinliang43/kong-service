package service_module

import (
	"database/sql"
	"dev/kong-service/models"
	"fmt"
)

type ServiceManager struct {
	ServiceDao *ServiceDao
}

func NewServiceManager(database *sql.DB) *ServiceManager {
	return &ServiceManager{
		ServiceDao: NewServiceDao(database),
	}
}

// Fetch all Services
func (sm ServiceManager) GetServices() []*models.Service {
	return sm.ServiceDao.GetAllServiceRecords()
}

// Fetch all Service records for a Service Id
func (sm ServiceManager) GetServiceAllRecords(serviceId int64) []*models.Service {
	return sm.ServiceDao.GetAllServiceRecordsByServiceId(serviceId)
}

// Fetch all Service versions for a Service Id
func (sm ServiceManager) GetServiceAllVersions(serviceId int64) []float64 {
	// TODO: Cacheing system for serviceId --> version (TTL: ~15 Seconds)
	return sm.ServiceDao.GetAllServiceVersionsByServiceId(serviceId)
}

// Fetch Service record (latest version) for a given Service Id
func (sm ServiceManager) GetServiceLatestVersion(serviceId int64) *models.Service {
	return sm.ServiceDao.GetLatestServiceRecordByServiceId(serviceId)
}

// Fetch Service record for a given Service Id and Version
func (sm ServiceManager) GetServiceVersion(serviceId int64, version float64) *models.Service {
	return sm.ServiceDao.getServiceRecordByServiceIdAndVersion(serviceId, version)
}

// Create new Service record
func (sm ServiceManager) CreateService(newService *models.Service) (*models.Service, error) {

	if newService.ServiceId == nil {
		// New service being created
		return sm.ServiceDao.CreateNewService(newService), nil
	} else {
		// Existing Service, new version
		versions := sm.ServiceDao.GetAllServiceVersionsByServiceId(*newService.ServiceId)

		if len(versions) == 0 {
			return nil, fmt.Errorf(
				"cannot create a new version record for non-existing serviceId:'%d'",
				*newService.ServiceId)
		}

		if versions[len(versions)-1] >= newService.Version {
			return nil, fmt.Errorf(
				"new records for existing Services must have a version that is higher than the most up to date version of the existing service. Provided serviceId:'%d', version:'%f'",
				*newService.ServiceId,
				newService.Version)

		}

		return sm.ServiceDao.CreateNewServiceVersion(newService), nil
	}

}
