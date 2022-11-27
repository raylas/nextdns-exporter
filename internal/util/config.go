package util

import (
	"fmt"
	"runtime/debug"

	"github.com/hashicorp/go-hclog"
)

const (
	Namespace = "nextdns"
	BaseURL   = "https://api.nextdns.io"
)

var (
	Log          hclog.Logger
	Version      string
	Port         string
	MetricsPath  string
	Profile      string
	APIKey       string
	ResultWindow string
	ResultLimit  string
)

// Initialize the configuration.
func init() {
	// Set up logging.
	level := GetEnv("LOG_LEVEL", "INFO")
	Log = hclog.New(&hclog.LoggerOptions{
		Level: hclog.LevelFromString(level),
	})

	// Retrieve version.
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		Version = info.Main.Version
	} else {
		Version = "dev"
	}

	// Set up exporter.
	Port = fmt.Sprintf(":%s", GetEnv("METRICS_PORT", "9948"))
	MetricsPath = GetEnv("METRICS_PATH", "/metrics")
	Profile = GetEnv("NEXTDNS_PROFILE", "")
	APIKey = GetEnv("NEXTDNS_API_KEY", "")
	ResultWindow = GetEnv("NEXTDNS_RESULT_WINDOW", "-5m")
	ResultLimit = GetEnv("NEXTDNS_RESULT_LIMIT", "50")
}
