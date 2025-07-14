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

// PostUrl хэндлер для добавления URL в задачу
func PostUrl(Hd HandlersData) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Log.Info("Only POST requests are allowed with this path!",
				zap.String("method", req.Method))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var err error
		// проверка заголовка запроса
		if req.Header.Get("content-type") != "application/json" {
			logger.Log.Info("Bad content-type header with this path!",
				zap.String("header", req.Header.Get("content-type")))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// получение индентификатора задачи из URL запроса
		TaskID := path.Base(req.URL.String())
		fmt.Println(TaskID)
		// проверяем существует ли такая задача
		URLs, has := Hd.Tasks[TaskID]
		if !has {
			logger.Log.Info("There's no task witch this TaskID",
				zap.String("taskID", TaskID))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// декодируем тело запроса
		var IncomeURL []models.URL
		err = json.NewDecoder(req.Body).Decode(&IncomeURL)
		if err != nil {
			logger.Log.Info("Bad request body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// проверяем не пустое ли тело запроса
		if len(IncomeURL) == 0 {
			logger.Log.Info("Request is empty")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// проверяем количество URL в запросе
		// если URL больше, чем разрешено,
		// отбрасываем лишнее
		rem := Hd.NumFiles - len(URLs)
		if rem > 0 {
			if rem > len(IncomeURL) {
				rem = len(IncomeURL)
			}
			IncomeURL = IncomeURL[:rem]
		} else {
			logger.Log.Info("Task is full")
			res.WriteHeader(http.StatusTooManyRequests)
			return
		}
		// отправляем URL в канал запросов на скачивание
		Hd.HdWG.Add(1)
		ThreadSend(IncomeURL, Hd, TaskID)

		//fmt.Println(IncomeURL)
		res.WriteHeader(http.StatusAccepted)
	}
}

// функция отправки полученных ссылок в канал для запросов на скачивание
func ThreadSend(inp []models.URL, Hd HandlersData, TaskID string) {
	go func() {
		defer Hd.HdWG.Done()
		for _, val := range inp {
			var Send models.ChanURLs
			Send.TaskID = TaskID
			Send.URL = val.RespURL
			Hd.reqChan <- Send
		}
	}()
}
