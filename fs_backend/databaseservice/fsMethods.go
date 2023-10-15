package databaseservice

import (
	"log"
	"strings"

	"fs_backend/apierrors"
	"fs_backend/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (gds *GraphDatabaseService) checkDirExistence(location string) (bool, error) {
	checkDirCypher := `
		MATCH (d:Directory) WHERE d.location = $location
		RETURN count(d) as count
	`
	checkDirCypherParams := map[string]any{
		"location": location,
	}
	checkDirRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		checkDirCypher, checkDirCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return false, err
	}
	dirCount, found := checkDirRes.Records[0].Get("count")
	if !found || dirCount.(int64) != 1 {
		return false, nil
	}
	return true, nil
}

func (gds *GraphDatabaseService) checkForSameFileName(parentLocation string, fileName string) (bool, error) {
	checkFileCypher := `
		MATCH (d:Directory)-[:CONTAINS]->(c)
		WHERE d.location = $parentLocation AND c.Name = $fileName
		RETURN count(c) as count
	`
	checkFileParams := map[string]any{
		"parentLocation": parentLocation,
		"fileName":       fileName,
	}
	checkFileRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		checkFileCypher, checkFileParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return false, err
	}
	contentsCount, found := checkFileRes.Records[0].Get("count")
	if !found || contentsCount.(int64) != 1 {
		return false, nil
	}
	return true, nil
}

func (gds *GraphDatabaseService) CreateDirectory(directory models.Directory) error {
	locationSplit := strings.Split(directory.Location, "/")
	parentlocation := strings.Join(locationSplit[:len(locationSplit)-1], "/")

	// Check whether folder already exists
	dirExistence, err := gds.checkDirExistence(directory.Location)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	if dirExistence {
		return apierrors.DirectoryWithSameNameAlreadyExists{
			ParentDirName: parentlocation,
			DirName:       directory.Name,
		}
	}

	createDirCypher := `
		MATCH (d:Directory) WHERE d.location = $parentLocation
		CREATE (d)-[:CONTAINS]->(newDir:Directory {
			id: $id,
			type: "directory",
			name: $dirName,
			location: $location,
			createdOn: $createdOn
		})
	`
	createDirCypherParams := map[string]any{
		"parentLocation": parentlocation,
		"id":             directory.Id,
		"dirName":        directory.Name,
		"location":       directory.Location,
		"createdOn":      directory.CreatedOn,
	}
	_, err = neo4j.ExecuteQuery(gds.ctx, gds.driver,
		createDirCypher, createDirCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) CreateFile(file models.File) error {
	locationSplit := strings.Split(file.Location, "/")
	parentLocation := strings.Join(locationSplit[:len(locationSplit)-1], "/")

	// Check if parent directory exists
	if dirExistence, err := gds.checkDirExistence(parentLocation); err != nil {
		log.Default().Println(err.Error())
		return err
	} else if !dirExistence {
		return apierrors.DirectoryNotFound{}
	}

	// Checking whether no same names in directory
	if sameFileNameExists, err := gds.checkForSameFileName(parentLocation, file.Name); err != nil {
		log.Default().Println(err.Error())
		return err
	} else if sameFileNameExists {
		return apierrors.FileWithSameNameAlreadyExists{}
	}

	// Create file
	createFileCypher := `
		MATCH (d:Directory) WHERE d.location = $parentLocation
		CREATE (d)-[:CONTAINS]->(newFile:File {
			id: $id,
			type: "file",
			name: $name,
			size: $size,
			location: $location,
			createdOn: $createdOn
		})
	`
	createFileParams := map[string]any{
		"parentLocation": parentLocation,
		"id":             file.Id,
		"name":           file.Name,
		"size":           file.Size,
		"location":       file.Location,
		"createdOn":      file.CreatedOn,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		createFileCypher, createFileParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) GetDirectoryDetails(location string) (models.Directory, error) {
	getDirCypher := `
		MATCH (d:Directory) WHERE d.location = $location
		RETURN d
	`
	getDirCypherParams := map[string]any{
		"location": location,
	}
	getDirRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getDirCypher, getDirCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.Directory{}, err
	}
	if len(getDirRes.Records) == 0 {
		return models.Directory{}, apierrors.DirectoryNotFound{}
	}
	dirRecord, found := getDirRes.Records[0].Get("d")
	if !found {
		return models.Directory{}, apierrors.DirectoryNotFound{}
	}
	directory := models.GetDirectoryFromRecord(dirRecord)
	return directory, nil
}

func (gds GraphDatabaseService) CountDirectoryContents(location string) (int64, error) {
	countDirCypher := `
		MATCH (d:Directory) WHERE d.location = $location
		OPTIONAL MATCH (d)-[:CONTAINS]->(c)
		RETURN count(c) as count
	`
	countDirCypherParams := map[string]any{
		"location": location,
	}
	countDirRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		countDirCypher, countDirCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return 0, err
	}
	count, found := countDirRes.Records[0].Get("count")
	if !found {
		return 0, apierrors.DirectoryNotFound{}
	}
	return count.(int64), nil
}

