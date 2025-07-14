package handlers

import (
	"archive/zip"
	"downloader/internal/models"
	"io"
	"os"
	"sync"
)

// структура с информацией для работы хэндлеров
type HandlersData struct {
	StorageAddr string
	NumTasks    int
	NumFiles    int
	Tasks       map[string]map[string]string
	Mu          *sync.Mutex
	reqChan     chan models.ChanURLs
	HdWG        *sync.WaitGroup
}

// инициализация стурктуры для работы хэндлеров
func InitHandlersData(Stor string,
	NumTasks, NumFiles int,
	reqChan chan models.ChanURLs,
	HdWG *sync.WaitGroup, MapMU *sync.Mutex) HandlersData {

	HD := new(HandlersData)
	HD.StorageAddr = Stor
	HD.NumTasks = NumTasks
	HD.NumFiles = NumFiles
	HD.Tasks = make(map[string]map[string]string, 3)
	HD.reqChan = reqChan
	HD.HdWG = HdWG
	HD.Mu = MapMU
	return *HD
}

// функция создания zip архива из файлов задачи
func CreateZipArchive(Hd *HandlersData, TaskID string, files map[string]string) (string, error) {
	// путь к zip архиву
	zipPath := Hd.StorageAddr + "/" + TaskID + ".zip"
	// создание выходного файла
	outFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	// создания writer-а в архив
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()
	// запись всех файлов из папки в архив
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
