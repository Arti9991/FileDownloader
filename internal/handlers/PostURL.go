package handlers

import (
	"downloader/internal/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// PostTask хэндлер для создания задачи.
func PostUrl(Hd HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!", zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var err error

		if req.Header.Get("content-type") != "application/json" {
			logger.Log.Info("Bad content-type header with this path!", zap.String("header", req.Header.Get("content-type")))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var IncomeURL []URL
		err = json.NewDecoder(req.Body).Decode(&IncomeURL)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(IncomeURL)
		// err = storage.NewTask(Hd.StorageAddr, TaskID)
		// if err != nil {
		// 	logger.Log.Error("Error in creating storage for task!", zap.Error(err))
		// 	res.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		res.WriteHeader(http.StatusCreated)
	}
}
