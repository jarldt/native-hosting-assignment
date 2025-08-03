package services

import (
	"database/sql"
	"fmt"
	"static-site-hosting/models"
)

type DatabaseService interface {
	Initialize() error
	InsertSite(name string) error
	InsertDeploymentLog(siteName, ipAddress, userAgent string) error
	GetDeploymentLogs(siteName string) ([]models.DeploymentLog, error)
	ListSites() ([]models.Site, error)
	DeleteSiteById(id int) error
	DeleteSiteByName(name string) error
}

type service struct {
	db *sql.DB
}

func (s *service) Initialize() error {
	// Create the sites table
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS sites (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	deployed_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return fmt.Errorf("error creating sites table: %v", err)
	}

	// Create the deployment_logs table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS deployment_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site_name TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		ip_address TEXT,
		user_agent TEXT)`)

	if err != nil {
		return fmt.Errorf("error creating deployment_logs table: %v", err)
	}

	return nil
}

func (s *service) InsertSite(name string) error {
	_, err := s.db.Exec("INSERT INTO sites (name) VALUES (?)", name)
	if err != nil {
		return fmt.Errorf("error inserting site: %v", err)
	}
	return nil
}

func (s *service) InsertDeploymentLog(siteName, ipAddress, userAgent string) error {
	_, err := s.db.Exec("INSERT INTO deployment_logs (site_name, ip_address, user_agent) VALUES (?, ?, ?)",
		siteName, ipAddress, userAgent)
	if err != nil {
		return fmt.Errorf("error inserting deployment log: %v", err)
	}
	return nil
}

func (s *service) GetDeploymentLogs(siteName string) ([]models.DeploymentLog, error) {
	rows, err := s.db.Query(`
			SELECT site_name, timestamp, ip_address, user_agent
			FROM deployment_logs
			WHERE site_name = ?
			ORDER BY timestamp DESC LIMIT 100`, siteName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.DeploymentLog
	for rows.Next() {
		var log models.DeploymentLog
		if err := rows.Scan(&log.SiteName, &log.Timestamp, &log.IPAddress, &log.UserAgent); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *service) ListSites() ([]models.Site, error) {
	rows, err := s.db.Query("SELECT name, deployed_at FROM sites ORDER BY deployed_at DESC")
	if err != nil {
		return nil, fmt.Errorf("error querying sites: %v", err)
	}
	defer rows.Close()

	var sites []models.Site
	for rows.Next() {
		var s models.Site
		err := rows.Scan(&s.Name, &s.DeployedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning site row: %v", err)
		}
		sites = append(sites, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over site rows: %v", err)
	}

	return sites, nil
}

func (s *service) DeleteSiteById(id int) error {
	_, err := s.db.Exec("DELETE FROM sites WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting site by id: %v", err)
	}
	return nil
}

func (s *service) DeleteSiteByName(name string) error {
	_, err := s.db.Exec("DELETE FROM sites WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("error deleting site by name: %v", err)
	}
	return nil
}

func NewDatabaseService(db *sql.DB) DatabaseService {
	return &service{db: db}
}
