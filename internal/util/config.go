package util

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
)

const (
	Namespace = "nextdns"
	BaseURL   = "https://api.nextdns.io"
)

var (
	Log         hclog.Logger
	Port        string
	MetricsPath string
	Profile     string
	APIKey      string
	FilterFrom  string
	ResultLimit string
)

// Initialize the configuration.
func init() {
	// Set up logging.
	level := GetEnv("LOG_LEVEL", "INFO")
	Log = hclog.New(&hclog.LoggerOptions{
		Level: hclog.LevelFromString(level),
	})

	// Set up exporter.
	Port = fmt.Sprintf(":%s", GetEnv("METRICS_PORT", "9948"))
	MetricsPath = GetEnv("METRICS_PATH", "/metrics")
	Profile = GetEnv("NEXTDNS_PROFILE", "")
	APIKey = GetEnv("NEXTDNS_API_KEY", "")
	FilterFrom = GetEnv("NEXTDNS_FILTER_FROM", "-1d")
	ResultLimit = GetEnv("NEXTDNS_RESULT_LIMIT", "50")
}
