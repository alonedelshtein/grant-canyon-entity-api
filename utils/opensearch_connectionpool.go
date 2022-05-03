package utils

import (
	"fmt"
	"math/rand"
	"net/url"

	// Import the Elasticsearch library packages

	"github.com/opensearch-project/opensearch-go"
	opensearchtransport "github.com/opensearch-project/opensearch-go/opensearchtransport"
)

//Connection pool
type CustomConnectionPool struct {
	urls []*url.URL
}

// Next returns a random connection.
func (cp *CustomConnectionPool) Next() (*opensearchtransport.Connection, error) {
	u := cp.urls[rand.Intn(len(cp.urls))]
	return &opensearchtransport.Connection{URL: u}, nil
}

func (cp *CustomConnectionPool) OnFailure(c *opensearchtransport.Connection) error {
	var index = -1
	for i, u := range cp.urls {
		if u == c.URL {
			index = i
		}
	}
	if index > -1 {
		cp.urls = append(cp.urls[:index], cp.urls[index+1:]...)
		return nil
	}
	return fmt.Errorf("connection not found")
}
func (cp *CustomConnectionPool) OnSuccess(c *opensearchtransport.Connection) error { return nil }
func (cp *CustomConnectionPool) URLs() []*url.URL                                  { return cp.urls }

func CreateOpensearchPool() (*opensearch.Client, error) {
	numberOfConnections := 10
	var urls []*url.URL

	for i := 0; i < numberOfConnections; i++ {
		urls = append(urls, &url.URL{Scheme: "https", Host: "search-grant-canyon-dev-2j6iinn6e2xztpg5wplcrwbhca.eu-central-1.es.amazonaws.com"})
	}
	return opensearch.NewClient(opensearch.Config{

		ConnectionPoolFunc: func(conns []*opensearchtransport.Connection, selector opensearchtransport.Selector) opensearchtransport.ConnectionPool {
			return &CustomConnectionPool{
				urls: urls,
			}
		},
	})
}
