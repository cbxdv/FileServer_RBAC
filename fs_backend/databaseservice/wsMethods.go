package databaseservice

import (
	"fs_backend/apierrors"
	"fs_backend/models"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (gds *GraphDatabaseService) CheckWorkspace(workspaceName string) (bool, error) {
	checkWorkspaceAvailabilityCypher := `
		MATCH (w:Workspace) WHERE w.name = $workspaceName
		RETURN count(w) as count
	`
	checkWorkspaceAvailabilityParams := map[string]any{
		"workspaceName": workspaceName,
	}
	checkWorkspaceAvailabilityRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		checkWorkspaceAvailabilityCypher, checkWorkspaceAvailabilityParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return false, err
	}
	workspaceCount, found := checkWorkspaceAvailabilityRes.Records[0].Get("count")
	if !found || workspaceCount.(int64) == 0 {
		return false, nil
	}
	return true, nil
}

func (gds GraphDatabaseService) CreateWorkspace(workspace models.Workspace, ownerAccountId string) error {
	// Check if workspace already exists
	if exists, err := gds.CheckWorkspace(workspace.Name); err != nil {
		log.Default().Println(err.Error())
		return err
	} else if exists {
		return apierrors.WorkspaceAlreadyExists{}
	}

	createWorkspaceCypher := `
		MATCH (owner:OwnerAccount) WHERE owner.id = $ownerAccountId
		CREATE (owner)-[:OWNS]->(w:Workspace {
			id: $workspaceId,
			name: $workspaceName
		})-[:STORES]->(wd:RootDirectory:Directory{
			id: $workspaceId,
			type: "directory",
			name: $workspaceName,
			createdOn: datetime(),
			location: $workspaceName
		})
	`
	createWorkspaceParams := map[string]any{
		"ownerAccountId": ownerAccountId,
		"workspaceId":    workspace.Id,
		"workspaceName":  workspace.Name,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		createWorkspaceCypher, createWorkspaceParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) DeleteWorkspace(workspaceName string) error {
	// Deleting all roles and cypher
	deleteWorkspaceCypher := `
		MATCH (w:Workspace) WHERE w.name = $workspaceName
		OPTIONAL MATCH p1=(w)-[*]->(c)
		OPTIONAL MATCH p2=(r:Role)-[*]->(w)
		OPTIONAL MATCH p3=(s:ServiceAccount)-[*]->(w)
		DETACH DELETE p3, p2, p1, w
	`
	deleteWorkspaceParams := map[string]any{
		"workspaceName": workspaceName,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		deleteWorkspaceCypher, deleteWorkspaceParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) UpdateWorkspaceName(oldWorkspaceName string, newWorkspaceName string) error {
	// Check if workspace already exists
	if exists, err := gds.CheckWorkspace(oldWorkspaceName); err != nil {
		log.Default().Println(err.Error())
		return err
	} else if !exists {
		return apierrors.WorkspaceNotFound{}
	}

	// Check if new workspace name is available
	if exists, err := gds.CheckWorkspace(newWorkspaceName); err != nil {
		log.Default().Println(err.Error())
		return err
	} else if exists {
		return apierrors.WorkspaceAlreadyExists{}
	}

	updateWorkspaceNameCypher := `
		MATCH (w:Workspace)-[:STORES]->(rd:RootDirectory) WHERE w.name = $workspaceName
		SET w.name = $newWorkspaceName, rd.name = $newWorkspaceName, rd.location = $newWorkspaceName
	`
	updateWorkspaceNameParams := map[string]any{
		"workspaceName":    oldWorkspaceName,
		"newWorkspaceName": newWorkspaceName,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		updateWorkspaceNameCypher, updateWorkspaceNameParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) GetWorkspaceOwner(workspaceName string) (models.OwnerAccount, error) {
	getWorkspaceOwnerCypher := `
		MATCH (owner:OwnerAccount)-[:OWNS]->(w:Workspace)
		WHERE w.name = $workspaceName
		RETURN owner
	`
	getWorkspaceOwnerParams := map[string]any{
		"workspaceName": workspaceName,
	}
	getWorkspaceOwnerRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getWorkspaceOwnerCypher, getWorkspaceOwnerParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.OwnerAccount{}, err
	}

	if len(getWorkspaceOwnerRes.Records) == 0 {
		return models.OwnerAccount{}, apierrors.WorkspaceNotFound{}
	}

	ownerRecord, found := getWorkspaceOwnerRes.Records[0].Get("owner")
	if !found {
		return models.OwnerAccount{}, apierrors.WorkspaceNotFound{}
	}
	owner := models.GetOnwerAccountFromRecord(ownerRecord)
	return owner, nil
}

func (gds GraphDatabaseService) GetServiceAccountWithWorkspace(username string, workspaceName string) (models.ServiceAccount, models.Workspace, error) {
	getAccountCypher := `
		MATCH (sa:ServiceAccount)-[:SERVICES]->(w:Workspace)
		WHERE sa.username = $username AND w.name = $workspaceName
		RETURN sa, w
	`
	getAccountCypherParams := map[string]any{
		"username":      username,
		"workspaceName": workspaceName,
	}
	getAccountRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getAccountCypher, getAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		return models.ServiceAccount{}, models.Workspace{}, err
	}

	if len(getAccountRes.Records) == 0 {
		return models.ServiceAccount{}, models.Workspace{}, apierrors.AccountNotFound{}
	}

	accountRecord, found := getAccountRes.Records[0].Get("sa")
	if !found {
		return models.ServiceAccount{}, models.Workspace{}, apierrors.AccountNotFound{}
	}
	account := models.GetServiceAccountFromRecord(accountRecord)

	workspaceRecord, found := getAccountRes.Records[0].Get("w")
	if !found {
		return models.ServiceAccount{}, models.Workspace{}, apierrors.WorkspaceNotFound{}
	}
	workspace := models.GetWorkspaceFromRecord(workspaceRecord)

	return account, workspace, nil
}

func (gds GraphDatabaseService) GetAllServiceAccountsInWorkspace(workspaceName string) ([]models.ServiceAccount, error) {
	getAccountsCypher := `
		MATCH (sa:ServiceAccount)-[:SERVICES]->(w:Workspace)
		WHERE w.name = $workspaceName
		RETURN COLLECT(sa) AS accounts
	`
	getAccountsCypherParams := map[string]any{
		"workspaceName": workspaceName,
	}
	getAccountsRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getAccountsCypher, getAccountsCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		return nil, err
	}
	accountsRecords, found := getAccountsRes.Records[0].Get("accounts")
	if !found {
		return make([]models.ServiceAccount, 0), nil
	}
	accounts := []models.ServiceAccount{}
	recordList := accountsRecords.([]any)
	for _, accRecord := range recordList {
		accounts = append(accounts, models.GetServiceAccountFromRecord(accRecord))
	}
	return accounts, nil
}

