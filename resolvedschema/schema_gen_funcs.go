package resolvedschema

import("encoding/json")

func RAMParaStruct(param1 *map[string]interface{}) RAM {
	var returnStruct RAM
	temp, err := json.Marshal(param1)
	if err != nil {
		return RAM{} 
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return RAM{}
	}
	return returnStruct
}