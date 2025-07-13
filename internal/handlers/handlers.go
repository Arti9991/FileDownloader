package handlers

import (
	"archive/zip"
	"downloader/internal/models"
	"io"
	"os"
)

type HandlersData struct {
	StorageAddr string
	Tasks       map[string]map[string]string
	reqChan     chan models.ChanURLs
}

func InitHandlersData(Stor string, reqChan chan models.ChanURLs) HandlersData {
	HD := new(HandlersData)
	HD.StorageAddr = Stor
	HD.Tasks = make(map[string]map[string]string, 3)
	HD.reqChan = reqChan
	return *HD
}

func CreateZipArchive(Hd *HandlersData, TaskID string, files map[string]string) (string, error) {

	zipPath := Hd.StorageAddr + "/" + TaskID + ".zip"

	// files: map[имя_в_архиве]путь_к_файлу
	outFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	for FileName, Status := range files {
		if Status != "DONE" {
			continue
		}
		filePath := Hd.StorageAddr + "/" + TaskID + "/" + FileName
		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}

		writer, err := zipWriter.Create(FileName)
		if err != nil {
			file.Close()
			return "", err
		}

		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return "", err
		}
	}

	return zipPath, nil
}
