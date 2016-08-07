package main

//package tools

import (
	"fmt"
	"github.com/buger/jsonparser"
	//"net/http"
	"strings"

	//"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/pquerna/ffjson/ffjson"
	//"github.com/mailru/easyjson/jwriter"
)

type TransformationTracer struct {
	/*
		this struct is meant to keep track of
	*/

	//trackers
	objects    map[string]interface{}   // map of {path: obj} of objects processed to be re-assembled
	newObjs    []map[string]interface{} // array that keeps track of {path: obj} to objects not yet processed
	pathToObjs []string                 // array of paths to objects in original dismantled obj

	arrays    map[string]interface{} // array that keeps track of {path: ary} to objects not yet processed
	newArrays []interface{}          // array that keeps track of arrays not yet processed
	//newArrays []map[string]interface{} // array that keeps track of arrays not yet processed

	//objects
	preprocessed []byte                 // original data to be remodled
	processed    map[string]interface{} // result of remodling

	structureType string // ex, json (currently only support json)
}

func dismantleObj(remodler *TransformationTracer) {
	for len(remodler.newObjs) > 0 {
		for currentIndex := 0; currentIndex < len(remodler.newObjs); currentIndex++ {
			for key, value := range remodler.newObjs[0] {
				if key == "-----_PATH_TO_OBJECT_-----" {
					continue
				}
				switch value.(interface{}).(type) {
				case string:
					val, _, _, err := jsonparser.Get(remodler.preprocessed, strings.Split((value).(string), ".")...)
					if err == nil {
						remodler.newObjs[0][key] = string(val)
					}
				case map[string]interface{}:
					remodler.newObjs = append(remodler.newObjs, value.(map[string]interface{}))
					if val, ok := remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"]; ok {
						remodler.newObjs[len(remodler.newObjs)-1]["-----_PATH_TO_OBJECT_-----"] = val.(string) + "." + key
					} else {
						remodler.newObjs[len(remodler.newObjs)-1]["-----_PATH_TO_OBJECT_-----"] = key
					}
					delete(remodler.newObjs[0], key)
				case []interface{}:
					for _, item := range value.([]interface{}) {
						remodler.newArrays = append(remodler.newArrays, item)
						generateArray(remodler)
					}
				default: // for any values that cant be a path or nested element keep whats set as default and move on
					continue
				}
			}
			pathToObj := remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"].(string)
			if pathToObj != "~~~~~-----%!%!%Root%!%!%-----~~~~~" {
				remodler.pathToObjs = append(remodler.pathToObjs, pathToObj)
			}
			remodler.objects[pathToObj] = remodler.newObjs[0]
			delete(remodler.objects[pathToObj].(map[string]interface{}), "-----_PATH_TO_OBJECT_-----")
			remodler.newObjs = remodler.newObjs[1:]
		}
	}
}

func generateArray(remodler *TransformationTracer) {
	for len(remodler.newArrays) > 0 {
		fmt.Println(remodler.newArrays[0])
		//jsonparser.ArrayEach(remodler.preprocessed, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		//	fmt.Println(jsonparser.Get(value, "url"))
		//}, strings.Split((remodler.newArrays[0]).(string), ".")...) //"person", "gravatar", "avatars")
		remodler.newArrays = remodler.newArrays[1:]
	}
	/*for item := range value.([]interface{}) {
		remodler.newArrays[d] = append(remodler.newArrays, item.(interface{}))
	}

	switch item.(interface{}).(type) {
	case map[string]interface{}:
		remodler.newArrays = append(remodler.newArrays, item.(map[string]interface{}))
	//dismantleArray(remodler)
	fmt.Println(item)
	}
	*/

}

func reassembleObj(remodler *TransformationTracer) {
	for len(remodler.pathToObjs) > 0 {
		index := len(remodler.pathToObjs) - 1
		path := strings.Split(remodler.pathToObjs[index], ".")
		parent := strings.Join(path[:len(path)-1], ".")

		remodler.objects[parent].(map[string]interface{})[path[len(path)-1]] = remodler.objects[remodler.pathToObjs[index]]
		delete(remodler.objects, remodler.pathToObjs[index])
		remodler.pathToObjs = remodler.pathToObjs[:index]
	}
	remodler.objects = remodler.objects["~~~~~-----%!%!%Root%!%!%-----~~~~~"].(map[string]interface{})
}

//func Remodel(w http.ResponseWriter, expected []byte, original []byte) {
func Remodel(expected []byte, original []byte) { //[]byte{
	expectedJsonItr := jlexer.Lexer{Data: expected}

	//map string interface representation of json input
	expectedJsonMSI := expectedJsonItr.Interface().(map[string]interface{})
	remodler := TransformationTracer{
		make(map[string]interface{}, 0),   // obj
		make([]map[string]interface{}, 0), // obj
		make([]string, 0),                 // pathToObjs
		make(map[string]interface{}, 0),   // arrays
		//make([]map[string]interface{}, 0), // arrays
		make([]interface{}, 0),
		original,                     // preprocessed
		make(map[string]interface{}), // processed
		"json", // string
	}

	//printJ(expectedJsonMSI)
	remodler.newObjs = append(remodler.newObjs, expectedJsonMSI)
	remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"] = "~~~~~-----%!%!%Root%!%!%-----~~~~~"
	dismantleObj(&remodler)
	//printJ(remodler.objects)
	printI(remodler.newArrays)
	reassembleObj(&remodler)
	printJ(remodler.objects)

	//buf, _ := ffjson.Marshal(&remodler.objects) //&expectedJsonMSI)
	//return buf
	//return expectedJsonMSI
	//fmt.Println(remodler.pathToObjs)

	//fmt.Fprintf(w, string(buf))
	//fmt.Println(string(buf))
}

