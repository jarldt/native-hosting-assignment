package models

// DeploymentLog represents one deployment event
type DeploymentLog struct {
	SiteName  string `json:"site_name"`
	Timestamp string `json:"timestamp"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

// Site represents a deployed site
type Site struct {
	Name       string `json:"name"`
	DeployedAt string `json:"deployed_at"`
}
