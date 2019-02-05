package config

import (
	"os"
	"strconv"
	"testing"
)

const (
	dbHostEnvName string = "DB_HOST"
	dbHostDefault string = "localhost"

	dbPortEnvName string = "DB_PORT"
	dbPortDefault string = "27017"

	dbNameEnvName string = "DB_NAME"
	dbNameDefault string = "test"

	svcContextTimeoutEnvName string = "CONTEXT_TIMEOUT"
	svcContextTimeoutDefault int    = 10

	svcPortEnvName string = "LISTEN_PORT"
	svcPortDefault string = "8080"

	svcDebugEnvName string = "DEBUG"
	svcDebugDefault bool   = true
)

func TestGet(t *testing.T) {
	cfgExpected := getExpectedConfig(t)
	if cfgExpected == nil {
		t.Errorf("expected config object is not poopulated")
	}

	cfgActual, err := Get()
	if err != nil {
		t.Errorf("unexpected errors while parsing config: %s", err.Error())
	}
	if cfgActual == nil {
		t.Errorf("actual config object is not poopulated")
	}

	if configsMatch := checkConfigsMatch(t, cfgExpected, cfgActual); !configsMatch {
		t.Error("Parsed config does not match expected config")
	}
}

func getExpectedConfig(t *testing.T) *Config {

	//constructing expected config struct
	cfgExpected := &Config{}

	//expected DB config
	cfgExpected.DB = DBConf{
		Host: os.Getenv(dbHostEnvName),
		Port: os.Getenv(dbPortEnvName),
		Name: os.Getenv(dbNameEnvName),
	}
	if cfgExpected.DB.Host == "" {
		cfgExpected.DB.Host = dbHostDefault
	}
	if cfgExpected.DB.Port == "" {
		cfgExpected.DB.Port = dbPortDefault
	}
	if cfgExpected.DB.Name == "" {
		cfgExpected.DB.Name = dbNameDefault
	}

	//Expected service config
	cfgExpected.Service = ServiceConf{
		CtxTimeout: 10,
		Debug:      true,
		Port:       os.Getenv(svcPortEnvName),
	}

	//svc.CtxTimeout
	if envVarStr, isSet := os.LookupEnv(svcContextTimeoutEnvName); isSet {
		envVar, err := strconv.ParseInt(envVarStr, 10, 0)
		if err != nil {
			t.Logf("env variable %s is set to %s, which cannot be parsed to an integer", svcContextTimeoutEnvName, envVarStr)
			cfgExpected.Service.CtxTimeout = svcContextTimeoutDefault
		} else {
			cfgExpected.Service.CtxTimeout = int(envVar)
		}
	} else {
		cfgExpected.Service.CtxTimeout = svcContextTimeoutDefault
	}

	//svc.Debug
	if envVarStr, isSet := os.LookupEnv(svcDebugEnvName); isSet {
		envVar, err := strconv.ParseBool(envVarStr)
		if err != nil {
			t.Logf("env variable %s is set to %s, which cannot be parsed to a boolean", svcDebugEnvName, envVarStr)
			cfgExpected.Service.Debug = svcDebugDefault
		} else {
			cfgExpected.Service.Debug = envVar
		}
	} else {
		cfgExpected.Service.Debug = svcDebugDefault
	}

	//svc.Port
	if cfgExpected.Service.Port == "" {
		cfgExpected.Service.Port = svcPortDefault
	}

	return cfgExpected
}

func checkConfigsMatch(t *testing.T, expected, actual *Config) bool {

	isOk := true

	if expected == nil {
		t.Log("failed to check the two configs as the expected config object is nil")
		isOk = false
	}

	if actual == nil {
		t.Log("failed to check the two configs as the actual config object is nil")
		isOk = false
	}

	if !isOk {
		//cannot proceed further
		return false
	}

	isOk = compareTwoStrings(t, "DB host", expected.DB.Host, actual.DB.Host) && isOk
	isOk = compareTwoStrings(t, "DB port", expected.DB.Port, actual.DB.Port) && isOk
	isOk = compareTwoStrings(t, "DB name", expected.DB.Name, actual.DB.Name) && isOk

	isOk = compareTwoIntegers(t, "Service ctx timeout", expected.Service.CtxTimeout, actual.Service.CtxTimeout) && isOk
	isOk = compareTwoStrings(t, "Service port", expected.Service.Port, actual.Service.Port) && isOk
	isOk = compareTwoBooleans(t, "Service debug", expected.Service.Debug, actual.Service.Debug) && isOk

	return isOk
}

func compareTwoStrings(t *testing.T, name, expected, actual string) bool {

	if expected != actual {
		t.Logf("%s values don't match; expected: %s, got: %s", name, expected, actual)
		return false
	}

	return true
}

func compareTwoIntegers(t *testing.T, name string, expected, actual int) bool {

	if expected != actual {
		t.Logf("%s values don't match; expected: %d, got: %d", name, expected, actual)
		return false
	}

	return true
}

func compareTwoBooleans(t *testing.T, name string, expected, actual bool) bool {

	if expected != actual {
		t.Logf("%s values don't match; expected: %t, got: %t", name, expected, actual)
		return false
	}

	return true
}
