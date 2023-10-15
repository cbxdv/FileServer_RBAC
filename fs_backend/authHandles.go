package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"fs_backend/apierrors"
	"fs_backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apifn *ApiConfig) HandleCheckAuth(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	// Enforcing only GET method
	if req.Method != http.MethodGet {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	resData := make(map[string]any)
	resData["account"] = map[string]any{
		"id":       claims.AccountId,
		"name":     claims.Name,
		"username": claims.Username,
		"isOwner":  claims.IsOwner,
	}
	JsonResponseWriter(res, resData, http.StatusOK)
}

func (apifn *ApiConfig) HandleOwnerAccountRegistration(res http.ResponseWriter, req *http.Request) {
	// Enforcing only POST method
	if req.Method != http.MethodPost {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Extracting params from body
	var params struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	// Checking whether the data is valid
	if params.Name == "" || params.Email == "" || params.Password == "" {
		log.Default().Println("Invalid request data")
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	// Checking whether an account already exists
	oldAccount, err := apifn.graphService.GetOwnerAccountWithEmail(params.Email)
	if err != nil {
		if errors.Is(err, apierrors.AccountNotFound{}) {
			log.Default().Println(err.Error(), "- new account")
		} else {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}

	}
	if oldAccount.Email == params.Email {
		log.Default().Println("Account already exists")
		ErrorResponseWriter(res, apierrors.ResErrOAAlreadyExists, http.StatusBadRequest)
		return
	}

	// Enforcing passwords with minimum length of 8
	if len(params.Password) < 8 {
		log.Default().Println("Weak password")
		ErrorResponseWriter(res, apierrors.ResErrWeakPassword, http.StatusBadRequest)
		return
	}

	// Generaing ID with max attempts of 10
	accId := ""
	for i := 1; i < 10; i++ {
		accId = uuid.New().String()
		isPresent, err := apifn.graphService.CheckOwnerAccountWithId(accId)
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
			return
		}
		if !isPresent {
			break
		}
	}

	// Creating new account
	newAccount := models.OwnerAccount{
		Id:       accId,
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}

	// Hashing password
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}
	newAccount.Password = string(hash)

	// Inserting data into db
	err = apifn.graphService.CreateOwnerAccount(newAccount)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	// Responding with a OK message
	resData := make(map[string]any)
	JsonResponseWriter(res, resData, http.StatusCreated)
	log.Default().Println("Registered account with id", newAccount.Id)
}

func (apifn ApiConfig) HandleOwnerAccountLogin(res http.ResponseWriter, req *http.Request) {
	// Enforcing only POST method
	if req.Method != http.MethodPost {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Extracting params from body
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	// Extracting account data from database
	account, err := apifn.graphService.GetOwnerAccountWithEmail(params.Email)
	if err != nil {
		log.Default().Println(err.Error())
		if errors.Is(err, apierrors.AccountNotFound{}) {
			ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusUnauthorized)
			return
		}
	}

	// Checking password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password))
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	// Generating JWT token
	tokenString, err := apifn.generateJwtToken(models.JWTData{
		Issuer:     "fs_backend",
		RemoteAddr: req.RemoteAddr,
		TokenId:    uuid.New().String(),
		AccountId:  account.Id,
		Name:       account.Name,
		Username:   account.Email,
		IsOwner:    true,
	})
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}
	resData := make(map[string]any)
	resData["account"] = map[string]any{
		"id":       account.Id,
		"name":     account.Name,
		"username": account.Email,
		"isOwner":  true,
	}
	resData["token"] = tokenString
	JsonResponseWriter(res, resData, http.StatusOK)
	log.Default().Println("Login by account with id", account.Id)
}

func (apifn ApiConfig) HandleServiceAccountLogin(res http.ResponseWriter, req *http.Request) {
	// Enforcing only POST method
	if req.Method != http.MethodPost {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Extracting params from body
	var params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	splits := strings.Split(params.Username, "@")
	if err != nil {
		ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusForbidden)
		return
	}
	username := splits[0]
	workspaceId := splits[1]

	// Extracting account data from database
	account, workspace, err := apifn.graphService.GetServiceAccountWithWorkspace(username, workspaceId)
	if err != nil {
		log.Default().Println(err.Error())
		if err != nil {
			log.Default().Println(err.Error())
			ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusForbidden)
			return
		}
	}

	// Checking password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password))
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusForbidden)
		return
	}

	// Generating JWT token
	tokenString, err := apifn.generateJwtToken(models.JWTData{
		Issuer:     "fs_backend",
		RemoteAddr: req.RemoteAddr,
		TokenId:    uuid.New().String(),
		AccountId:  account.Id,
		Name:       account.Name,
		Username:   params.Username,
		IsOwner:    false,
	})
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	resData := make(map[string]any)
	resData["token"] = tokenString
	resData["account"] = map[string]any{
		"id":        account.Id,
		"name":      account.Name,
		"username":  account.Username,
		"isOwner":   false,
		"workspace": workspace.Name,
	}
	JsonResponseWriter(res, resData, http.StatusOK)
	log.Default().Println("Login by account with id", account.Id)
}

func (apifn ApiConfig) HandleOwnerAccountChangePassword(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if req.Method != http.MethodPatch {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Extracting params from body
	var params struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidData, http.StatusBadRequest)
		return
	}

	// Extracting account data from database
	account, err := apifn.graphService.GetOwnerAccountWithId(claims.AccountId)
	if err != nil {
		log.Default().Println(err.Error())
		if errors.Is(err, apierrors.AccountNotFound{}) {
			ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusForbidden)
			return
		}
	}

	// Checking password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.OldPassword))
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrInvalidCredentials, http.StatusForbidden)
		return
	}

	// Hashing password
	hash, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), 14)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	// Updating password
	err = apifn.graphService.UpdateOwnerAccountPassword(claims.AccountId, string(hash))
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}
	resp := make(map[string]any)
	JsonResponseWriter(res, resp, http.StatusOK)
}

func (apifn ApiConfig) HandleOwnerAccountDelete(res http.ResponseWriter, req *http.Request, claims models.JWTData) {
	if req.Method != http.MethodDelete {
		ErrorResponseWriter(res, apierrors.ResErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Check whether any workspace exists for the account
	workspaces, err := apifn.graphService.GetWorkspacesForOwner(claims.AccountId)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	if len(workspaces) != 0 {
		ErrorResponseWriter(res, apierrors.ResErrAccountHasWorkspace, http.StatusBadRequest)
		return
	}

	// Delete account
	err = apifn.graphService.DeleteOwnerAccount(claims.AccountId)
	if err != nil {
		log.Default().Println(err.Error())
		ErrorResponseWriter(res, apierrors.ResErrServerError, http.StatusInternalServerError)
		return
	}

	// Respond with OK
	resData := make(map[string]any)
	JsonResponseWriter(res, resData, http.StatusOK)
	log.Default().Println("Account with id", claims.AccountId, "deleted")
}
