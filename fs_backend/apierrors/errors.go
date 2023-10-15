package apierrors

import "fmt"

/* ------------------------------- JWT Errors ------------------------------- */

type NoJWTToken struct{}

func (NoJWTToken) Error() string {
	return "No token found in request"
}

type InvalidToken struct{}

func (InvalidToken) Error() string {
	return "Invalid Token"
}

/* ----------------------------- Account Errors ----------------------------- */

type AccountWithEmailAlreadyExists struct {
	Email string
}

func (uae AccountWithEmailAlreadyExists) Error() string {
	return "User with email " + uae.Email + " already exists"
}

type AccountNotFound struct{}

func (AccountNotFound) Error() string {
	return "User not found"
}

/* ---------------------------- Workspace Errors ---------------------------- */

type WorkspaceAlreadyExists struct{}

func (WorkspaceAlreadyExists) Error() string {
	return "Workspace already exists"
}

type WorkspaceNotFound struct{}

func (WorkspaceNotFound) Error() string {
	return "Workspace not found"
}

/* ---------------------------- Directory Errors ---------------------------- */

type DirectoryNotFound struct {
	DirName string
}

func (err DirectoryNotFound) Error() string {
	if err.DirName == "" {
		return "Directory not found"
	}
	return "Directory with name " + err.DirName + " not found"
}

type DirectoryWithSameNameAlreadyExists struct {
	ParentDirName string
	DirName       string
}

func (err DirectoryWithSameNameAlreadyExists) Error() string {
	return "Directory with name " + err.DirName + " already exists in " + err.ParentDirName
}

/* ------------------------------- File Errors ------------------------------ */

type FileNotFound struct {
	FileName string
}

func (err FileNotFound) Error() string {
	if err.FileName == "" {
		return "File not found"
	}
	return "File with name " + err.FileName + " not found"
}

type FileWithSameNameAlreadyExists struct {
	ParentDirName string
	FileName      string
}

func (err FileWithSameNameAlreadyExists) Error() string {
	return "File with name " + err.FileName + " already exists in directory " + err.ParentDirName
}

/* ------------------------ Upload Properties Errors ------------------------ */

type UploadIdNotFound struct {
	UploadId string
}

func (err UploadIdNotFound) Error() string {
	return fmt.Sprintf("UploadId %s not found", err.UploadId)
}

/* ------------------------------- Role Errors ------------------------------ */

type RoleNotFound struct {
	RoleId string
}

func (err RoleNotFound) Error() string {
	return fmt.Sprintf("Role with id %s not found", err.RoleId)
}
