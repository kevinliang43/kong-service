package service_module

import (
	"database/sql"
	"dev/kong-service/models"
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

func (sld ServiceLatestDao) GetAllServices() []*models.ServiceLatest {

	query := `SELECT service_id, latest_record_id, versions FROM services_latest;`

	rows, err := sld.database.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.ServiceLatest{}
		}
		log.Fatal(err)
	}
	defer rows.Close()
	return parseServiceLatestRows(rows)
}

func (sld ServiceLatestDao) GetService(sid int64) *models.ServiceLatest {
	var (
		serviceId      int64
		latestRecordId string
		versions       int64
	)

	err := sld.database.QueryRow("SELECT service_id, latest_record_id, versions FROM services_latest WHERE service_id = $1", sid).Scan(&serviceId, &latestRecordId, &versions)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}

	return &models.ServiceLatest{
		ServiceId: &serviceId,
		Id:        latestRecordId,
		Versions:  versions,
	}
}

func (sld ServiceLatestDao) CreateService(sr *models.ServiceRecord) *models.ServiceLatest {
	var (
		serviceId      int64
		latestRecordId string
		versions       int64
	)

	err := sld.database.QueryRow("INSERT INTO services_latest(service_id, latest_record_id, versions) VALUES($1, $2, $3) RETURNING service_id, latest_record_id, versions", sr.ServiceId, sr.Id, 1).Scan(&serviceId, &latestRecordId, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return &models.ServiceLatest{
		ServiceId: &serviceId,
		Id:        latestRecordId,
		Versions:  versions,
	}

}

func (sld ServiceLatestDao) UpdateService(sr *models.ServiceRecord) *models.ServiceLatest {

	var (
		serviceId      int64
		latestRecordId string
		versions       int64
	)

	existingService := sld.GetService(*sr.ServiceId)

	if existingService == nil {
		return nil
	}

	updateQuery := `UPDATE services_latest SET latest_record_id=$1, versions=$2 WHERE service_id=$3 RETURNING service_id, latest_record_id, versions;`

	err := sld.database.QueryRow(updateQuery, sr.Id, existingService.Versions+1, sr.ServiceId).Scan(&serviceId, &latestRecordId, &versions)
	if err != nil {
		log.Fatal(err)
	}

	return &models.ServiceLatest{
		ServiceId: &serviceId,
		Id:        latestRecordId,
		Versions:  versions,
	}

}

func parseServiceLatestRows(rows *sql.Rows) []*models.ServiceLatest {
	services := []*models.ServiceLatest{}

	for rows.Next() {
		var (
			serviceId      int64
			latestRecordId string
			versions       int64
		)

		rows.Scan(&serviceId, &latestRecordId, &versions)

		services = append(services, &models.ServiceLatest{
			ServiceId: &serviceId,
			Id:        latestRecordId,
			Versions:  versions,
		})
	}

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return services
}
