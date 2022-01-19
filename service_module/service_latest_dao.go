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
