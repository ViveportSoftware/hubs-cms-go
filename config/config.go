package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type envVariable struct {
	Version               string `env:"FULL_VERSION" envDefault:"1.0.0"`
	Port                  string `env:"GO_HTTP_PORT,required"`
	LogLevel              string `env:"LOG_LEVEL" envDefault:"ERROR"`
	Environment           string `env:"ENVIRONMENT" envDefault:"DEVELOP"`
	MastodonBaseURI       string `env:"MASTODON_BASE_URI,required"`
	DirectusBaseURI       string `env:"DIRECTUS_BASE_URI,required"`
	DirectusAdminEmail    string `env:"DIRECTUS_ADMIN_EMAIL,required"`
	DirectusAdminPassword string `env:"DIRECTUS_ADMIN_PASSWORD,required"`
	HubsBaseURI           string `env:"HUBS_BASE_URI,required"`
	EventBackupInterval   string `env:"EVENT_BACKUP_INTERVAL" envDefault:"@daily"`
}

func (r envVariable) Validate() bool {

	port, err := strconv.ParseInt(EnvVariable.Port, 10, 16)
	if err != nil || port <= 0 || port > 65535 {
		log.Fatalf("ERR: required environment variable \"GO_HTTP_PORT\" should be 0~65535")
		return false
	}

	if d := strings.ToLower(EnvVariable.LogLevel); d != "error" && d != "debug" && d != "warn" && d != "info" {
		log.Fatalf("ERR: required environment variable \"LOG_LEVEL\" should be \"ERROR|DEBUG|WARN|INFO\"")
		return false
	}

	if d := strings.ToLower(EnvVariable.Environment); d != "develop" && d != "stage" && d != "production" {
		log.Fatalf("ERR: required environment variable \"ENVIRONMENT\" should be \"DEVELOP|STAGE|PRODUCTION\"")
		return false
	}

	if EnvVariable.MastodonBaseURI == "" {
		log.Fatalf("ERR: required environment variable \"MASTODON_BASE_URI\" is not set")
		return false
	}

	mastodonUri, err := url.Parse(EnvVariable.MastodonBaseURI)
	if err != nil {
		log.Fatalf("ERR: required environment variable \"MASTODON_BASE_URI\" is invalid")
		return false
	}
	DefaultMastodonAccountDomain = fmt.Sprintf("@%s", mastodonUri.Host)

	if EnvVariable.DirectusBaseURI == "" {
		log.Fatalf("ERR: required environment variable \"DIRECTUS_BASE_URI\" is not set")
		return false
	}

	if EnvVariable.DirectusAdminEmail == "" {
		log.Fatalf("ERR: required environment variable \"DIRECTUS_ADMIN_EMAIL\" is not set")
		return false
	}

	if EnvVariable.DirectusAdminPassword == "" {
		log.Fatalf("ERR: required environment variable \"DIRECTUS_ADMIN_PASSWORD\" is not set")
		return false
	}

	return true
}

// EnvVariable is the environment variables config
var EnvVariable envVariable
var DefaultMastodonAccountDomain string

// Setup loads and validates config from environment variables
func Setup() {

	// parse env
	if err := env.Parse(&EnvVariable); err != nil {
		log.Fatalf("ERR: %v\n", err)
	}

	// validate env
	if !EnvVariable.Validate() {
		os.Exit(1)
	}
}

// IsDevEnv should enable gin debug mode for more logs
func IsDevEnv() bool {
	d := strings.ToLower(EnvVariable.Environment)
	return d == "develop"
}
