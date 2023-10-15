package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"fs_backend/apierrors"
	"fs_backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apifn ApiConfig) handleCheckWorkspaceAvailability(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	// Enforcing only GET method
	if req.Method != http.MethodGet {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Only owner accounts are allowed
	if !claims.IsOwner {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
		return
	}

	var params struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}
	workspaceExists, err := apifn.graphService.CheckWorkspace(params.Name)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}
	resData := make(map[string]any)
	resData["available"] = !workspaceExists
	JsonResponseWriter(res, resData, http.StatusOK)
}

func (apifn ApiConfig) handleWorkspaceOperations(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	// Enforcing only POST, DELETE and PATCH methods
	if req.Method != http.MethodGet && req.Method != http.MethodPut && req.Method != http.MethodDelete && req.Method != http.MethodPatch {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Only owner accounts are allowed
	if !claims.IsOwner {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusForbidden)
		return
	}

	if req.Method == http.MethodGet {
		workspaces, err := apifn.graphService.GetWorkspacesForOwner(claims.AccountId)
		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		resp := make(map[string]any)
		resp["workspaces"] = workspaces
		JsonResponseWriter(res, resp, http.StatusOK)
		return
	}
	if req.Method == http.MethodPut {
		var params struct {
			WorkspaceName string `json:"workspaceName"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}
		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}

		newWorkspace := models.Workspace{
			Id:   uuid.New().String(),
			Name: params.WorkspaceName,
		}

		err = apifn.graphService.CreateWorkspace(newWorkspace, claims.AccountId)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceAlreadyExists{}) {
				ErrorResponseWriter(res, apierrors.ResErrWorkspaceAlreadyExists, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		err = apifn.fileService.CreateWorkspaceSpace(params.WorkspaceName)
		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		resp := make(map[string]any)
		resp["workspace"] = newWorkspace
		JsonResponseWriter(res, resp, http.StatusCreated)
		return
	}
	if req.Method == http.MethodDelete {
		var params struct {
			WorkspaceName string `json:"workspaceName"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}
		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}

		// Checking whether the account in the token is the owner of the workspace
		ownerInDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if ownerInDb.Id != claims.AccountId {
			ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
			return
		}

		err = apifn.graphService.DeleteWorkspace(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		// Deleting workspace files
		err = apifn.fileService.DeleteWorkspaceSpace(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		JsonResponseWriter(res, map[string]any{}, http.StatusOK)
		return
	}
}

func (apifn ApiConfig) handleWorkspaceAccountOperations(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if !claims.IsOwner {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
	}

	if req.Method == http.MethodGet {
		query := req.URL.Query()
		workspaceName := query.Get("workspace")
		if workspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}
		accounts, err := apifn.graphService.GetAllServiceAccountsInWorkspace(workspaceName)
		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}
		resData := map[string]any{}
		resData["accounts"] = accounts
		JsonResponseWriter(res, resData, http.StatusOK)
	}
	if req.Method == http.MethodPut {
		// Checking whether the user is the owner of the workspace
		var params struct {
			WorkspaceName string `json:"workspaceName"`
			Name          string `json:"name"`
			Username      string `json:"username"`
			Password      string `json:"password"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)

		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}
		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}
		// Fetching owner details from database
		ownerInDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if ownerInDb.Id != claims.AccountId {
			ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
			return
		}
		if params.Name == "" || params.Username == "" || params.Password == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		// Check whether the workspace exists
		workspaceExists, err := apifn.graphService.CheckWorkspace(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if !workspaceExists {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}

		// Check whether the service account already exists
		account, _, err := apifn.graphService.GetServiceAccountWithWorkspace(params.Username, params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
				return
			} else if errors.Is(err, apierrors.AccountNotFound{}) {
				log.Default().Println("Account not found : New Service Account")
				// Do nothing
			} else {
				ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
				return
			}
		}
		if account.Username == params.Username {
			ErrorResponseWriter(res, apierrors.ResErrSAAlreadyExists, http.StatusBadRequest)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		newServiceAccount := models.ServiceAccount{
			Id:                  uuid.New().String(),
			Name:                params.Name,
			Username:            params.Username,
			Password:            string(passwordHash),
			ShouldResetPassword: true,
			LinkedEmail:         "",
		}

		err = apifn.graphService.CreateServiceAccount(newServiceAccount, params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		resData := make(map[string]any)
		resData["newServiceAccount"] = map[string]any{
			"id":       newServiceAccount.Id,
			"name":     newServiceAccount.Name,
			"username": newServiceAccount.Username,
		}
		JsonResponseWriter(res, resData, http.StatusCreated)
		return
	}
	if req.Method == http.MethodDelete {
		// Checking whether the user is the owner of the workspace
		var params struct {
			WorkspaceName string `json:"workspaceName"`
			Username      string `json:"username"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)

		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}
		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}
		// Fetching owner details from database
		ownerInDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if ownerInDb.Id != claims.AccountId {
			ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
			return
		}
		err = apifn.graphService.DeleteServiceAccount(params.WorkspaceName, params.Username)
		if err != nil {
			log.Default().Println(err.Error())
			if errors.Is(err, apierrors.WorkspaceNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
				return
			}
			if errors.Is(err, apierrors.AccountNotFound{}) {
				ErrorResponseWriter(res, apierrors.ResErrSANotFound, http.StatusBadRequest)
				return
			}
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		JsonResponseWriter(res, map[string]any{}, http.StatusOK)
	}
}
