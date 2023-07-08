package config_test

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	flag2 "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	want := &config.Config{Address: "", ConnectionDB: "", Accrual: ""}
	got := config.NewConfig()
	assert.Equal(t, want, got)
}

func TestProcessEnvNoError(t *testing.T) {
	var wantErr error
	gotErr := config.ProcessEnvServer(config.NewConfig())

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)
}

var errTestProcessEnvError = errors.New("env: expected a pointer to a Struct")

func TestProcessEnvError(t *testing.T) {
	wantErr := fmt.Errorf("failed to parse an environment, error: %w", errTestProcessEnvError)
	gotErr := config.ProcessEnvServer(nil)

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)
}

func TestLoadConfigEnvError(t *testing.T) {
	gotErr := config.LoadConfig(nil, config.ProcessEnvServer)

	assert.Error(t, gotErr)
}

func TestConfigStringer(t *testing.T) {
	want := "Address: b \n ConnectionDB: c \n Accrual: a \n"
	cfg := config.NewConfig()
	cfg.Accrual = "a"
	cfg.Address = "b"
	cfg.ConnectionDB = "c"
	got := cfg.String()

	assert.Equal(t, want, got)
}

func TestLoadEmptyConfigNoErr(t *testing.T) {
	want := &config.Config{
		Address:      "",
		ConnectionDB: "",
		Accrual:      "",
	}
	got := config.NewConfig()
	err := config.LoadConfig(got, config.ProcessEnvServer)
	assert.NoError(t, err, "error must be nil")
	assert.Equal(t, want, got, "Configs - got: %v, want: %v", got, want)
}

type argTestConfig struct {
	envAddress  string
	envDB       string
	envAccrual  string
	flagAddress string
	flagDB      string
	flagAccrual string
}

var wantConfigEnv = &config.Config{
	Address: "127.0.0.1:59483", ConnectionDB: "db_uri", Accrual: "accrual",
}

var testsCasesInitConfig = []struct {
	name    string
	args    argTestConfig
	want    *config.Config
	wantErr error
}{
	{
		name: "env", args: argTestConfig{ //nolint:exhaustruct
			envAddress: "127.0.0.1:59483",
			envDB:      "db_uri",
			envAccrual: "accrual",
		},
		want: wantConfigEnv,
	},
	{
		name: "flags",
		args: argTestConfig{ //nolint:exhaustruct
			flagAddress: "127.0.0.1:59483",
			flagDB:      "db_uri",
			flagAccrual: "accrual",
		},
		want: wantConfigEnv,
	},
	{
		name: "flags&env, priority", want: wantConfigEnv,
		args: argTestConfig{
			envAddress: "wrong_address", envDB: "wrong_uri", envAccrual: "",
			flagAddress: "127.0.0.1:59483",
			flagDB:      "db_uri",
			flagAccrual: "accrual",
		},
	},
}

func TestAgentInitConfig(t *testing.T) {
	for _, test := range testsCasesInitConfig {
		test := test
		t.Run(test.name, func(t *testing.T) {
			envArgsInitConfig(t, config.EnvKeyAddress, test.args.envAddress) // ENV setup
			envArgsInitConfig(t, "DATABASE_URI", test.args.envDB)
			envArgsInitConfig(t, "ACCRUAL_SYSTEM_ADDRESS", test.args.envAccrual)
			osArgOrig := os.Args
			flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
			flag2.CommandLine.SetOutput(io.Discard)
			os.Args = make([]string, 0)
			os.Args = append(os.Args, osArgOrig[0])
			appendArgsInitConfig(t, &os.Args, "-a", test.args.flagAddress)
			appendArgsInitConfig(t, &os.Args, "-r", test.args.flagAccrual)
			appendArgsInitConfig(t, &os.Args, "-d", test.args.flagDB)
			t.Cleanup(func() { os.Args = osArgOrig })

			got := config.NewConfig()
			gotErr := config.LoadConfig(got, config.ProcessEnvServer)

			if test.wantErr != nil {
				assert.EqualErrorf(t, gotErr, test.wantErr.Error(), "Configs - got error: %v, want: %v", gotErr, test.wantErr)
			} else {
				assert.NoError(t, gotErr, "Configs - got error: %v, want: %v", gotErr, test.wantErr)
			}

			assert.Equal(t, test.want, got, "Configs - got: %v, want: %v", got, test.want)
		})
	}
}

func appendArgsInitConfig(t *testing.T, target *[]string, key string, value string) {
	t.Helper()
	if value != "" {
		*target = append(*target, key)
		*target = append(*target, value)
	}
}

func envArgsInitConfig(t *testing.T, key string, value string) {
	t.Helper()
	if value != "" {
		origValue := os.Getenv(key)
		err := os.Setenv(key, value)
		log.Println("new "+key+":", value, " err:", err)
		assert.NoError(t, err)
		t.Cleanup(func() { _ = os.Setenv(key, origValue) })
	}
}
