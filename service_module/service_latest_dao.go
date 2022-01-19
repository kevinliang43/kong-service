package service_module

import (
	"database/sql"
	"dev/kong-service/models"
	"fmt"
	"log"
)

type ServiceLatestDao struct {
	database *sql.DB
}

func NewServiceLatestDao(database *sql.DB) *ServiceLatestDao {
	return &ServiceLatestDao{
		database: database,
	}
}

func (sld ServiceLatestDao) GetAllServices() []*models.Service {

	query := `SELECT service_id, latest_record_id, name, description, version, versions FROM services_latest;`

	rows, err := sld.database.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Service{}
		}
		log.Fatal(err)
	}
	defer rows.Close()
	return parseServiceLatestRows(rows)
}

func (sld ServiceLatestDao) SearchServices(ssr *models.ServicesSearchRequest) []*models.Service {
	baseQuery := `SELECT service_id, latest_record_id, name, description, version, versions FROM services_latest `

	if ssr.NameFilter != nil {
		baseQuery = baseQuery + fmt.Sprintf("WHERE name LIKE '%%%s%%' ", *ssr.NameFilter)
	}

	if ssr.Limit != nil {
		baseQuery = baseQuery + fmt.Sprintf("LIMIT %d ", *ssr.Limit)
	}

	if ssr.Offset != nil {
		baseQuery = baseQuery + fmt.Sprintf("OFFSET %d ", *ssr.Offset)
	}

	baseQuery = baseQuery + ";"

	rows, err := sld.database.Query(baseQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Service{}
		}
		log.Fatal(err)
	}
	defer rows.Close()

	return parseServiceLatestRows(rows)

}

func (sld ServiceLatestDao) GetService(sid int64) *models.Service {
	var (
		serviceId      int64
		latestRecordId string
		name           string
		description    string
		version        float64
		versions       int64
	)

	err := sld.database.QueryRow("SELECT service_id, latest_record_id, name, description, version, versions FROM services_latest WHERE service_id = $1", sid).Scan(&serviceId, &latestRecordId, &name, &description, &version, &versions)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return &models.Service{
		ServiceId:   serviceId,
		Id:          latestRecordId,
		Name:        name,
		Description: description,
		Version:     version,
		Versions:    versions,
	}
}

func (sld ServiceLatestDao) CreateService(sr *models.ServiceRecord) *models.Service {
	var (
		serviceId      int64
		latestRecordId string
		name           string
		description    string
		version        float64
		versions       int64
	)

	err := sld.database.QueryRow("INSERT INTO services_latest(service_id, latest_record_id, name, description, version, versions) VALUES($1, $2, $3) RETURNING service_id, latest_record_id, versions", sr.ServiceId, sr.Id, sr.Name, sr.Description, sr.Version, 1).Scan(&serviceId, &latestRecordId, &name, &description, &version, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return &models.Service{
		Id:          latestRecordId,
		ServiceId:   serviceId,
		Name:        name,
		Description: description,
		Version:     version,
		Versions:    versions,
	}

}

func (sld ServiceLatestDao) UpdateService(sr *models.ServiceRecord) *models.Service {

	var (
		serviceId      int64
		latestRecordId string
		name           string
		description    string
		version        float64
		versions       int64
	)

	existingService := sld.GetService(*sr.ServiceId)

	if existingService == nil {
		return nil
	}

	updateQuery := `UPDATE services_latest SET latest_record_id=$1, name=$2, description=$3, version=$4, versions=$5 WHERE service_id=$6 RETURNING service_id, latest_record_id, name, description, version, versions;`

	err := sld.database.QueryRow(updateQuery, sr.Id, sr.Name, sr.Description, sr.Version, existingService.Versions+1, sr.ServiceId).Scan(&serviceId, &latestRecordId, &name, &description, &version, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return &models.Service{
		Id:          latestRecordId,
		ServiceId:   serviceId,
		Name:        name,
		Description: description,
		Version:     version,
		Versions:    versions,
	}

}

func parseServiceLatestRows(rows *sql.Rows) []*models.Service {
	services := []*models.Service{}

	for rows.Next() {
		var (
			serviceId      int64
			latestRecordId string
			name           string
			description    string
			version        float64
			versions       int64
		)

		rows.Scan(&serviceId, &latestRecordId, &name, &description, &version, &versions)

		services = append(services, &models.Service{
			ServiceId:   serviceId,
			Id:          latestRecordId,
			Name:        name,
			Description: description,
			Version:     version,
			Versions:    versions,
		})
	}

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return services
}
