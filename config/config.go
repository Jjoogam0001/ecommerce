// Package config allows for customisation of the risk API
package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"os"
	"path"
	"path/filepath"
	"strings"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/labstack/gommon/log"
)

// AppConfig contains the application configuration.
type AppConfig struct {
	Database   DatabaseConfig
	HTTPServer HTTPServer
	LogLevel   string
	Swagger    swagger
}

// DatabaseConfig provides the connection string and other params for configuring the database.
type DatabaseConfig struct {
	ConnectionString string
	LocalDocker      DockerConfig
}

// DockerConfig is the docker configuration mapping.
type DockerConfig struct {
	Enabled          bool
	Image            string
	Name             string
	RemoveOnFinish   bool
	IgnoreOnConflict bool
}
type swagger struct {
	Host string
}
type Feed struct {
	Name string
	Path string
}

// HTTPServer defines the server http configuration.
type HTTPServer struct {
	Addr           string
	RequestTimeout int
}

// HostConfig is the host basic configuration.
type HostConfig struct {
	env string
}

func (c HostConfig) String() string {
	return fmt.Sprintf("host config env: %s", c.env)
}

// BuildHost generates the host config from the flags and configures the log level.
func BuildHost() HostConfig {
	var host HostConfig

	loadFlags(func(s *flag.FlagSet) {
		s.StringVar(&host.env, "env", "local", "Environment")
	})

	fmt.Fprintf(os.Stderr, "%s\n", host)

	return host
}

// OverridingEnv overrides the env value.
func (c HostConfig) OverridingEnv(env string) HostConfig {
	c.env = env
	return c
}

func loadFlags(flags func(*flag.FlagSet)) {
	set := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	set.SetOutput(io.Discard)
	flags(set)

	err := set.Parse(os.Args[1:])
	set.SetOutput(nil)

	if errors.Is(err, flag.ErrHelp) {
		set.Usage()
		os.Exit(0)
	}
}

// ReadAppConfig loads the configuration defined in config.*.json files based on the env value of host config.
func (c HostConfig) ReadAppConfig() (AppConfig, error) {
	var a AppConfig

	env := strings.ToLower(c.env)
	log.Infof("Loading configuration for environment [%v]", env)

	cwd, _ := os.Getwd()
	dirPath, _ := findPath("config", cwd)

	configFile := "config.json"
	config := path.Join(dirPath, configFile)
	if err := a.loadJSON(config); err != nil {
		return a, fmt.Errorf("cannot parse default config file from path %s for env %s in folder (and parents) %s", config, env, cwd)
	}

	if env != "" {
		configFile = "config." + env + ".json"
		config := path.Join(dirPath, configFile)
		if err := a.loadJSON(config); err != nil {
			return a, fmt.Errorf("cannot parse config file from path %s for env %s in folder (and parents) %s", config, env, cwd)
		}
	}

	log.SetLevel(map[string]log.Lvl{
		"DEBUG": log.DEBUG,
		"INFO":  log.INFO,
		"WARN":  log.WARN,
		"ERROR": log.ERROR,
		"OFF":   log.OFF,
	}[strings.ToUpper(a.LogLevel)])

	return a, nil
}

func (c *AppConfig) loadJSON(filePath string) error {
	log.Infof("Loading file %s", filePath)

	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("Ignoring config file %v: %w", filePath, err)
	}

	file, err := os.Open(filePath) // nolint:gosec
	if err != nil {
		return fmt.Errorf("Error opening file %v: %w", filePath, err)
	}

	byteValue, _ := io.ReadAll(file)

	if err = file.Close(); err != nil {
		return fmt.Errorf("Error closing file %v: %w", filePath, err)
	}

	if err = json.Unmarshal(byteValue, c); err != nil {
		return fmt.Errorf("Error unmarshalling file %v: %w", filePath, err)
	}

	log.Infof("Success loading config: %v", c)

	return nil
}

// StartDatabase connects to the database configured in the connection string and migrates it to the latest version.
func (c AppConfig) StartDatabase(ctx context.Context) (*pgxpool.Pool, error) {

	db, err := sql.Open("pgx", c.Database.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting", err)
	}
	defer db.Close()

	return pgxpool.Connect(ctx, c.Database.ConnectionString)
}

func findPath(target string, dir string) (string, error) {

	p := path.Join(dir, target)

	_, err := os.Stat(p)

	if err == nil {
		return p, nil
	}

	if dir == "/" {
		return "", os.ErrNotExist
	}

	return findPath(target, filepath.Dir(dir))
}
