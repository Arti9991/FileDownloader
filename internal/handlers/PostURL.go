package handlers

import (
	"downloader/internal/logger"
	"downloader/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"go.uber.org/zap"
)

// PostTask хэндлер для создания задачи.
func PostUrl(Hd HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!",
				zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var err error

		if req.Header.Get("content-type") != "application/json" {
			logger.Log.Info("Bad content-type header with this path!",
				zap.String("header", req.Header.Get("content-type")))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// получаем индентификатор из URL запроса
		TaskID := path.Base(req.URL.String())
		fmt.Println(TaskID)

		URLs, has := Hd.Tasks[TaskID]
		if !has || len(URLs) == 3 {
			logger.Log.Info("There's no task witch this TaskID or this task completed",
				zap.String("taskID", TaskID))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// if len(Hd.Tasks) >= 3 {
		// 	logger.Log.Info("List this tasks is full", zap.String("taskID", TaskID))
		// 	res.WriteHeader(http.StatusTooManyRequests)
		// 	return
		// }

		var IncomeURL []models.URL
		err = json.NewDecoder(req.Body).Decode(&IncomeURL)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(IncomeURL)-len(URLs) < 0 || len(IncomeURL) == 0 {
			logger.Log.Info("Too many URLs in request")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, inc := range IncomeURL {
			var Send models.ChanURLs
			Send.TaskID = TaskID
			Send.URL = inc.RespURL
			Hd.Tasks[TaskID][inc.RespURL] = "PROCEED"
			Hd.reqChan <- Send
		}

		//fmt.Println(IncomeURL)
		res.WriteHeader(http.StatusAccepted)
	}
}
