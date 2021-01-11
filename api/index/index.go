package main

import (
        "context"
        "fmt"
        
        "github.com/olivere/elastic"
)

const (
        POST_INDEX = "post"
        USER_INDEX = "user"
		ES_URL = "http://10.128.0.2:9200"
)

func main() {
    // Obtain a es client
    client, err := elastic.NewClient(
        elastic.SetURL(ES_URL),
        elastic.SetBasicAuth("elastic", "12345678"))
    
    // Handler an error of connection
    if err != nil {
        panic(err)
    }

    // Use the IndexExist to check if a specific index (POST_INDEX) exists
    exists, err := client.IndexExists(POST_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        // Mapping is like schema in oracle, it is metadata of the database
        mapping := `{
                        "mappings": {
                                "properties": {
                                        "user":     { "type": "keyword" },
                                        "message":  { "type": "text" },
                                        "url":      { "type": "keyword", "index": false },
                                        "type":     { "type": "keyword", "index": false }
                                }
                        }
                }`
        _, err := client.CreateIndex(POST_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    exists, err = client.IndexExists(USER_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    if !exists {
        mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword", "index": false},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`

        // Create a index with mapping and context
        // What is context in cs? It's the information you want provide for this call, like a header for a http request
        // It is typically used to pass things not necessarily tied directly to a method call, but could still be pertinent. 
        // A layperson way of describing it might be "stuff you may care about"
        _, err = client.CreateIndex(USER_INDEX).Body(mapping).Do(context.Background())
        if err != nil {
            panic(err)
        }
    }

    fmt.Println("Indexes are created.")
}