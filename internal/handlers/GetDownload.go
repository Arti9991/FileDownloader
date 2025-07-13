package handlers

import (
	"downloader/internal/logger"
	"fmt"
	"net/http"
	"os"
	"path"

	"go.uber.org/zap"
)

// PostTask хэндлер для создания задачи.
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

		URLs, has := Hd.Tasks[TaskID]
		if !has {
			logger.Log.Info("There's no task witch this TaskID",
				zap.String("taskID", TaskID))
			res.WriteHeader(http.StatusNoContent)
			return
		}

		ZipPath, err := CreateZipArchive(&Hd, TaskID, URLs)
		if err != nil {
			logger.Log.Error("Error in creating zip", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		os.RemoveAll(Hd.StorageAddr + "/" + TaskID)

		// добавить логику уления файла архива после отправки (возомжно при шатдауне)
		// сами по себе папки с файлами будут удаляться вместе с отправкой архива
		// также добавить возможность удаления задания
		res.Header().Set("Content-Disposition", "attachment; filename="+TaskID+".zip")
		res.Header().Set("Content-Type", "application/zip")
		http.ServeFile(res, req, ZipPath)
	}
}
