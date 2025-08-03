package mock_services

import (
	"static-site-hosting/models"
	"static-site-hosting/services"
)

type fakeDB struct{}

// DeleteSiteById implements services.DatabaseService.
func (f *fakeDB) DeleteSiteById(id int) error {
	return nil
}

// DeleteSiteByName implements services.DatabaseService.
func (f *fakeDB) DeleteSiteByName(name string) error {
	return nil
}

// GetDeploymentLogs implements services.DatabaseService.
func (f *fakeDB) GetDeploymentLogs(siteName string) ([]models.DeploymentLog, error) {
	return []models.DeploymentLog{
		{SiteName: siteName, IPAddress: "127.0.0.1", UserAgent: "TestAgent", Timestamp: "2023-10-01T12:00:00Z"},
	}, nil
}

// Initialize implements services.DatabaseService.
func (f *fakeDB) Initialize() error {
	return nil
}

// InsertDeploymentLog implements services.DatabaseService.
func (f *fakeDB) InsertDeploymentLog(siteName string, ipAddress string, userAgent string) error {
	return nil
}

// InsertSite implements services.DatabaseService.
func (f *fakeDB) InsertSite(name string) error {
	return nil
}

// ListSites implements services.DatabaseService.
func (f *fakeDB) ListSites() ([]models.Site, error) {
	sites := []models.Site{
		{Name: "site1", DeployedAt: "2023-10-01T12:00:00Z"},
		{Name: "site2", DeployedAt: "2023-10-02T12:00:00Z"},
		{Name: "site3", DeployedAt: "2023-10-03T12:00:00Z"},
		{Name: "site4", DeployedAt: "2023-10-04T12:00:00Z"},
	}

	return sites, nil
}

func NewFakeDB() services.DatabaseService {
	return &fakeDB{}
}
