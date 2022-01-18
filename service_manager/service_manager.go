package service_manager

import (
	"database/sql"
	"dev/kong-service/models"
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
func (sm ServiceManager) CreateService(newService *models.Service) *models.Service {

	if newService.ServiceId == nil {
		// New service being created
		return sm.ServiceDao.CreateNewService(newService)
	} else {
		// Existing Service, new version
		// TODO: ensure service exists
		versions := sm.ServiceDao.GetAllServiceVersionsByServiceId(*newService.ServiceId)

		if len(versions) == 0 {
			// TODO: ensure that new version record has a higher version
			// throw error
		}

		if versions[len(versions)-1] > newService.Version {
			// TODO: validate that version number for new version is higher than existing most-uptodate version

		}

		//TODO : check if versions contains the new version, if so error (version exists)
		// maybe this is redundant since we are sorting and checking the biggest one already
		return sm.ServiceDao.CreateNewServiceVersion(newService)
	}

}
