import requests

es_api = {
	"host": "http://localhost",
	"port": ":8001",
	"schema": {
		"api": {
			"http_method": requests.post,
			"path": "/api/",
			"mapping": '''{
				"mappings": {
					"users": {
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
							}
						}
					}
				}
			}''',
		},
	}
}


for es_schema in es_api["schema"]:
	print es_schema
	url = es_api["host"] + es_api["port"] + es_api["schema"][es_schema]["path"]
	resp = es_api["schema"][es_schema]["http_method"](url, data = es_api["schema"][es_schema]["mapping"])
	print resp.text

