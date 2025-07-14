package handlers

import (
	"crypto/rand"
	"downloader/internal/logger"
	"downloader/internal/storage"
	"net/http"

	"go.uber.org/zap"
)

// PostTask хэндлер для создания задачи
func PostTask(Hd HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var err error

		// если количество задач меньше заданного (трех)
		if len(Hd.Tasks) < Hd.NumTasks {
			Hd.Mu.Lock()
			defer Hd.Mu.Unlock()
			// создаем ID задачи
			TaskID := rand.Text()[:8]
			// создаем карту для файлов в задаче с их статусами
			URLs := make(map[string]string, Hd.NumTasks)
			Hd.Tasks[TaskID] = URLs
			// создание папки для файлов задачи в хранилище
			err = storage.NewTask(Hd.StorageAddr, TaskID)
			if err != nil {
				logger.Log.Error("Error in creating storage for task!", zap.Error(err))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			// возвращаем статус о создании и ID задачи
			ansStr := "\nTask ID is: " + TaskID + "\n"
			res.WriteHeader(http.StatusCreated)
			res.Header().Set("content-type", "text/plain")
			res.Write([]byte(ansStr))
		} else {
			res.WriteHeader(http.StatusTooManyRequests)
		}
	}
}
