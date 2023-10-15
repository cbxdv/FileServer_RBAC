package main

import (
	"encoding/json"
	"fs_backend/apierrors"
	"log"
	"net/http"
)

func ErrorResponseWriter(res http.ResponseWriter, errCode string, statusCode int) {
	resp := make(map[string]any)
	resp["error"] = map[string]string{
		"code":        errCode,
		"description": apierrors.GetErrorCodeDescription(errCode),
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln("Error when JSON marshal : ", err.Error())
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(statusCode)
	res.Write(data)
}

func JsonResponseWriter(res http.ResponseWriter, dataMap map[string]any, statusCode int) {
	if len(dataMap) != 0 {
		data, err := json.Marshal(dataMap)
		if err != nil {
			log.Fatalln("Error when JSON marshal : ", err.Error())
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.WriteHeader(statusCode)
		res.Write(data)
	}
}
