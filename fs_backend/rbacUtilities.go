package main

import "fs_backend/models"

var resolutionMap map[bool]map[bool]bool = map[bool]map[bool]bool{
	true: {
		true:  true,
		false: true,
	},
	false: {
		true:  true,
		false: false,
	},
}

func resolveRoles(roles []models.Role) models.Role {
	if len(roles) == 1 {
		return roles[0]
	}
	newRole := roles[0]
	for _, role := range roles[1:] {
		newRole.CanRead = resolutionMap[newRole.CanRead][role.CanRead]
		newRole.CanCreate = resolutionMap[newRole.CanCreate][role.CanCreate]
		newRole.CanRename = resolutionMap[newRole.CanRename][role.CanRename]
		newRole.CanDelete = resolutionMap[newRole.CanDelete][role.CanDelete]
	}
	return newRole
}
