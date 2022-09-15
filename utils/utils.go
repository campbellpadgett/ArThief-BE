package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

	return nil
}

func (c *Config) SetUpRailway() {
	c.Server.Host = os.Getenv("host")
	c.Server.Port = os.Getenv("port")

	c.Database.DBname = os.Getenv("dbname")
	c.Database.User = os.Getenv("user")
	c.Database.Password = os.Getenv("password")
	c.Database.TimeZone = os.Getenv("timezone")
	c.Database.SSLMODE = os.Getenv("sslmode")
	c.Database.DBport = os.Getenv("dbport")

	c.SecretKey = os.Getenv("secretkey")
}