func (gds GraphDatabaseService) GetFileDetails(fileLocation string) (models.File, error) {
	getFileCypher := `
		MATCH (f:File) WHERE f.location = $fileLocation
		RETURN f
	`
	getFileParams := map[string]any{
		"fileLocation": fileLocation,
	}
	getFileRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getFileCypher, getFileParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.File{}, err
	}
	fileRecord, found := getFileRes.Records[0].Get("f")
	if !found {
		return models.File{}, apierrors.FileNotFound{}
	}
	file := models.GetFileFromRecord(fileRecord)
	return file, nil
}

func (gds GraphDatabaseService) GetDirectoryAndContentDetails(dirLocation string) (models.DirectoryWithContents, error) {
	getDirAndContentCypher := `
		MATCH (parent:Directory) WHERE parent.location = $dirLocation
		OPTIONAL MATCH (parent)-[:CONTAINS]->(c)
		RETURN parent, collect(c) as content
	`
	getDirAndContentCypherParams := map[string]any{
		"dirLocation": dirLocation,
	}
	getDirAndContentRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getDirAndContentCypher, getDirAndContentCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return models.DirectoryWithContents{}, err
	}
	if len(getDirAndContentRes.Records) == 0 {
		return models.DirectoryWithContents{}, apierrors.DirectoryNotFound{}
	}
	parentRecord, found := getDirAndContentRes.Records[0].Get("parent")
	if !found {
		return models.DirectoryWithContents{}, apierrors.DirectoryNotFound{}
	}
	parent := models.GetDirectoryFromRecord(parentRecord)
	record, found := getDirAndContentRes.Records[0].Get("content")
	contentList := []interface{}{}
	if found {
		recordList := record.([]any)
		for _, r := range recordList {
			label := r.(neo4j.Node).Labels[0]
			if label == "Directory" {
				contentList = append(contentList, models.GetDirectoryFromRecord(r))
			} else if label == "File" {
				contentList = append(contentList, models.GetFileFromRecord(r))
			}
		}
	}

	return models.DirectoryWithContents{
		Id:        parent.Id,
		Type:      "Directory",
		Name:      parent.Name,
		CreatedOn: parent.CreatedOn,
		Location:  parent.Location,
		Contents:  contentList,
	}, nil
}

func (gds GraphDatabaseService) DeleteDirectory(location string) error {
	deleteDirCypher := `
		MATCH (d:Directory) WHERE d.location = $location
		DETACH DELETE d
	`
	deleteDirCypherParams := map[string]any{
		"location": location,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		deleteDirCypher, deleteDirCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) DeleteFile(fileLocation string) error {
	deleteFileCypher := `
		MATCH (f:File) WHERE f.location = $fileLocation
		DETACH DELETE f
	`
	deleteFileParams := map[string]any{
		"fileLocation": fileLocation,
	}
	_, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		deleteFileCypher, deleteFileParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (gds GraphDatabaseService) GetSharedDirsAndFiles(accId string, workspace string) ([]any, error) {
	getSharedCypher := `
		MATCH (sa:ServiceAccount{id:$accId})-[:SERVICES]->(:Workspace{name:$workspace})
		MATCH (sa)-[:HAS_ROLE]->(r:Role)
		MATCH (r)-[:MANAGES]->(c:Directory|File)
		RETURN COLLECT(c) as contents
	`
	getSharedCypherParams := map[string]any{
		"accId":     accId,
		"workspace": workspace,
	}
	getSharedRes, err := neo4j.ExecuteQuery(gds.ctx, gds.driver,
		getSharedCypher, getSharedCypherParams,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, err
	}
	sharedRecord, found := getSharedRes.Records[0].Get("contents")
	if !found {
		return []any{}, nil
	}
	sharedList := []any{}
	for _, r := range sharedRecord.([]any) {
		label := r.(neo4j.Node).Labels[0]
		if label == "Directory" {
			sharedList = append(sharedList, models.GetDirectoryFromRecord(r))
		} else if label == "File" {
			sharedList = append(sharedList, models.GetFileFromRecord(r))
		}
	}
	return sharedList, nil
}
