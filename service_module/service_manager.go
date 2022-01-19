package service_module

import (
	"database/sql"
	"dev/kong-service/models"
	"fmt"
)

type ServiceManager struct {
	ServiceDao       *ServiceDao
	ServiceLatestDao *ServiceLatestDao
}

func NewServiceManager(database *sql.DB) *ServiceManager {
	return &ServiceManager{
		ServiceDao:       NewServiceDao(database),
		ServiceLatestDao: NewServiceLatestDao(database),
	}
}

// Fetch all Services
func (sm ServiceManager) GetAllServices() []*models.Service {
	return sm.ServiceLatestDao.GetAllServices()
}

// Fetch Service
func (sm ServiceManager) GetService(serviceId int64) *models.Service {
	return sm.ServiceLatestDao.GetService(serviceId)
}

// Search Services
func (sm ServiceManager) SearchServices(ssr *models.ServicesSearchRequest) (*models.ServicesSearchResponse, error) {

	// validate Limit and Offset params
	if ssr.Limit != nil && (*ssr.Limit > 100 || *ssr.Limit < 0) {
		return nil, fmt.Errorf("'limit' parameter must be between 0 and 100, inclusive. Provided: '%d'",
			*ssr.Limit)
	}
	if ssr.Offset != nil && *ssr.Offset < 0 {
		return nil, fmt.Errorf("'offset' parameter must be greater than 0. Provided: '%d'",
			*ssr.Offset)
	}

	services := sm.ServiceLatestDao.SearchServices(ssr)

	return &models.ServicesSearchResponse{
		Services:   services,
		NextOffset: ssr.GetNextOffset(int64(len(services)))}, nil
}

// Fetch all Service records for a Service Id
func (sm ServiceManager) GetServiceAllRecords(serviceId int64) []*models.ServiceRecord {
	return sm.ServiceDao.GetAllServiceRecordsByServiceId(serviceId)
}

// Fetch Service record for a given Service Id and Version
func (sm ServiceManager) GetServiceVersion(serviceId int64, version float64) *models.ServiceRecord {
	return sm.ServiceDao.getServiceRecordByServiceIdAndVersion(serviceId, version)
}

// Fetch all Service versions for a Service Id
func (sm ServiceManager) GetServiceAllVersions(serviceId int64) []float64 {
	// TODO: Cacheing system for serviceId --> version (TTL: ~15 Seconds)
	return sm.ServiceDao.GetAllServiceVersionsByServiceId(serviceId)
}

// Create new Service record
func (sm ServiceManager) CreateService(newServiceRecord *models.ServiceRecord) (*models.ServiceRecord, error) {

	if newServiceRecord.ServiceId == nil {
		// New service being created
		createdService := sm.ServiceDao.CreateNewService(newServiceRecord)
		sm.ServiceLatestDao.CreateService(createdService)

		return createdService, nil

	} else {
		// Existing Service, new version
		versions := sm.ServiceDao.GetAllServiceVersionsByServiceId(*newServiceRecord.ServiceId)

		if len(versions) == 0 {
			return nil, fmt.Errorf(
				"cannot create a new version record for non-existing serviceId:'%d'",
				*newServiceRecord.ServiceId)
		}

		if versions[len(versions)-1] >= newServiceRecord.Version {
			return nil, fmt.Errorf(
				"new records for existing Services must have a version that is higher than the most up to date version of the existing service. Provided serviceId:'%d', version:'%f'",
				*newServiceRecord.ServiceId,
				newServiceRecord.Version)

		}

		createdServiceRecord := sm.ServiceDao.CreateNewServiceVersion(newServiceRecord)
		sm.ServiceLatestDao.UpdateService(createdServiceRecord)

		return createdServiceRecord, nil
	}

}
