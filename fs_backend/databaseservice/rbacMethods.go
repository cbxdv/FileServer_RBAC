package databaseservice

import (
	"fs_backend/apierrors"
	"fs_backend/models"
	"log"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (gds GraphDatabaseService) CreateNewRole(role models.Role, workspaceName string) error {
	createNewRoleCypher := `
		MATCH (w:Workspace{name: $workspaceName})
		CREATE (:Role {
			id:          $roleId,
			name:        $roleName,
			description: $roleDescription,
			canRead:     $canRead,
			canCreate:   $canCreate,
			canRename:   $canRename,
			canDelete:   $canDelete
		})-[:ROLLED_IN]->(w)
	`
	createNewRoleParams := map[string]interface{}{
		"workspaceName":   workspaceName,
		"roleId":          role.Id,
		"roleName":        role.Name,
		"roleDescription": role.Description,
		"canRead":         role.CanRead,
		"canCreate":       role.CanCreate,
		"canRename":       role.CanRename,
		"canDelete":       role.CanDelete,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		createNewRoleCypher, createNewRoleParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) DeleteRole(roleId string, workspaceName string) error {
	deleteRoleCypher := `
		MATCH (r:Role)-[:ROLLED_IN]->(w:Workspace)
		WHERE w.name = $workspaceName AND r.id = $roleId
		DETACH DELETE r
	`
	deleteRoleParams := map[string]interface{}{
		"workspaceName": workspaceName,
		"roleId":        roleId,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		deleteRoleCypher, deleteRoleParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) ChcekRoleSAAssignment(roleId string, servieAccountId string) (bool, error) {
	checkRole := `
		MATCH (r:Role{id: $roleId})<-[rr:HAS_ROLE]-(a:ServiceAccount{id: $accountId})
		RETURN count(rr) AS count
	`
	checkRoleParams := map[string]interface{}{
		"roleId":    roleId,
		"accountId": servieAccountId,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		checkRole, checkRoleParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return false, err
	}
	if len(recordsRes.Records) == 0 {
		return false, nil
	}
	count, found := recordsRes.Records[0].Get("count")
	if found {
		return count.(int64) > 0, nil
	}
	return false, nil
}

func (gds GraphDatabaseService) CheckRoleFSAssignment(roleId string, location string) (bool, error) {
	checkRoleCypher := `
		MATCH (r:Role{id: $roleId})-[]->(i:File|Directory{location: $location})
		RETURN count(r) AS count
	`
	checkRoleCypherParams := map[string]interface{}{
		"roleId":   roleId,
		"location": location,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		checkRoleCypher, checkRoleCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return false, err
	}
	if len(recordsRes.Records) == 0 {
		return false, nil
	}
	count, found := recordsRes.Records[0].Get("count")
	if found {
		return count.(int64) > 0, nil
	}
	return false, nil
}

func (gds GraphDatabaseService) AssignRoleToServiceAccount(roleId string, accountId string) error {
	addRoleToAccountCypher := `
		MATCH (r:Role) WHERE r.id = $roleId
		MATCH (s:ServiceAccount) WHERE s.id = $accountId
		CREATE (r)<-[:HAS_ROLE]-(s)
	`
	addRoleToAccountCypherParams := map[string]interface{}{
		"roleId":    roleId,
		"accountId": accountId,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		addRoleToAccountCypher, addRoleToAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) RemoveRoleFromServiceAccount(roleId string, accId string) error {
	removeRoleFromAccountCypher := `
		MATCH (a:ServiceAccount)-[rr:HAS_ROLE]->(r:Role)
		WHERE a.id = $accId AND r.id = $roleId
		DELETE rr
	`
	removeRoleFromAccountParams := map[string]interface{}{
		"accId":  accId,
		"roleId": roleId,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		removeRoleFromAccountCypher, removeRoleFromAccountParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) UpdateRole(updatedRole models.Role) error {
	updateRoleCypher := `
		MATCH (r:Role) WHERE r.id = $roleId
		SET r.name = $roleName, r.description = $roleDesc, r.canRead = $canRead, r.canCreate = $canCreate, r.canRename = $canRename, r.canDelete = $canDelete
	`
	updateRoleCypherParams := map[string]interface{}{
		"roleId":    updatedRole.Id,
		"roleName":  updatedRole.Name,
		"roleDesc":  updatedRole.Description,
		"canRead":   updatedRole.CanRead,
		"canCreate": updatedRole.CanCreate,
		"canRename": updatedRole.CanRename,
		"canDelete": updatedRole.CanDelete,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		updateRoleCypher, updateRoleCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) AddRoleToItem(roleId string, location string) error {
	addRoleCypher := `
		MATCH (r:Role) WHERE r.id = $roleId
		MATCH (i:Directory|File) WHERE i.location = $location
		CREATE (r)-[:MANAGES]->(i)
	`
	addRoleCypherParams := map[string]interface{}{
		"roleId":   roleId,
		"location": location,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		addRoleCypher, addRoleCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) RemoveRoleFromItem(roleId string, location string) error {
	removeRoleCypher := `
		MATCH (r:Role)-[m:MANAGES]->(i:Directory|File)
		WHERE r.id = $roleId AND i.location = $location
		DELETE m
	`
	removeRoleCypherParams := map[string]interface{}{
		"roleId":   roleId,
		"location": location,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		removeRoleCypher, removeRoleCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) GetNearestRole(accountId string, location string) ([]models.Role, error) {
	locationSplit := strings.Split(location, "/")
	workspaceName := locationSplit[0]

	getNearestRolesCypher := ""
	getNearestRolesCypherParams := map[string]any{
		"workspaceName": workspaceName,
		"accountId":     accountId,
		"location":      location,
	}
	if len(locationSplit) == 1 {
		getNearestRolesCypher = `
			MATCH (child:Directory|File{location: $location})
			MATCH (r:Role)-[:MANAGES]->(child)
			MATCH (:ServiceAccount{id: $accountId})-[:HAS_ROLE]->(r)
			RETURN collect(r) AS avaRoles
		`
	} else {
		getNearestRolesCypher = `
			MATCH (root:Directory{location: $workspaceName})
			MATCH (child:Directory|File{location: $location})
			MATCH p1=SHORTESTPATH((root)-[:CONTAINS*]->(child))
				WITH nodes(p1) AS ns, child AS child
			MATCH (r:Role)-[:MANAGES]->(f)
			MATCH (r)<-[:HAS_ROLE]-(a:ServiceAccount{id: $accountId})
				WHERE ANY(n IN [(r)-[:MANAGES]->(f) | f] WHERE n IN ns)
				WITH r AS r, child AS child, MIN(COALESCE(length(SHORTESTPATH((child)<-[*]-(r))))) AS min
				WHERE COALESCE(length(SHORTESTPATH((child)<-[*]-(r)))) = min
			RETURN collect(r) AS avaRoles
		`
	}

	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getNearestRolesCypher, getNearestRolesCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return []models.Role{}, err
	}
	roles := []models.Role{}
	rolesRecords, found := recordsRes.Records[0].Get("avaRoles")
	if found {
		recordList := rolesRecords.([]any)
		for _, roleRecord := range recordList {
			roles = append(roles, models.GetRoleFromRecord(roleRecord))
		}
	}
	return roles, nil
}

func (gds GraphDatabaseService) GetAllRolesInWorkspace(workspaceName string) ([]models.Role, error) {
	getRolesCypher := `
		MATCH (w:Workspace{name: $workspaceName})<-[:ROLLED_IN]-(r:Role)
		RETURN collect(r) AS r
	`
	getRolesCypherParams := map[string]interface{}{
		"workspaceName": workspaceName,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getRolesCypher, getRolesCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return []models.Role{}, err
	}
	roles := []models.Role{}
	rolesRecords, found := recordsRes.Records[0].Get("r")
	if found {
		recordList := rolesRecords.([]any)
		for _, roleRecord := range recordList {
			roles = append(roles, models.GetRoleFromRecord(roleRecord))
		}
	}
	return roles, nil
}

func (gds GraphDatabaseService) GetRole(workspaceName string, roleId string) (models.Role, error) {
	getRoleCypher := `
		MATCH (r:Role{id: $roleId})-[:ROLLED_IN]->(w:Workspace{name: $workspaceName})
		RETURN r
	`
	getRoleCypherParams := map[string]interface{}{
		"workspaceName": workspaceName,
		"roleId":        roleId,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getRoleCypher, getRoleCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.Role{}, err
	}
	if len(recordsRes.Records) == 0 {
		return models.Role{}, apierrors.RoleNotFound{}
	}
	roleRecord, found := recordsRes.Records[0].Get("r")
	if !found {
		return models.Role{}, apierrors.RoleNotFound{}
	}
	return models.GetRoleFromRecord(roleRecord), nil
}

func (gds GraphDatabaseService) GetRoleDetailsWithSAInWorkspace(workspaceName string, roleId string) (models.RoleWithUsers, error) {
	getRoleDetails := `
		MATCH (w:Workspace{name: $workspaceName})<-[:ROLLED_IN]-(r:Role{id: $roleId})
		OPTIONAL MATCH (r)<-[:HAS_ROLE]-(a:ServiceAccount)
		RETURN r, collect(a) AS users
	`
	getRoleDetailsParams := map[string]interface{}{
		"workspaceName": workspaceName,
		"roleId":        roleId,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getRoleDetails, getRoleDetailsParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.RoleWithUsers{}, err
	}
	if len(recordsRes.Records) == 0 {
		return models.RoleWithUsers{}, apierrors.RoleNotFound{}
	}
	roleWithUsers := models.RoleWithUsers{
		Role:            models.Role{},
		ServiceAccounts: []models.ServiceAccount{},
	}
	roleRecord, found := recordsRes.Records[0].Get("r")
	if found {
		roleWithUsers.Role = models.GetRoleFromRecord(roleRecord)
	}
	usersRecords, found := recordsRes.Records[0].Get("users")
	if found {
		recordList := usersRecords.([]any)
		for _, userRecord := range recordList {
			roleWithUsers.ServiceAccounts = append(roleWithUsers.ServiceAccounts, models.GetServiceAccountFromRecord(userRecord))
		}
	}
	return roleWithUsers, nil
}

func (gds GraphDatabaseService) GetItemAllRoles(location string) ([]models.Role, error) {
	getAllRolesCypher := `
		MATCH (r:Role)-[:MANAGES]->(i:File|Directory)
		WHERE i.location=$location
		RETURN COLLECT(r) as roles
	`
	getAllRolesCypherParams := map[string]interface{}{
		"location": location,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getAllRolesCypher, getAllRolesCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return []models.Role{}, err
	}
	roles := []models.Role{}
	rolesRecord, found := recordsRes.Records[0].Get("roles")
	if found {
		rolesList := rolesRecord.([]any)
		for _, roleRecord := range rolesList {
			roles = append(roles, models.GetRoleFromRecord(roleRecord))
		}
	}
	return roles, nil
}

func (gds GraphDatabaseService) GetAllRolesAccount(accountId string) ([]models.Role, error) {
	getRolesCypher := `
		MATCH (sa:ServiceAccount)-[:HAS_ROLE]->(r:Role)
		WHERE sa.id=$accountId
		RETURN COLLECT(r) AS roles
	`
	getRolesCypherParams := map[string]interface{}{
		"accountId": accountId,
	}
	recordsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getRolesCypher, getRolesCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return []models.Role{}, err
	}
	roles := []models.Role{}
	rolesRecord, found := recordsRes.Records[0].Get("roles")
	if found {
		rolesList := rolesRecord.([]any)
		for _, roleRecord := range rolesList {
			roles = append(roles, models.GetRoleFromRecord(roleRecord))
		}
	}
	return roles, nil
}
