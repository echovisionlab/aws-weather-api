package testutil

import (
	"github.com/echovisionlab/aws-weather-api/pkg/constants"
	"testing"
)

func ResetDatabaseEnv(t *testing.T) {
	t.Setenv(constants.DatabaseHost, "")
	t.Setenv(constants.DatabasePort, "")
	t.Setenv(constants.DatabaseUser, "")
	t.Setenv(constants.DatabasePass, "")
	t.Setenv(constants.DatabaseName, "")
}