func printJ(JsonMSI map[string]interface{}) {
	buf, _ := ffjson.Marshal(&JsonMSI)
	fmt.Println(string(buf))
}

func printI(JsonMSI interface{}) {
	buf, _ := ffjson.Marshal(&JsonMSI)
	fmt.Println(string(buf))
}

//func RemodelJ(w http.ResponseWriter, r *http.Request) {
func main() {
	a := []byte(`{
		"a": 1,
		"b": 1.333,
		"e": true,
		"c": false,
		"d": "success",
		"values": [{
			"bool": "path.to.bool",
			"float": "path.to.float",
			"int": "path.to.int",
			"obj": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"string": "path.to.string"
			},
			"string": "path.to.string"
		}],
		"bool": "path.to.bool",
		"float": "path.to.float",
		"int": "path.to.int",
		"obj1": {
			"arrayElem": "array.0.string",
			"bool": "path.to.bool",
			"float": "path.to.float",
			"int": "path.to.int",
			"obj": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"obj": {
					"arrayElem": "array.0.string",
					"bool": "path.to.bool",
					"float": "path.to.float",
					"int": "path.to.int",
					"string": "path.to.string"
				},
				"string": "path.to.string"
			},
			"string": "path.to.string"
		},
		"obj2": {
			"arrayElem": "array.0.string",
			"bool": "path.to.bool",
			"float": "path.to.float",
			"int": "path.to.int",
			"obj": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"obj": {
					"arrayElem": "array.0.string",
					"bool": "path.to.bool",
					"float": "path.to.float",
					"int": "path.to.int",
					"string": "path.to.string"
				},
				"string": "path.to.string"
			},
			"string": "path.to.string",
			"obj1": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"obj": {
					"arrayElem": "array.0.string",
					"bool": "path.to.bool",
					"float": "path.to.float",
					"int": "path.to.int",
					"obj2": {
						"arrayElem": "array.0.string",
						"bool": "path.to.bool",
						"float": "path.to.float",
						"int": "path.to.int",
						"obj": {
							"arrayElem": "array.0.string",
							"bool": "path.to.bool",
							"float": "path.to.float",
							"int": "path.to.int",
							"obj": {
								"arrayElem": "array.0.string",
								"bool": "path.to.bool",
								"float": "path.to.float",
								"int": "path.to.int",
								"string": "path.to.string"
							},
							"string": "path.to.string"
						},
						"string": "path.to.string"
					},
					"obj": {
						"arrayElem": "array.0.string",
						"bool": "path.to.bool",
						"float": "path.to.float",
						"int": "path.to.int",
						"string": "path.to.string"
					},
					"string": "path.to.string"
				},
				"string": "path.to.string"
			}
		},
		"obj3": {
			"arrayElem": "array.0.string",
			"bool": "path.to.bool",
			"float": "path.to.float",
			"int": "path.to.int",
			"obj": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"string": "path.to.string"
			},
			"obj1": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"string": "path.to.string"
			},
			"obj2": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"string": "path.to.string"
			},
			"obj3": {
				"arrayElem": "array.0.string",
				"bool": "path.to.bool",
				"float": "path.to.float",
				"int": "path.to.int",
				"obj": {
					"arrayElem": "array.0.string",
					"bool": "path.to.bool",
					"float": "path.to.float",
					"int": "path.to.int",
					"obj": {
						"arrayElem": "array.0.string",
						"bool": "path.to.bool",
						"float": "path.to.float",
						"int": "path.to.int",
						"string": "path.to.string",
						"obj": {
							"arrayElem": "array.0.string",
							"bool": "path.to.bool",
							"float": "path.to.float",
							"int": "path.to.int",
							"string": "path.to.string",
							"obj": {
								"arrayElem": "array.0.string",
								"bool": "path.to.bool",
								"float": "path.to.float",
								"int": "path.to.int",
								"string": "path.to.string",
								"obj": {
									"arrayElem": "array.0.string",
									"bool": "path.to.bool",
									"float": "path.to.float",
									"int": "path.to.int",
									"string": "path.to.string",
									"obj": {
										"arrayElem": "array.0.string",
										"bool": "path.to.bool",
										"float": "path.to.float",
										"int": "path.to.int",
										"string": "path.to.string"
									},
									"obj1": {
										"arrayElem": "array.0.string",
										"bool": "path.to.bool",
										"float": "path.to.float",
										"int": "path.to.int",
										"string": "path.to.string",
										"obj": {
											"arrayElem": "array.0.string",
											"bool": "path.to.bool",
											"float": "path.to.float",
											"int": "path.to.int",
											"string": "path.to.string"
										}
									}
								}
							}
						}
					},
					"string": "path.to.string"
				},
				"string": "path.to.string"
			}
		}
	}`)
	in := []byte(`{
          "path": {
            "to": {
              "string": "i am a string",
              "float": "i am a float",
              "bool": "i am a bool",
              "int": "i am an int",

              "array": [{
                "string": "i am a string",
                "float": "i am a float",
                "bool": "i am a bool",
                "int": "i am an int"
              },{
                "string": "i am a string",
                "float": "i am a float",
                "bool": "i am a bool",
                "int": "i am an int"
              }],
            }
          }
        }`)
	/*a = []byte(`{
	  "array": [{
		"########ORIGINAL_PATH_TO_ARRAY########": "path.to.array",
		"string": "string",
		"float": "float",
		"bool": "bool",
		"int": "int"
	  }]
	}`)*/

	//Remodel(w, a, in)
	Remodel(a, in)
	//return ""
}
