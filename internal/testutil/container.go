package testutil

import (
	"context"
	"github.com/echovisionlab/aws-weather-api/pkg/constants"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"math/rand"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"
)

var ProjectRootPattern = regexp.MustCompile(`^(.*aws-weather-api).*`)

func ShutdownContainer(ctx context.Context, t *testing.T, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("failed to shutdown container: %+v", container)
	}
}

type GeneratedContainerReport struct {
	Container testcontainers.Container
	Host      string
	Port      string
}

type PostgresContainer struct {
	testcontainers.Container
	Host   string
	Port   string
	User   string
	Pass   string
	DbName string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func SetupPostgres(ctx context.Context, t *testing.T) testcontainers.Container {
	name, user, pass := RandStringBytes(10), RandStringBytes(10), RandStringBytes(10)
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		postgres.WithInitScripts(path.Join(TestDataPath, "db.sql")),
		postgres.WithDatabase(name),
		postgres.WithUsername(user),
		postgres.WithPassword(pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("failed to retrieve container endpoint: %+v", err)
	}

	parts := strings.Split(endpoint, ":")

	t.Setenv(constants.DatabaseName, name)
	t.Setenv(constants.DatabaseUser, user)
	t.Setenv(constants.DatabasePass, pass)
	t.Setenv(constants.DatabaseHost, parts[0])
	t.Setenv(constants.DatabasePort, parts[1])

	log.Println("database user: ", user)
	log.Println("database pass: ", pass)

	return container
}
