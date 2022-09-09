package testutils

import (
	"fmt"
	"log"
	"os"

	"github.com/davidklassen/testing-workshop/pkg/dbutils"
	"github.com/ory/dockertest/v3"
)

func Cleanup(code int, pool *dockertest.Pool, network *dockertest.Network, resources ...*dockertest.Resource) {
	for _, resource := range resources {
		if resource != nil {
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("failed to purge resource: %s", err.Error())
			}
		}
	}
	if network != nil {
		if err := network.Close(); err != nil {
			log.Fatalf("failed to close network: %s", err.Error())
		}
	}
	os.Exit(code)
}

func CreatePostgres(pool *dockertest.Pool, network *dockertest.Network, user, pass, dbName string) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "pg",
		Repository: "postgres",
		Tag:        "14.2",
		Env: []string{
			"PGUSER=" + user,
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + pass,
			"POSTGRES_DB=" + dbName,
		},
		Networks: []*dockertest.Network{network},
	})
	if err != nil {
		return resource, fmt.Errorf("failed to run postgres: %w", err)
	}

	if err = pool.Retry(func() error {
		connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, pass, resource.GetPort("5432/tcp"), dbName)
		_, err := dbutils.OpenPostgres(connStr)
		return err
	}); err != nil {
		return resource, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return resource, nil
}
