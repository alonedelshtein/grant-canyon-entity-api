package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	// Import the Elasticsearch library packages

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/opensearch-project/opensearch-go"
)

const queryByTerm = `{
	"query": {
	  "bool": {
		"should": [
				{ "match": { "name":   "%s"  }},
				{ "match": { "program": "%s" }}
			],
			"minimum_should_match": 1,
			"filter": [
				{ "term":  { "fund": "%s" }}
			]
		}
	},
	"sort": [
	  {
	    "due_date": {
	      "order": "desc"
	    }
	  }
	]
  }`

func PostBulk(client *opensearch.Client, index string, docs string) (map[string]interface{}, error) {
	ctx := context.Background()

	var mapResp map[string]interface{}

	// Instantiate a request object
	req := esapi.BulkRequest{
		Index:   index,
		Body:    strings.NewReader(docs),
		Refresh: "true",
	}

	// Return an API response object from request
	res, _ := req.Do(ctx, client)
	defer res.Body.Close()
	var err error
	if err = json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	return mapResp, err
}

func SearchAllIndex(client *opensearch.Client, index string) []map[string]interface{} {
	var query = `{"query": {"match_all" : {}}}`

	return SearchByQuery(client, index, query)
}

func SearchById(client *opensearch.Client, index string, docId string) []map[string]interface{} {
	var query = fmt.Sprintf(`{"query": {"ids" : {"values" : ["%s"]}}}`, docId)

	return SearchByQuery(client, index, query)
}

func SearchByTerm(client *opensearch.Client, index string, term string, fund string) []map[string]interface{} {
	var query = fmt.Sprintf(queryByTerm, term, term, fund)

	return SearchByQuery(client, index, query)
}

func SearchByQuery(client *opensearch.Client, index string, query string) []map[string]interface{} {

	ctx := context.Background()

	var mapResp map[string]interface{}
	var buf bytes.Buffer
	var b strings.Builder
	b.WriteString(query)
	read := strings.NewReader(b.String())

	if err := json.NewEncoder(&buf).Encode(read); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(read),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithSize(1000),
	)

	if err != nil {
		log.Fatalf("Elasticsearch Search() API ERROR: %s", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	sources := extractAllSources(mapResp)

	return sources
}

func extractAllSources(mapResp map[string]interface{}) []map[string]interface{} {
	sources := make([]map[string]interface{}, 0)

	if _, ok := mapResp["hits"]; !ok {
		return sources
	}

	// Iterate the document "hits" returned by API call
	for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {

		// Parse the attributes/fields of the document
		doc := hit.(map[string]interface{})

		// The "_source" data is another map interface nested inside of doc
		source := doc["_source"].(map[string]interface{})

		sources = append(sources, source)
	} // end of response iteration

	return sources
}
