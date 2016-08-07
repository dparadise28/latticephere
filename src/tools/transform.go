//package main

package tools

import (
	"fmt"
	"github.com/buger/jsonparser"
	"net/http"
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
	currentLocation string // path to set of values currently being processed

	objects     map[string]interface{}   // array that keeps track of {path: obj} to objects processed
	newObjs     []map[string]interface{} // array that keeps track of {path: obj} to objects not yet processed
	moreObjects bool                     // check updated on tracer pass throughs to see if
	// any more objects have been found and another pass is necissary

	pathToObjs []string // array of paths to objects in original dismantled obj

	arrays     []map[string]interface{} // array that keeps track of {path: ary} to objects not yet processed
	moreArrays int                      // check updated on tracer pass throughs to see if
	// any more arrays have been found and another pass is necissary

	//objects
	preprocessed []byte                 // original data to be remodled
	processed    map[string]interface{} // result of remodling

	structureType string // ex, json (currently only support json)
}

func dismantleObj(remodler *TransformationTracer) {
	for len(remodler.newObjs) > 0 {
		remodler.moreObjects = false
		for currentIndex := 0; currentIndex < len(remodler.newObjs); currentIndex++ {
			for key, value := range remodler.newObjs[0] {
				if key == "-----_PATH_TO_OBJECT_-----" {
					continue
				}
				switch value.(interface{}).(type) {
				case int:
					remodler.newObjs[0][key] = value.(int)
				case bool:
					remodler.newObjs[0][key] = value.(bool)
				case float32:
					remodler.newObjs[0][key] = value.(float32)
				case float64:
					remodler.newObjs[0][key] = value.(float64)
				case string:
					val, _, _, _ := jsonparser.Get(remodler.preprocessed, strings.Split((value).(string), ".")...)
					remodler.newObjs[0][key] = string(val)
				case map[string]interface{}:
					remodler.newObjs = append(remodler.newObjs, value.(map[string]interface{}))
					if val, ok := remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"]; ok {
						remodler.newObjs[len(remodler.newObjs)-1]["-----_PATH_TO_OBJECT_-----"] = val.(string) + "." + key
					} else {
						remodler.newObjs[len(remodler.newObjs)-1]["-----_PATH_TO_OBJECT_-----"] = key
					}
					delete(remodler.newObjs[0], key)
				default:
					if true {
						break
					} else {
						break
					}
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

func Remodel(w http.ResponseWriter, expected []byte, original []byte) {
	//func Remodel(expected []byte, original []byte) { //[]byte{
	expectedJsonItr := jlexer.Lexer{Data: expected}

	//map string interface representation of json input
	expectedJsonMSI := expectedJsonItr.Interface().(map[string]interface{})
	remodler := TransformationTracer{
		"", // path
		make(map[string]interface{}, 0),   // obj
		make([]map[string]interface{}, 0), // obj
		true,                              // obj check
		make([]string, 0),                 // pathToObjs
		make([]map[string]interface{}, 0), // arrays
		0,                            // array check
		original,                     // preprocessed
		make(map[string]interface{}), // processed
		"json", // string
	}

	//printJ(expectedJsonMSI)
	remodler.newObjs = append(remodler.newObjs, expectedJsonMSI)
	remodler.newObjs[0]["-----_PATH_TO_OBJECT_-----"] = "~~~~~-----%!%!%Root%!%!%-----~~~~~"
	dismantleObj(&remodler)
	//printJ(remodler.objects)
	reassembleObj(&remodler)
	//printJ(remodler.objects)

	buf, _ := ffjson.Marshal(&remodler.objects) //&expectedJsonMSI)
	//return buf
	//return expectedJsonMSI
	//fmt.Println(remodler.pathToObjs)

	fmt.Fprintf(w, string(buf))
	//fmt.Println(string(buf))
}

func printJ(JsonMSI map[string]interface{}) {
	buf, _ := ffjson.Marshal(&JsonMSI)
	fmt.Println(string(buf))
}

func RemodelJ(w http.ResponseWriter, r *http.Request) {
	//func main() {
	a := []byte(`{
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
	Remodel(w, a, in)
	//Remodel(a, in)
	//return ""
}
