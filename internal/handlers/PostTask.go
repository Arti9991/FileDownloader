package handlers

import (
	"crypto/rand"
	"downloader/internal/logger"
	"downloader/internal/storage"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// PostTask хэндлер для создания задачи.
func PostTask(Hd HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var err error

		if len(Hd.Tasks) < 3 {
			TaskID := rand.Text()[:8]
			fmt.Println(TaskID)

			Hd.Tasks = append(Hd.Tasks, TaskID)
			err = storage.NewTask(Hd.StorageAddr, TaskID)
			if err != nil {
				logger.Log.Error("Error in creating storage for task!", zap.Error(err))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			ansStr := "\nTask ID is: " + TaskID + "\n"
			res.WriteHeader(http.StatusCreated)
			res.Header().Set("content-type", "text/plain")
			res.Write([]byte(ansStr))
		} else {
			res.WriteHeader(http.StatusTooManyRequests)
		}
	}
}
