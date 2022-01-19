package service_module

import (
	"database/sql"
	"dev/kong-service/models"
	"fmt"
	"log"
	"strings"
)

type ServiceDao struct {
	database *sql.DB
}

func NewServiceDao(database *sql.DB) *ServiceDao {
	return &ServiceDao{
		database: database,
	}
}

func (sd ServiceDao) GetAllServiceRecordsByServiceId(sid int64) []*models.ServiceRecord {
	query := `SELECT id, service_id, name, description, version FROM services WHERE service_id=$1;`

	rows, err := sd.database.Query(query, sid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.ServiceRecord{}
		}
		log.Fatal(err)
	}
	defer rows.Close()
	return parseServiceRows(rows)
}

func (sd ServiceDao) getServiceRecordByServiceIdAndVersion(sid int64, v float64) *models.ServiceRecord {
	var (
		id          string
		serviceId   int64
		name        string
		description string
		version     float64
	)

	err := sd.database.QueryRow("SELECT id, service_id, name, description, version FROM services WHERE service_id = $1 AND version = $2", sid, v).Scan(&id, &serviceId, &name, &description, &version)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return &models.ServiceRecord{
		Id:          id,
		ServiceId:   &serviceId,
		Name:        name,
		Description: description,
		Version:     version,
	}
}

func (sd ServiceDao) GetAllServiceVersionsByServiceId(sid int64) []float64 {
	query := `SELECT version FROM services WHERE service_id=$1;`
	rows, err := sd.database.Query(query, sid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []float64{}
		}
		log.Fatal(err)
	}
	defer rows.Close()

	versions := []float64{}

	for rows.Next() {
		var version float64
		rows.Scan(&version)
		versions = append(versions, version)
	}
	return versions

}

func (sd ServiceDao) GetServiceRecordsByRecordIds(serviceRecordIds []string) []*models.ServiceRecord {
	query := `SELECT id, service_id, name, description, version FROM services WHERE id = ANY($1);`
	param := "{" + strings.Join(serviceRecordIds, ",") + "}"

	rows, err := sd.database.Query(query, param)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.ServiceRecord{}
		}
		log.Fatal(err)
	}
	defer rows.Close()
	return parseServiceRows(rows)
}

func (sd ServiceDao) getServiceRecordByRecordId(serviceRecordId string) *models.ServiceRecord {
	var (
		id          string
		serviceId   int64
		name        string
		description string
		version     float64
	)

	err := sd.database.QueryRow("SELECT id, service_id, name, description, version FROM services WHERE id = $1", serviceRecordId).Scan(&id, &serviceId, &name, &description, &version)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return &models.ServiceRecord{
		Id:          id,
		ServiceId:   &serviceId,
		Name:        name,
		Description: description,
		Version:     version,
	}

}

func (sd ServiceDao) CreateNewService(s *models.ServiceRecord) *models.ServiceRecord {
	var id string

	err := sd.database.QueryRow("INSERT INTO services(name, description, version) VALUES($1, $2, $3) RETURNING id", s.Name, s.Description, s.Version).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	newService := sd.getServiceRecordByRecordId(id)

	fmt.Printf("Created new service with ID=%s\n", id)

	return newService

}

func (sd ServiceDao) CreateNewServiceVersion(s *models.ServiceRecord) *models.ServiceRecord {

	var (
		id        string
		serviceId int64
		version   float64
	)

	err := sd.database.QueryRow("INSERT INTO services(service_id, name, description, version) VALUES($1, $2, $3, $4) RETURNING id, service_id, version", s.ServiceId, s.Name, s.Description, s.Version).Scan(&id, &serviceId, &version)
	if err != nil {
		log.Fatal(err)
	}

	newService := sd.getServiceRecordByRecordId(id)

	fmt.Printf("Created new service version with ID=%s, ServiceId=%d Version=%f\n", id, serviceId, version)

	return newService
}

func parseServiceRows(rows *sql.Rows) []*models.ServiceRecord {
	services := []*models.ServiceRecord{}

	for rows.Next() {
		var (
			id          string
			serviceId   int64
			name        string
			description string
			version     float64
		)

		rows.Scan(&id, &serviceId, &name, &description, &version)

		services = append(services, &models.ServiceRecord{
			Id:          id,
			ServiceId:   &serviceId,
			Name:        name,
			Description: description,
			Version:     version,
		})
	}

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return services
}
