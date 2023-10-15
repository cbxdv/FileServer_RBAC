package fileservice

import (
	"bufio"
	"errors"
	"fs_backend/models"
	"io"
	"log"
	"os"
	"strings"
)

type FileService struct {
	internalLocation  string
	DownloadChunkSize int
}

func (fs *FileService) Initialize() {
	location := os.Getenv("STORAGE_LOCATION")
	if location == "" {
		location = "_storage/"
	}
	fs.internalLocation = location

	// Creating temp directory
	if _, err := os.Stat(fs.internalLocation); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fs.internalLocation, os.ModePerm)
		if err != nil {
			log.Default().Println(err.Error())
		}
	}

	// Chunk size for download
	fs.DownloadChunkSize = 1024 * 1024
}

func (fs FileService) GetInternalLocation() string {
	return fs.internalLocation
}

func (fs FileService) Cleanup() {
	os.RemoveAll(fs.internalLocation)
}

func (fs FileService) CreateWorkspaceSpace(workspaceName string) error {
	err := os.Mkdir(fs.internalLocation+"/"+workspaceName, os.ModePerm)
	return err
}

func (fs FileService) DeleteWorkspaceSpace(workspaceName string) error {
	err := os.RemoveAll(fs.internalLocation + "/" + workspaceName)
	return err
}

func (fs FileService) getFileLocation(fileProperties models.File) string {
	locationSplit := strings.Split(fileProperties.Location, "/")
	workspaceName := locationSplit[0]
	fileInternalLocation := fs.internalLocation + "/" + workspaceName + "/" + fileProperties.Id
	return fileInternalLocation
}

func (fs FileService) ReadChunkFromFile(fileProperties models.File, chunkNumber int) ([]byte, error) {
	file, err := os.Open(fs.getFileLocation(fileProperties))
	if err != nil {
		return []byte{}, err
	}

	start := fs.DownloadChunkSize * (chunkNumber - 1)

	if chunkNumber != 1 {
		file.Seek(int64(start), io.SeekStart)
	}

	var chunk []byte
	if start+fs.DownloadChunkSize > fileProperties.Size {
		chunk = make([]byte, fileProperties.Size-start)
	} else {
		chunk = make([]byte, fs.DownloadChunkSize)
	}

	_, err = bufio.NewReader(file).Read(chunk)
	if err != nil {
		return []byte{}, err
	}

	return chunk, nil
}

func (fs FileService) WriteChunkToFile(fileProperties models.File, chunk []byte) error {
	file, err := os.OpenFile(fs.getFileLocation(fileProperties), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	defer file.Close()

	_, err = file.Write(chunk)
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}

func (fs FileService) DeleteFileFromInternalLocation(fileProperties models.File) error {
	if _, err := os.Stat(fs.internalLocation); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fs.internalLocation, os.ModePerm)
		if err != nil {
			log.Default().Println(err.Error())
			return err
		}
	}
	err := os.Remove(fs.getFileLocation(fileProperties))
	if err != nil {
		log.Default().Println(err.Error())
		return err
	}
	return nil
}
