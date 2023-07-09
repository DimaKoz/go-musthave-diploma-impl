package gophermart

import (
	"io"
	"os"
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/config"
	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/util"
	flag2 "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupConfigOk(t *testing.T) {
	want := &config.Config{
		Address:      ":8080",
		ConnectionDB: "ej",
		Accrual:      "",
	}

	origValueAddress := os.Getenv(config.EnvKeyAddress)
	err := os.Setenv(config.EnvKeyAddress, ":8080") //nolint:tenv
	require.NoError(t, err)
	origValuePathDB := os.Getenv(config.EnvKeyDatabaseURI)
	err = os.Setenv(config.EnvKeyDatabaseURI, "ej") //nolint:tenv
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Setenv(config.EnvKeyAddress, origValueAddress)
		require.NoError(t, err)
		err = os.Setenv(config.EnvKeyDatabaseURI, origValuePathDB)
		require.NoError(t, err)
	})
	got := config.NewConfig()
	err = setupConfig(got, config.ProcessEnvServer)
	assert.NoError(t, err)
	assert.Equal(t, want, got, "setupConfig() = %v, want %v", got, want)
}

func TestSetupConfigErr(t *testing.T) {
	processEnv := func(config *config.Config) error {
		return os.ErrPermission // an any err
	}
	gotErr := setupConfig(nil, processEnv)
	assert.Error(t, gotErr)
}

func TestSetupConfigEmptyAddress(t *testing.T) {
	osArgOrig := os.Args
	flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
	flag2.CommandLine.SetOutput(io.Discard)
	os.Args = make([]string, 0)
	os.Args = append(os.Args, osArgOrig[0])
	t.Cleanup(func() { os.Args = osArgOrig })

	cfg := config.NewConfig()
	gotErr := setupConfig(cfg, config.ProcessEnvServer)
	assert.Error(t, gotErr)
}

func TestRunEmptyAddress(t *testing.T) {
	osArgOrig := os.Args
	flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
	flag2.CommandLine.SetOutput(io.Discard)
	os.Args = make([]string, 0)
	os.Args = append(os.Args, osArgOrig[0])
	t.Cleanup(func() { os.Args = osArgOrig })

	output := util.CaptureOutput(func() {
		Run()
	})
	assert.Contains(t, output, "server address is empty")
	t.Log("log:", output)
}

func TestRunEmptyPathDB(t *testing.T) {
	osArgOrig := os.Args
	flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
	flag2.CommandLine.SetOutput(io.Discard)
	os.Args = make([]string, 0)
	os.Args = append(os.Args, osArgOrig[0])
	os.Args = append(os.Args, "-a")
	os.Args = append(os.Args, ":8080")
	t.Cleanup(func() { os.Args = osArgOrig })

	output := util.CaptureOutput(func() {
		Run()
	})
	assert.Contains(t, output, "db uri is empty")
	t.Log("log:", output)
}
