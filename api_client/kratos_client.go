package api_client

import (
	"net/url"
	"sync"

	"github.com/ory/client-go"
	ory "github.com/ory/client-go"
)

var (
	oryClientOnce sync.Once
	oryClient     *ory.APIClient
)

func InitOryClient(url url.URL) *ory.APIClient {

	oryClientOnce.Do(func() { // <-- atomic, does not allow repeating

		cfg := client.NewConfiguration()
		cfg.Debug = true
		cfg.Servers = client.ServerConfigurations{
			{URL: url.String()},
		}

		oryClient = client.NewAPIClient(cfg)

	})

	return oryClient
}

func OryClient() *ory.APIClient {
	return oryClient
}

var (
	oryAdminClientOnce sync.Once
	oryAdminClient     *ory.APIClient
)

func InitOryAdminClient(url url.URL) *ory.APIClient {

	oryAdminClientOnce.Do(func() { // <-- atomic, does not allow repeating

		cfg := client.NewConfiguration()
		cfg.Servers = client.ServerConfigurations{
			{URL: url.String()},
		}

		oryAdminClient = client.NewAPIClient(cfg)

	})

	return oryAdminClient
}

func OryAdminClient() *ory.APIClient {
	return oryAdminClient
}
