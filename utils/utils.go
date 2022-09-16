package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ParsedRequestData map[string]string

func ParseFormData(reqBody io.ReadCloser) (ParsedRequestData, error) {
	data, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read request.body: ")
	}

	var reqData ParsedRequestData
	if err := json.Unmarshal(data, &reqData); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal data: ")
	}

	return reqData, nil
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	DBname   string
	User     string
	Password string
	TimeZone string
	SSLMODE  string
	DBport   string
}

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	SecretKey string
}

func (c *Config) SetUpViper(configFile, path, format string) error {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(path)
	// allows env variabless to be read
	viper.AutomaticEnv()
	viper.SetConfigType(format)

	if err := viper.ReadInConfig(); err != nil {
		m := fmt.Sprintf("SetUpViper unable to read %v.%v", configFile, format)
		return errors.Wrap(err, m)
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return errors.Wrap(err, "Unable to unmarshal env variables")
	}

	if err := os.Setenv("HOST", c.Server.Host); err != nil {
		return errors.Wrap(err, "c.Server.Host: ")
	}
	if err := os.Setenv("PORT", c.Server.Port); err != nil {
		return errors.Wrap(err, "c.Server.Port: ")
	}

	if err := os.Setenv("dbname", c.Database.DBname); err != nil {
		return errors.Wrap(err, "c.Database.DBname: ")
	}
	if err := os.Setenv("user", c.Database.User); err != nil {
		return errors.Wrap(err, "c.Database.User: ")
	}
	if err := os.Setenv("password", c.Database.Password); err != nil {
		return errors.Wrap(err, "c.Database.Password: ")
	}
	if err := os.Setenv("timezone", c.Database.TimeZone); err != nil {
		return errors.Wrap(err, "c.Database.TimeZone: ")
	}
	if err := os.Setenv("sslmode", c.Database.SSLMODE); err != nil {
		return errors.Wrap(err, "c.Database.SSLMODE: ")
	}
	if err := os.Setenv("dbport", c.Database.DBport); err != nil {
		return errors.Wrap(err, "c.Database.DBport: ")
	}

	if err := os.Setenv("secretkey", c.SecretKey); err != nil {
		return errors.Wrap(err, "c.SecretKey: ")
	}

	return nil
}

func (c *Config) SetUpRailway() {
	c.Server.Host = os.Getenv("HOST")
	c.Server.Port = os.Getenv("PORT")

	c.Database.DBname = os.Getenv("dbname")
	c.Database.User = os.Getenv("user")
	c.Database.Password = os.Getenv("password")
	c.Database.TimeZone = os.Getenv("timezone")
	c.Database.SSLMODE = os.Getenv("sslmode")
	c.Database.DBport = os.Getenv("dbport")

	c.SecretKey = os.Getenv("secretkey")
}

// Takes env variables and creates dsn for gorm database connection
func (c *Config) CreateDSN() (string, error) {

	switch {
	case c.Server.Host == "":
		return "", errors.New("c.Server.Host is empty")
	case c.Database.User == "":
		return "", errors.New("c.Database.User is empty")
	case c.Database.Password == "":
		return "", errors.New("c.Database.Password is empty")
	case c.Database.DBname == "":
		return "", errors.New("c.Database.DBname is empty")
	case c.Database.DBport == "":
		return "", errors.New("c.Database.DBport is empty")
	case c.Database.SSLMODE == "":
		return "", errors.New("c.Database.SSLMODE is empty")
	case c.Database.TimeZone == "":
		return "", errors.New("c.Database.TimeZone is empty")
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Server.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.DBname,
		c.Database.DBport,
		c.Database.SSLMODE,
		c.Database.TimeZone,
	), nil
}

// sets up database and env variables
func SetupConfiguration(test bool) (*gorm.DB, error) {
	var configFile, path string
	if test == true {
		configFile, path = "../config.yaml", "../"
	} else {
		configFile, path = "config.yaml", "."
	}

	var config Config
	// checking if yaml file with config declartions exists, otherwise, use the env varibale provided by railway
	if _, err := os.Stat(configFile); err == nil {
		config.SetUpViper("config", path, "yaml")
	} else if errors.Is(err, os.ErrNotExist) {
		config.SetUpRailway()
	}

	dsn, err := config.CreateDSN()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
