package util

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
)

const (
	Namespace = "nextdns"
	BaseURL   = "https://api.nextdns.io"
)

var (
	Log          hclog.Logger
	Port         string
	MetricsPath  string
	ResultWindow string
	ResultLimit  string
	Profile      string
	APIKey       string
)

// Initialize the configuration.
func init() {
	// Set up logging.
	level := DefaultEnv("LOG_LEVEL", "INFO")
	Log = hclog.New(&hclog.LoggerOptions{
		Level: hclog.LevelFromString(level),
	})

	// Retrieve configuration, or use defaults.
	Port = fmt.Sprintf(":%s", DefaultEnv("METRICS_PORT", "9948"))
	MetricsPath = DefaultEnv("METRICS_PATH", "/metrics")
	ResultWindow = DefaultEnv("NEXTDNS_RESULT_WINDOW", "-5m")
	ResultLimit = DefaultEnv("NEXTDNS_RESULT_LIMIT", "50")

	// Required configuration.
	var set bool
	Profile, set = os.LookupEnv("NEXTDNS_PROFILE")
	if !set {
		Log.Error("NEXTDNS_PROFILE must be set")
		os.Exit(1)
	}
	APIKey, set = os.LookupEnv("NEXTDNS_API_KEY")
	if !set {
		Log.Error("NEXTDNS_API_KEY must be set")
		os.Exit(1)
	}
}
