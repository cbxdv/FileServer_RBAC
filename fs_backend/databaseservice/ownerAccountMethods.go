package databaseservice

import (
	"errors"
	"log"

	"fs_backend/apierrors"
	"fs_backend/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (gds *GraphDatabaseService) GetOwnerAccountWithEmail(email string) (models.OwnerAccount, error) {
	if email == "" {
		return models.OwnerAccount{}, errors.New("invalid email")
	}
	getAccountCypher := `
		MATCH (account:OwnerAccount) WHERE account.email = $email
		RETURN account
	`
	getAccountCypherParams := map[string]any{
		"email": email,
	}
	getAccountResponse, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getAccountCypher, getAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.OwnerAccount{}, err
	}
	if len(getAccountResponse.Records) == 0 {
		return models.OwnerAccount{}, apierrors.AccountNotFound{}
	}
	accRecord, found := getAccountResponse.Records[0].Get("account")
	if !found {
		return models.OwnerAccount{}, apierrors.AccountNotFound{}
	}
	user := models.GetOnwerAccountFromRecord(accRecord)
	return user, nil
}

func (gds *GraphDatabaseService) GetOwnerAccountWithId(accountId string) (models.OwnerAccount, error) {
	if accountId == "" {
		return models.OwnerAccount{}, errors.New("invalid id")
	}
	getAccountCypher := `
		MATCH (account:OwnerAccount) WHERE account.id = $accountId
		RETURN account
	`
	getAccountCypherParams := map[string]any{
		"accountId": accountId,
	}
	getAccountResponse, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getAccountCypher, getAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.OwnerAccount{}, err
	}
	if len(getAccountResponse.Records) == 0 {
		return models.OwnerAccount{}, apierrors.AccountNotFound{}
	}
	accRecord, found := getAccountResponse.Records[0].Get("account")
	if !found {
		return models.OwnerAccount{}, apierrors.AccountNotFound{}
	}
	user := models.GetOnwerAccountFromRecord(accRecord)
	return user, nil
}

func (gds GraphDatabaseService) CheckOwnerAccountWithId(accountId string) (bool, error) {
	if accountId == "" {
		return false, errors.New("invalid id")
	}
	checkAccountCypher := `
		MATCH (a:OwnerAccount) WHERE a.id = $accountId
		RETURN count(a) as count
	`
	checkAccountCypherParams := map[string]any{
		"accountId": accountId,
	}
	if gds.driver == nil {
		log.Fatalln("Driver is nil")
	}
	if gds.ctx == nil {
		log.Fatalln("Context is nil")
	}

	checkAccountResponse, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		checkAccountCypher, checkAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return false, err
	}
	accountCount, found := checkAccountResponse.Records[0].Get("count")
	if !found {
		return false, nil
	}
	if accountCount.(int64) != 1 {
		return false, nil
	}
	return true, nil
}

func (gds *GraphDatabaseService) CreateOwnerAccount(account models.OwnerAccount) error {
	if account.Name == "" || account.Email == "" || account.Password == "" {
		return errors.New("invalid user details")
	}

	// Checking whether an account already exists with the email provided
	accountResponse, err := gds.GetOwnerAccountWithEmail(account.Email)
	if err != nil {
		log.Default().Println(err.Error(), ": creating new user")
		if errors.Is(err, apierrors.AccountNotFound{}) {
			// Do nothing as no account with the same email exists
		} else {
			return err
		}
	}
	if accountResponse.Email == account.Email {
		return apierrors.AccountWithEmailAlreadyExists{Email: account.Email}
	}

	createOACypher := `
		CREATE (:OwnerAccount{
			id:       $id,
			name:     $name,
			email:    $email,
			password: $password
		})
	`
	createOACypherParams := map[string]any{
		"id":       account.Id,
		"name":     account.Name,
		"email":    account.Email,
		"password": account.Password,
	}
	_, err = neo4j.ExecuteQuery(gds.ctx, gds.driver,
		createOACypher, createOACypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds *GraphDatabaseService) UpdateOwnerAccountPassword(accountId string, newPassword string) error {
	if accountId == "" || newPassword == "" {
		return errors.New("invalid id or password")
	}
	updateAccountCypher := `
		MATCH (a:OwnerAccount) WHERE a.id = $accountId
		SET a.password = $newPassword
	`
	updateAccountCypherParams := map[string]any{
		"accountId":   accountId,
		"newPassword": newPassword,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		updateAccountCypher, updateAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds *GraphDatabaseService) DeleteOwnerAccount(accountId string) error {
	if accountId == "" {
		return errors.New("invalid id")
	}
	deleteAccountCypher := `
		MATCH (a:OwnerAccount) WHERE a.id = $accountId
		DETACH DELETE a
	`
	deleteAccountCypherParams := map[string]any{
		"accountId": accountId,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		deleteAccountCypher, deleteAccountCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}
