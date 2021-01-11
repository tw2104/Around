package main

import (
    "context"
    "fmt"

    "github.com/olivere/elastic"
)

const (
        ES_URL = "http://10.128.0.2:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	// create a connection
    client, err := elastic.NewClient(
        elastic.SetURL(ES_URL),
        elastic.SetBasicAuth("elastic", "12345678"))
    if err != nil {
        return nil, err
    }

    searchResult, err := client.Search().
        Index(index). // search in index "twitter"
        Query(query). // specify the query
        Pretty(true). // pretty print request and response json
        Do(context.Background()) // execute
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

func saveToES(i interface{}, index string, id string) error {
    // create a connection
    client, err := elastic.NewClient(
        elastic.SetURL(ES_URL),
        elastic.SetBasicAuth("elastic", "12345678"))
    if err != nil {
        return err
    }

    _, err = client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())

    if err != nil {
        return err
    }

    fmt.Printf("Post is saved to ElasticSearch: %s\n", id)
    return nil
}

