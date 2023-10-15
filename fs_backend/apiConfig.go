package main

import (
	"log"
	"os"

	"fs_backend/databaseservice"
	"fs_backend/fileservice"
	"fs_backend/transferpropertiesservice"
)

type ApiConfig struct {
	ServerPort string
	jwtSecret  string

	graphService         databaseservice.GraphDatabaseService
	transferPropsService transferpropertiesservice.TransferPropertiesService
	fileService          fileservice.FileService
}

func (apifn *ApiConfig) initialize() {
	apifn.readEnv()
	apifn.graphService.Connect()
	apifn.transferPropsService.Start()
	apifn.fileService.Initialize()
}

func (apifn *ApiConfig) readEnv() {
	apifn.jwtSecret = os.Getenv("JWT_SECRET")
	if apifn.jwtSecret == "" {
		log.Default().Println("Needed a JWT secret")
	}
}

func (apifn ApiConfig) close() {
	apifn.transferPropsService.Stop()
	apifn.graphService.Close()
	apifn.fileService.Cleanup()
}
