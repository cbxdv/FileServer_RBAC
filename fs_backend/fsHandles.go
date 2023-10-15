package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fs_backend/apierrors"
	"fs_backend/models"

	"github.com/google/uuid"
)

func (apifn *ApiConfig) HandleDirectoryQuery(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	// Enforcing only GET, PUT and DELETE methods
	if (req.Method != http.MethodGet) && (req.Method != http.MethodPut) && (req.Method != http.MethodDelete) {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query()
	location := query.Get("location")

	// Checking whether the data is valid
	if location == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidLocation, http.StatusBadRequest)
		return
	}

	locationSplit := strings.Split(location, "/")
	workspaceName := locationSplit[0]

	// Checking whether it is the owner of the workspace
	workspaceOwner, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodGet {
		if workspaceOwner.Id != claims.AccountId {
			nearestRoles, err := apifn.graphService.GetNearestRole(claims.AccountId, location)
			if err != nil {
				log.Default().Println(err.Error())
				ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
				return
			}
			if len(nearestRoles) == 0 {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
			nearestRole := resolveRoles(nearestRoles)
			if !(nearestRole.CanRead) {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
		}

		// Getting the directory and its contents
		directoryAndContents, err := apifn.graphService.GetDirectoryAndContentDetails(location)
		if err != nil {
			if errors.Is(err, apierrors.DirectoryNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidLocation, http.StatusBadRequest)
				return
			}
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		resData := make(map[string]any)
		resData["directoryAndContents"] = directoryAndContents
		JsonResponseWriter(res, resData, http.StatusOK)
	}
	if req.Method == http.MethodPut {
		// Owner doesn't need to check for permissions
		if workspaceOwner.Id != claims.AccountId {
			// Checking whether account has permission to create directory in the parent directory
			nearestRoles, err := apifn.graphService.GetNearestRole(claims.AccountId, location)
			if err != nil {
				log.Default().Println(err.Error())
				ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
				return
			}
			if len(nearestRoles) == 0 {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
			nearestRole := resolveRoles(nearestRoles)
			if !nearestRole.CanCreate {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
		}

		// Parsing the request body
		var params struct {
			NewDirectoryName string `json:"newDirectoryName"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		// Checking whether the data is valid
		if params.NewDirectoryName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		newDirectory := models.Directory{
			Id:        uuid.New().String(),
			Type:      "directory",
			Name:      params.NewDirectoryName,
			Location:  location + "/" + params.NewDirectoryName,
			CreatedOn: time.Now().UTC(),
		}

		err = apifn.graphService.CreateDirectory(newDirectory)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.DirectoryWithSameNameAlreadyExists{DirName: newDirectory.Name, ParentDirName: location}) {
				ErrorResponseWriter(res, apierrors.ResErrDirAlreadyExists, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		resData := make(map[string]any)
		resData["newDirectory"] = newDirectory
		JsonResponseWriter(res, resData, http.StatusCreated)
		return
	}
	if req.Method == http.MethodDelete {
		// Owner doesn't need to check for permissions
		if workspaceOwner.Id != claims.AccountId {
			if workspaceOwner.Id != claims.AccountId {
				nearestRoles, err := apifn.graphService.GetNearestRole(claims.AccountId, location)
				if err != nil {
					log.Default().Println(err.Error())
					ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
					return
				}
				if len(nearestRoles) == 0 {
					ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
					return
				}
				nearestRole := resolveRoles(nearestRoles)
				if !(nearestRole.CanDelete) {
					ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
					return
				}
			}
		}

		if len(locationSplit) == 1 {
			// Cannot delete root
			ErrorResponseWriter(res, apierrors.ResErrInvalidLocation, http.StatusBadRequest)
			return
		}

		// Checking whether the directory has content
		count, err := apifn.graphService.CountDirectoryContents(location)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		// Directory with contents cannot be deleted
		if count != 0 {
			ErrorResponseWriter(res, apierrors.ResErrDirNotEmpty, http.StatusBadRequest)
			return
		}
		// Deleting the directory
		err = apifn.graphService.DeleteDirectory(location)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		resData := make(map[string]any)
		JsonResponseWriter(res, resData, http.StatusOK)
		return
	}
}

func (apifn *ApiConfig) HandleFileQuery(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if (req.Method != http.MethodGet) && (req.Method != http.MethodPost) && (req.Method != http.MethodDelete) {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query()
	location := query.Get("location")

	// Checking whether the data is valid
	if location == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidLocation, http.StatusBadRequest)
		return
	}

	locationSplit := strings.Split(location, "/")
	workspaceName := locationSplit[0]

	// Checking whether it is the owner of the workspace
	workspaceOwner, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodGet {
		file, err := apifn.graphService.GetFileDetails(location)
		if err != nil {
			if errors.Is(err, apierrors.FileNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrResourceNotFound, http.StatusNotFound)
				return
			}
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		chunkTotal := math.Ceil(float64(file.Size) / float64(apifn.fileService.DownloadChunkSize))
		// Creating upload properties for internal use
		downloadProperties := models.FileTransferProperties{
			LinkId:         uuid.New().String(),
			FileProperties: file,
			LinkGenerated:  time.Now(),
		}
		apifn.transferPropsService.Set(downloadProperties)
		resData := map[string]any{}
		resData["chunkSize"] = apifn.fileService.DownloadChunkSize
		resData["chunkTotal"] = chunkTotal
		resData["downloadLink"] = downloadProperties.LinkId
		JsonResponseWriter(res, resData, http.StatusOK)
		return
	}
	if req.Method == http.MethodPost {
		// Parsing the request body
		var params struct {
			Name string `json:"name"`
			Size int    `json:"size"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		// Checking whether the data is valid
		if params.Name == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		// Checking permissions
		if workspaceOwner.Id != claims.AccountId {
			nearestRoles, err := apifn.graphService.GetNearestRole(claims.AccountId, location)
			if err != nil {
				log.Default().Println(err.Error())
				ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
				return
			}
			if len(nearestRoles) == 0 {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
			nearestRole := resolveRoles(nearestRoles)
			if !(nearestRole.CanCreate) {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
		}

		if params.Size == 0 {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		newFileId := uuid.New().String()

		// Creating upload properties for internal use
		uploadProperties := models.FileTransferProperties{
			LinkId: uuid.New().String(),
			FileProperties: models.File{
				Id:        newFileId,
				Type:      "file",
				Name:      params.Name,
				CreatedOn: time.Now().UTC(),
				Size:      params.Size,
				Location:  location + "/" + params.Name,
			},
			LinkGenerated: time.Now(),
		}

		// Storing upload details
		apifn.transferPropsService.Set(uploadProperties)

		// Sending the upload link
		resData := make(map[string]any)
		resData["uploadLink"] = uploadProperties.LinkId
		JsonResponseWriter(res, resData, http.StatusOK)
		return
	}
	if req.Method == http.MethodDelete {
		// Checking permissions
		if workspaceOwner.Id != claims.AccountId {
			nearestRoles, err := apifn.graphService.GetNearestRole(claims.AccountId, location)
			if err != nil {
				log.Default().Println(err.Error())
				ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
				return
			}
			if len(nearestRoles) == 0 {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
				return
			}
			nearestRole := resolveRoles(nearestRoles)
			if !(nearestRole.CanDelete) {
				ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
				return
			}
		}

		fileFromDb, err := apifn.graphService.GetFileDetails(location)
		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		err = apifn.graphService.DeleteFile(location)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		}

		// Deleting the file from storage
		go func() {
			apifn.fileService.DeleteFileFromInternalLocation(fileFromDb)
		}()

		res.WriteHeader(http.StatusOK)
	}
}

func (apifn ApiConfig) handleFileUpload(res http.ResponseWriter, req *http.Request, _ models.JWTData) {
	if req.Method != http.MethodPost {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Getting the upload ID from the url path
	uploadId := strings.TrimPrefix(req.URL.Path, "/fs/upload/")

	// Getting upload properties
	properties, err := apifn.transferPropsService.Get(uploadId)
	if err != nil {
		log.Default().Println(err.Error())
		if errors.Is(err, apierrors.UploadIdNotFound{}) {
			ErrorResponseWriter(res, apierrors.ResErrInvalidUploadId, http.StatusBadRequest)
			return
		}
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	// Reading header
	chunkTotalStr := req.Header.Get("Chunk-Total")
	chunkCurrentStr := req.Header.Get("Chunk-Current")

	if chunkTotalStr == "" || chunkCurrentStr == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	chunkTotal, err := strconv.Atoi(chunkTotalStr)
	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}
	chunkCurrent, err := strconv.Atoi(chunkCurrentStr)
	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	var chunkData struct {
		Data string `json:"data"`
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&chunkData)
	if err != nil {
		log.Default().Println(err)
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	data, err := base64.StdEncoding.DecodeString(chunkData.Data)
	if err != nil {
		log.Default().Println(err)
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	err = apifn.fileService.WriteChunkToFile(properties.FileProperties, data)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	resData := make(map[string]any)
	if chunkCurrent == chunkTotal {
		// If upload has completed, then create a record in database
		err := apifn.graphService.CreateFile(properties.FileProperties)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		resData["newFile"] = properties.FileProperties
	}

	JsonResponseWriter(res, resData, http.StatusOK)
}

func (apifn ApiConfig) handleFileDownload(res http.ResponseWriter, req *http.Request, _ models.JWTData) {
	if req.Method != http.MethodGet {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Getting the upload ID from the url path
	downloadId := strings.TrimPrefix(req.URL.Path, "/fs/download/")

	// Getting upload details
	apifn.transferPropsService.Get(downloadId)

	chunkTotalStr := req.Header.Get("Chunk-Total")
	chunkCurrentStr := req.Header.Get("Chunk-Current")

	if chunkTotalStr == "" || chunkCurrentStr == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	chunkTotal, err := strconv.Atoi(chunkTotalStr)
	if err != nil {
		log.Default().Println(err)
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}
	chunkCurrent, err := strconv.Atoi(chunkCurrentStr)
	if err != nil {
		log.Default().Println(err)
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	// Getting download properties
	properties, err := apifn.transferPropsService.Get(downloadId)
	if err != nil {
		log.Default().Println(err.Error())
		if errors.Is(err, apierrors.UploadIdNotFound{}) {
			ErrorResponseWriter(res, apierrors.ResErrInvalidUploadId, http.StatusBadRequest)
			return
		}
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	chunk, err := apifn.fileService.ReadChunkFromFile(properties.FileProperties, chunkCurrent)
	if err != nil {
		log.Default().Println(err)
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if chunkCurrent == chunkTotal {
		log.Default().Println("Chunk Total")
	}

	res.WriteHeader(200)
	res.Write(chunk)
}

func (apifn ApiConfig) HandleFSShared(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if req.Method != http.MethodGet {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}

	workspaceName := req.URL.Query().Get("workspace")
	if len(workspaceName) == 0 {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	sharedContent, err := apifn.graphService.GetSharedDirsAndFiles(claims.AccountId, workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}
	resData := make(map[string]any)
	resData["sharedContent"] = sharedContent
	JsonResponseWriter(res, resData, http.StatusOK)
}

func (apifn ApiConfig) handleDirDetailsQuery(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	location := req.URL.Query().Get("location")

	if len(location) == 0 {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	locationSplit := strings.Split(location, "/")
	workspaceName := locationSplit[0]

	// Checking whether it is the owner of the workspace
	workspaceOwner, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if workspaceOwner.Id != claims.AccountId {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
		return
	}

	// Getting the file
	dir, err := apifn.graphService.GetDirectoryDetails(location)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	// Getting the roles
	roles, err := apifn.graphService.GetItemAllRoles(location)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	resData := make(map[string]any)
	resData["directory"] = dir
	resData["roles"] = roles
	JsonResponseWriter(res, resData, http.StatusOK)
}

func (apifn ApiConfig) handleFileDetailsQuery(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	location := req.URL.Query().Get("location")

	if len(location) == 0 {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	locationSplit := strings.Split(location, "/")
	workspaceName := locationSplit[0]

	// Checking whether it is the owner of the workspace
	workspaceOwner, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if workspaceOwner.Id != claims.AccountId {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
		return
	}

	// Getting the file
	file, err := apifn.graphService.GetFileDetails(location)
	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	// Getting the roles
	roles, err := apifn.graphService.GetItemAllRoles(location)
	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	resData := make(map[string]any)
	resData["file"] = file
	resData["getRolesCypher"] = roles
	JsonResponseWriter(res, resData, http.StatusOK)
}
