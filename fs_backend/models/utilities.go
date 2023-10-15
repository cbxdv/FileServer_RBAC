package models

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetOnwerAccountFromRecord(record any) OwnerAccount {
	att := record.(neo4j.Node).Props
	return OwnerAccount{
		Id:       att["id"].(string),
		Name:     att["name"].(string),
		Email:    att["email"].(string),
		Password: att["password"].(string),
	}
}

func GetServiceAccountFromRecord(record any) ServiceAccount {
	att := record.(neo4j.Node).Props
	return ServiceAccount{
		Id:                  att["id"].(string),
		Username:            att["username"].(string),
		Name:                att["name"].(string),
		LinkedEmail:         att["linkedEmail"].(string),
		Password:            att["password"].(string),
		ShouldResetPassword: att["shouldResetPassword"].(bool),
	}
}

func GetDirectoryFromRecord(record any) Directory {
	att := record.(neo4j.Node).Props
	return Directory{
		Id:        att["id"].(string),
		Type:      "directory",
		Name:      att["name"].(string),
		Location:  att["location"].(string),
		CreatedOn: att["createdOn"].(time.Time),
	}
}

func GetFileFromRecord(record any) File {
	att := record.(neo4j.Node).Props
	return File{
		Id:        att["id"].(string),
		Type:      "file",
		Name:      att["name"].(string),
		Size:      int(att["size"].(int64)),
		Location:  att["location"].(string),
		CreatedOn: att["createdOn"].(time.Time),
	}
}

func GetRoleFromRecord(record any) Role {
	att := record.(neo4j.Node).Props
	return Role{
		Id:          att["id"].(string),
		Name:        att["name"].(string),
		Description: att["description"].(string),
		CanRead:     att["canRead"].(bool),
		CanCreate:   att["canCreate"].(bool),
		CanRename:   att["canRename"].(bool),
		CanDelete:   att["canDelete"].(bool),
	}
}

func GetWorkspaceFromRecord(record any) Workspace {
	att := record.(neo4j.Node).Props
	return Workspace{
		Id:   att["id"].(string),
		Name: att["name"].(string),
	}
}
