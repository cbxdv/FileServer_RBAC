package models

import "time"

type Workspace struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OwnerAccount struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ServiceAccount struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Username            string `json:"username"`
	LinkedEmail         string `json:"linkedEmail"`
	ShouldResetPassword bool   `json:"shouldResetPassword"`
	Password            string `json:"password"`
}

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CanRead     bool   `json:"canRead"`
	CanCreate   bool   `json:"canCreate"`
	CanRename   bool   `json:"canRename"`
	CanDelete   bool   `json:"canDelete"`
}

type Directory struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedOn time.Time `json:"createdOn"`
}

type File struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	Location  string    `json:"location"`
	CreatedOn time.Time `json:"createdOn"`
}

type DirectoryWithContents struct {
	Id        string        `json:"id"`
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	CreatedOn time.Time     `json:"createdOn"`
	Location  string        `json:"location"`
	Contents  []interface{} `json:"contents"`
}

type RoleWithUsers struct {
	Role            Role             `json:"role"`
	ServiceAccounts []ServiceAccount `json:"accounts"`
}

type FileTransferProperties struct {
	FileProperties File
	LinkId         string
	LinkGenerated  time.Time
}

type JWTData struct {
	Issuer     string
	RemoteAddr string
	TokenId    string
	AccountId  string
	Name       string
	Username   string
	IsOwner    bool
}
