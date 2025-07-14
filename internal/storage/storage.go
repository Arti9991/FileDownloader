package storage

import (
	"downloader/internal/logger"
	"os"
	"strings"

	"go.uber.org/zap"
)

// функция для создания директории хранилища
func NewStor(StorageDir string) error {
	// срздаем директорию хранилища, указанную в конфигурации
	err := os.Mkdir(StorageDir, 0644)
	if err != nil {
		if strings.Contains(err.Error(), "exist") {
			return nil
		} else {
			logger.Log.Error("Error in creating directory", zap.Error(err))
			return err
		}
	}
	return nil
}

// функция для создания директории для задачи
func NewTask(StorageDir string, TaskDir string) error {
	// срздаем директорию задачи
	err := os.Mkdir(StorageDir+"/"+TaskDir, 0644)
	if err != nil {
		if strings.Contains(err.Error(), "exist") {
			return nil
		} else {
			logger.Log.Error("Error in creating directory", zap.Error(err))
			return err
		}
	}
	return nil
}
