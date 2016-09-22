package main

import (
	"fmt"
	"http"
)

var host string = "http://localhost"
var port string = ":8001"

func main() {
	schema := {
		"api": {
			"http_method": "POST",
			"path": "/api/",
			"mapping": `{
				"mappings": {
					"users": {
						"properties": {
							"dynamic": "true",
							"properties": {
								"first_name": {
									"type": "string",
									"index": "no"
								},
								"last_name": {
									"type": "string",
									"index": "no"
								},
								"password": {
									"type": "string",
									"index": "no"
								},
								"email": {
									"type": "string",
									"index": "not_analyzed"
								},
								"id": {
									"type": "string",
									"index": "not_analyzed"
								},
								"active": {
									"type": "string",
									"index": "not_analyzed"
								},
								"supplemental": {
									"type": "object",
									"index": "no"
								}
							}
						}
					}
				}
			}`,
		},
	}

//func main() {
	fmt.Println(port)
	fmt.Println(host)
	for es_schema_index, properties := range schema {
		req, err := http.NewRequest(properties["http_method"], host + path + port, properties["mapping"])
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}
