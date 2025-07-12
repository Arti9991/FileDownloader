package storage

import (
	"downloader/internal/logger"
	"os"
	"strings"

	"go.uber.org/zap"
)

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

func NewTask(StorageDir string, TaskDir string) error {
	// срздаем директорию хранилища, указанную в конфигурации
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
