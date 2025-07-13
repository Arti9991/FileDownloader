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
func GetStatus(Hd HandlersData) http.HandlerFunc {
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

		URLs, has := Hd.Tasks[TaskID]
		if !has {
			logger.Log.Info("There's no task witch this TaskID",
				zap.String("taskID", TaskID))
			res.WriteHeader(http.StatusNoContent)
			return
		}
		var info models.ResponceDownload
		for key, val := range URLs {
			var buf models.URLInfo
			buf.Status = val
			buf.URL = key
			info.URLsInfo = append(info.URLsInfo, buf)
		}
		if len(URLs) == 3 {
			info.DownloadURL = "http://localhost:8080/download/" + TaskID
		}
		// кодирование тела ответа.
		out, err := json.MarshalIndent(info, "", " ")
		if err != nil {
			logger.Log.Info("Wrong responce body", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		//fmt.Println(IncomeURL)
		res.Header().Set("content-type", "application/json")
		res.Write(out)
		//res.WriteHeader(http.StatusOK)
	}
}