func (gds GraphDatabaseService) CreateServiceAccount(account models.ServiceAccount, workspaceName string) error {
	// Check if workspace exists
	if exists, err := gds.CheckWorkspace(workspaceName); err != nil {
		log.Default().Println(err.Error())
		return err
	} else if !exists {
		return apierrors.WorkspaceNotFound{}
	}

	createServiceAccCypher := `
		MATCH (w:Workspace) WHERE w.name = $workspaceName
		CREATE (su:ServiceAccount {
			id: $accountId,
			name: $name,
			username: $username,
			linkedEmail: $linkedEmail,
			password: $password,
			shouldResetPassword: $shouldResetPassword
		})-[r:SERVICES]->(w)
	`
	createServiceAccCypherParams := map[string]any{
		"workspaceName":       workspaceName,
		"accountId":           account.Id,
		"name":                account.Name,
		"username":            account.Username,
		"password":            account.Password,
		"linkedEmail":         account.LinkedEmail,
		"shouldResetPassword": account.ShouldResetPassword,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		createServiceAccCypher, createServiceAccCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) DeleteServiceAccount(workspaceName string, username string) error {
	deleteSericeAccCypher := `
		MATCH (sa:ServiceAccount)-[r:SERVICES]->(w:Workspace)
		WHERE sa.username = $username AND w.name = $workspaceName
		DETACH DELETE sa
	`
	deleteSericeAccCypherParams := map[string]any{
		"username":      username,
		"workspaceName": workspaceName,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		deleteSericeAccCypher, deleteSericeAccCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) GetWorkspacesForOwner(ownerAccountId string) ([]models.Workspace, error) {
	getWorkspacesCypher := `
		MATCH (owner:OwnerAccount)-[:OWNS]->(w:Workspace)
		WHERE owner.id = $ownerAccountId
		RETURN COLLECT(w) AS workspaces
	`
	getWorkspacesCypherParams := map[string]any{
		"ownerAccountId": ownerAccountId,
	}
	getWorkspacesRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getWorkspacesCypher, getWorkspacesCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return []models.Workspace{}, err
	}
	workspaces := []models.Workspace{}
	workspaceRecords, found := getWorkspacesRes.Records[0].Get("workspaces")
	if found {
		recordList := workspaceRecords.([]any)
		for _, wsRecord := range recordList {
			workspaces = append(workspaces, models.GetWorkspaceFromRecord(wsRecord))
		}
	}

	return workspaces, nil
}
