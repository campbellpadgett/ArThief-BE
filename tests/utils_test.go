package tests

import (
	"AT-BE/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperSetup(t *testing.T) {
	s := utils.ServerConfig{}
	d := utils.DatabaseConfig{}
	c := utils.Config{
		Server:    s,
		Database:  d,
		SecretKey: "",
	}

	if err := c.SetUpViper("test_config", ".", "yaml"); err != nil {
		t.Errorf("[Error] SetUpViper failed: %s", err)
	}

	assert.True(t, c.Server.Host == "localhost")
	assert.True(t, c.SecretKey == "secretkey")
	assert.True(t, c.Server.Port == "8000")
	assert.True(t, c.Database.Password == "password")
}

func TestRailwaySetup(t *testing.T) {
	os.Setenv("host", "__host__")
	os.Setenv("port", "__port__")

	os.Setenv("dbname", "__dbname__")
	os.Setenv("user", "__user__")
	os.Setenv("password", "__password__")
	os.Setenv("timezone", "__timezone__")
	os.Setenv("sslmode", "__sslmode__")
	os.Setenv("dbport", "__dbport__")

	os.Setenv("secretkey", "__secretkey__")

	s := utils.ServerConfig{}
	d := utils.DatabaseConfig{}
	c := utils.Config{
		Server:    s,
		Database:  d,
		SecretKey: "",
	}

	c.SetUpRailway()

	assert.True(t, c.Database.DBname == "__dbname__")
	assert.True(t, c.Database.Password == "__password__")
	assert.True(t, c.Server.Port == "__port__")
	assert.True(t, c.SecretKey == "__secretkey__")
}
