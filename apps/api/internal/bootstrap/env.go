package bootstrap

import "github.com/aisha-platform/aisha/apps/api/internal/config"

// LoadConfig is the only process-start configuration boundary.
func LoadConfig() (config.Config, error) {
	return config.Load()
}
