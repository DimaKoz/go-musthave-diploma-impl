package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	flag2 "github.com/spf13/pflag"
)

const (
	EnvKeyAddress     = "RUN_ADDRESS"
	EnvKeyDatabaseURI = "DATABASE_URI"
)

// Config represents a config of the server.
type Config struct {
	Address      string `env:"RUN_ADDRESS"`
	ConnectionDB string `env:"DATABASE_URI"`
	Accrual      string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// NewConfig creates an instance of Config.
func NewConfig() *Config {
	return &Config{Address: "", ConnectionDB: "", Accrual: ""}
}

// ProcessEnv receives and sets up the Config.
type ProcessEnv func(config *Config) error

// LoadConfig loads data to the passed Config.
func LoadConfig(cfg *Config, processing ProcessEnv) error {
	if err := processing(cfg); err != nil {
		return fmt.Errorf("config: cannot process ENV variables: %w", err)
	}

	processFlags(cfg)

	return nil
}

func processFlags(cfg *Config) {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	aFlag, dFlag, rFlag := "", "", ""

	addFlags(cfg, &aFlag, &rFlag, &dFlag)

	flag2.Parse()

	cfg.Address = aFlag
	cfg.ConnectionDB = dFlag
	cfg.Accrual = rFlag
}

func addFlags(cfg *Config, aFlag *string, rFlag *string, dFlag *string) {
	flag2.StringVarP(aFlag, "a", "a", cfg.Address, "")
	flag2.StringVarP(rFlag, "r", "r", cfg.Accrual, "")
	flag2.StringVarP(dFlag, "d", "d", cfg.ConnectionDB, "")
}

func ProcessEnvServer(config *Config) error {
	log.Println(os.Environ())

	opts := env.Options{ //nolint:exhaustruct
		OnSet: func(tag string, value interface{}, isDefault bool) {
			log.Printf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}

	if err := env.Parse(config, opts); err != nil {
		return fmt.Errorf("failed to parse an environment, error: %w", err)
	}

	return nil
}

func (cfg Config) String() string {
	return fmt.Sprintf("Address:[%s] \n ConnectionDB:[%s] \n Accrual:[%s] \n",
		cfg.Address, cfg.ConnectionDB, cfg.Accrual)
}
