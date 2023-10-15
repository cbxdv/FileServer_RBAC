package main

import (
	"encoding/json"
	"errors"
	"fs_backend/apierrors"
	"fs_backend/models"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (apifn ApiConfig) HandleGetAllRolesInWorkspace(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if req.Method != http.MethodGet {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	params := req.URL.Query()
	workspaceName := params.Get("workspaceName")

	if workspaceName == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
		return
	}

	ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if ownerIdDb.Id != claims.AccountId {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
		return
	}

	roles, err := apifn.graphService.GetAllRolesInWorkspace(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	resData := make(map[string]any)
	resData["roles"] = roles
	JsonResponseWriter(res, resData, http.StatusOK)
}

func (apifn ApiConfig) HandleRolesOperations(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if req.Method == http.MethodGet {
		var params struct {
			WorkspaceName string `json:"workspaceName"`
			RoleId        string `json:"roleId"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)

		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}
		if params.RoleId == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		if ownerIdDb.Id != claims.AccountId {
			ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
			return
		}

		role, err := apifn.graphService.GetRoleDetailsWithSAInWorkspace(params.WorkspaceName, params.RoleId)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		resData := make(map[string]any)
		resData["role"] = role
		JsonResponseWriter(res, resData, http.StatusOK)
	}

	if req.Method == http.MethodPut {
		var params struct {
			WorkspaceName string `json:"workspaceName"`
			RoleName      string `json:"name"`
			RoleDesc      string `json:"description"`
			CanRead       bool   `json:"canRead"`
			CanCreate     bool   `json:"canCreate"`
			CanRename     bool   `json:"canRename"`
			CanDelete     bool   `json:"canDelete"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)

		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}
		if params.RoleName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		if claims.AccountId != ownerIdDb.Id {
			ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
			return
		}

		role := models.Role{
			Id:          uuid.NewString(),
			Name:        params.RoleName,
			Description: params.RoleDesc,
			CanRead:     params.CanRead || false,
			CanCreate:   params.CanCreate || false,
			CanRename:   params.CanRename || false,
			CanDelete:   params.CanDelete || false,
		}

		err = apifn.graphService.CreateNewRole(role, params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		resData := make(map[string]any)
		resData["role"] = role
		JsonResponseWriter(res, resData, http.StatusCreated)
	}

	if req.Method == http.MethodDelete {
		var params struct {
			WorkspaceName string `json:"workspaceName"`
			RoleId        string `json:"roleId"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)

		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		if params.WorkspaceName == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
			return
		}
		if params.RoleId == "" {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		if claims.AccountId != ownerIdDb.Id {
			ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
			return
		}

		err = apifn.graphService.DeleteRole(params.RoleId, params.WorkspaceName)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		JsonResponseWriter(res, map[string]any{}, http.StatusOK)
	}

	if req.Method == http.MethodPatch {
		var params struct {
			WorkspaceName string `json:"workspaceName"`
			RoleId        string `json:"id"`
			RoleName      string `json:"name"`
			RoleDesc      string `json:"description"`
			CanRead       bool   `json:"canRead"`
			CanCreate     bool   `json:"canCreate"`
			CanRename     bool   `json:"canRename"`
			CanDelete     bool   `json:"canDelete"`
		}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)

		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
			return
		}

		// Get old role details
		role, err := apifn.graphService.GetRole(params.WorkspaceName, params.RoleId)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		if params.RoleName != "" {
			role.Name = params.RoleName
		}
		if params.RoleDesc != "" {
			role.Description = params.RoleDesc
		}
		role.CanCreate = params.CanCreate
		role.CanRead = params.CanRead
		role.CanRename = params.CanRename
		role.CanDelete = params.CanDelete

		err = apifn.graphService.UpdateRole(role)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		resData := make(map[string]any)
		resData["role"] = role
		JsonResponseWriter(res, resData, http.StatusOK)
	}
}

func (apifn ApiConfig) HandleAssignRoleToSA(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	var params struct {
		WorkspaceName    string `json:"workspaceName"`
		RoleId           string `json:"roleId"`
		ServiceAccountId string `json:"accountId"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	if params.WorkspaceName == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidWorkspaceName, http.StatusBadRequest)
		return
	}
	if params.RoleId == "" || params.ServiceAccountId == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(params.WorkspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if claims.AccountId != ownerIdDb.Id {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
		return
	}

	if req.Method == http.MethodPost {
		// Check whether role exists
		roleAssignment, err := apifn.graphService.ChcekRoleSAAssignment(params.RoleId, params.ServiceAccountId)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if roleAssignment {
			ErrorResponseWriter(res, apierrors.ResErrRoleAlreadyAssigned, http.StatusBadRequest)
			return
		}

		err = apifn.graphService.AssignRoleToServiceAccount(params.RoleId, params.ServiceAccountId)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

		resData := make(map[string]any)
		resData["success"] = true
		JsonResponseWriter(res, resData, http.StatusOK)
	}
	if req.Method == http.MethodDelete {
		err = apifn.graphService.RemoveRoleFromServiceAccount(params.RoleId, params.ServiceAccountId)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		JsonResponseWriter(res, map[string]any{}, http.StatusOK)
	}
}

func (apifn ApiConfig) HandleGetRoleFSPermissions(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	var params struct {
		Location string `json:"location"`
		RoleId   string `json:"roleId"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)

	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	if params.Location == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}
	if params.RoleId == "" {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	locationSplit := strings.Split(params.Location, "/")
	workspaceName := locationSplit[0]

	ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if claims.AccountId != ownerIdDb.Id {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
		return
	}

	// Getting the role and its details
	role, err := apifn.graphService.GetRole(workspaceName, params.RoleId)
	if err != nil {
		log.Default().Println(err.Error())
		if errors.Is(err, apierrors.RoleNotFound{}) {
			ErrorResponseWriter(res, apierrors.ResErrRoleNotFound, http.StatusBadRequest)
			return
		}
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if req.Method == http.MethodPost {
		// Getting roles
		roleAssignment, err := apifn.graphService.CheckRoleFSAssignment(params.RoleId, params.Location)
		if err != nil {
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if roleAssignment {
			ErrorResponseWriter(res, apierrors.ResErrRoleAlreadyAssigned, http.StatusBadRequest)
			return
		}
		err = apifn.graphService.AddRoleToItem(role.Id, params.Location)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		JsonResponseWriter(res, map[string]any{}, http.StatusCreated)
		return
	}
	if req.Method == http.MethodDelete {
		err = apifn.graphService.RemoveRoleFromItem(params.RoleId, params.Location)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		JsonResponseWriter(res, map[string]any{}, http.StatusOK)
	}
}

func (apifn ApiConfig) HandleGetAllAccountRoles(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if req.Method != http.MethodGet {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	accountId := req.URL.Query().Get("accountId")
	if len(accountId) == 0 {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	workspaceName := req.URL.Query().Get("workspaceName")
	if len(workspaceName) == 0 {
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	ownerIdDb, err := apifn.graphService.GetWorkspaceOwner(workspaceName)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if ownerIdDb.Id != claims.AccountId {
		ErrorResponseWriter(res, apierrors.ResErrPermissionDenied, http.StatusUnauthorized)
		return
	}

	roles, err := apifn.graphService.GetAllRolesAccount(accountId)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	resData := make(map[string]any)
	resData["roles"] = roles
	JsonResponseWriter(res, resData, http.StatusOK)
}
