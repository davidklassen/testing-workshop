//go:build integration

package greeter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/davidklassen/testing-workshop/pkg/dbutils"
	"github.com/davidklassen/testing-workshop/pkg/testutils"
	"github.com/ory/dockertest/v3"
)

var (
	httpPort string
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to connect to docker: %s", err.Error())
	}

	network, err := pool.CreateNetwork("greeter")
	if err != nil {
		log.Fatalf("failed to create docker network: %s", err.Error())
	}

	// setup postgres container
	user, pass, dbName := "greeter", "greeter", "greeter"
	pgContainer, err := testutils.CreatePostgres(pool, network, user, pass, dbName)
	if err != nil {
		log.Printf("failed to create potgres: %s", err.Error())
		testutils.Cleanup(1, pool, network, pgContainer)
	}

	// run migrations
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %s", err.Error())
	}

	projectRoot := path.Join(pwd, "..", "..")
	migrations := path.Join(projectRoot, "migrations")
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, pass, pgContainer.GetPort("5432/tcp"), dbName)
	err = dbutils.MigratePostgres(connStr, dbName, migrations)
	if err != nil {
		log.Printf("failed to migrate database: %s", err.Error())
		testutils.Cleanup(1, pool, network, pgContainer)
	}

	// setup greeter container
	greeterContainer, err := createGreeter(pool, network, projectRoot)
	if err != nil {
		log.Printf("failed to create greeter: %s", err.Error())
		testutils.Cleanup(1, pool, network, pgContainer, greeterContainer)
	}
	httpPort = greeterContainer.GetPort("8080/tcp")

	// run tests
	testutils.Cleanup(m.Run(), pool, network, pgContainer, greeterContainer)
}

func createGreeter(pool *dockertest.Pool, network *dockertest.Network, projectRoot string) (*dockertest.Resource, error) {
	buildOpts := &dockertest.BuildOptions{
		Dockerfile: "cmd/greeter/Dockerfile",
		ContextDir: projectRoot,
	}
	runOpts := &dockertest.RunOptions{
		Name:     "greeter",
		Networks: []*dockertest.Network{network},
	}

	resource, err := pool.BuildAndRunWithBuildOptions(buildOpts, runOpts)
	if err != nil {
		return resource, fmt.Errorf("failed to build and run: %w", err)
	}

	httpPort = resource.GetPort("8080/tcp")

	if err = pool.Retry(func() error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/readyz", httpPort))
		if err != nil {
			return err
		}

		if resp.StatusCode > 299 {
			return errors.New("got http error code")
		}

		return nil
	}); err != nil {
		return resource, fmt.Errorf("failed to connect to container: %w", err)
	}

	return resource, nil
}
