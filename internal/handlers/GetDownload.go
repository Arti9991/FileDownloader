package handlers

import (
	"downloader/internal/logger"
	"fmt"
	"net/http"
	"os"
	"path"

	"go.uber.org/zap"
)

// GetDownload хэндлер для архивирования и загрузки задачи.
func GetDownload(Hd HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			logger.Log.Info("Only GET requests are allowed with this path!",
				zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var err error

		// получаем индентификатор из URL запроса
		TaskID := path.Base(req.URL.String())
		fmt.Println(TaskID)
		// проверяем существует ли такая задача
		URLs, has := Hd.Tasks[TaskID]
		if !has {
			logger.Log.Info("There's no task witch this TaskID",
				zap.String("taskID", TaskID))
			res.WriteHeader(http.StatusNoContent)
			return
		}
		// Создаем ZIP архив с файлами задачи
		ZipPath, err := CreateZipArchive(&Hd, TaskID, URLs)
		if err != nil {
			logger.Log.Error("Error in creating zip", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// удаляем задачу из списка и дирректорию с задачей
		os.RemoveAll(Hd.StorageAddr + "/" + TaskID)
		delete(Hd.Tasks, TaskID)

		// выдаем ZIP архив в ответ
		res.Header().Set("Content-Disposition", "attachment; filename="+TaskID+".zip")
		res.Header().Set("Content-Type", "application/zip")
		http.ServeFile(res, req, ZipPath)
	}
}
