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
	latestServices := sm.ServiceLatestDao.GetAllServices()

	// Get latest version records
	recordIds := []string{}
	serviceToVersionNumMap := make(map[string]int64)

	for _, v := range latestServices {
		recordIds = append(recordIds, v.Id)
		serviceToVersionNumMap[v.Id] = v.Versions
	}
	serviceRecords := sm.ServiceDao.GetServiceRecordsByRecordIds(recordIds)

	// Return Services

	servicesResponse := []*models.Service{}

	for _, v := range serviceRecords {
		servicesResponse = append(servicesResponse, &models.Service{
			Id:          v.Id,
			ServiceId:   *v.ServiceId,
			Name:        v.Name,
			Description: v.Description,
			Version:     v.Version,
			Versions:    serviceToVersionNumMap[v.Id],
		})

	}

	return servicesResponse
}

// Fetch Service
func (sm ServiceManager) GetService(serviceId int64) *models.Service {
	latestService := sm.ServiceLatestDao.GetService(serviceId)
	serviceRecord := sm.ServiceDao.getServiceRecordByRecordId(latestService.Id)

	return &models.Service{
		Id:          serviceRecord.Id,
		ServiceId:   *serviceRecord.ServiceId,
		Name:        serviceRecord.Name,
		Description: serviceRecord.Description,
		Version:     serviceRecord.Version,
		Versions:    latestService.Versions,
	}

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
		return sm.ServiceDao.CreateNewService(newServiceRecord), nil
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

		return sm.ServiceDao.CreateNewServiceVersion(newServiceRecord), nil
	}

}
