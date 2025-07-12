package server

import (
	"downloader/internal/config"
	"downloader/internal/handlers"
	"downloader/internal/logger"
	"downloader/internal/models"
	"downloader/internal/storage"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	config.Config
	handlers.HandlersData
}

func StartServer() error {
	var err error
	serv := new(Server)
	reqCh := make(chan models.ChanURLs, 3)
	serv.Config, err = config.InitConfig()
	if err != nil {
		return err
	}

	logger.Initialize(serv.InFileLog)

	logger.Log.Info("Logger initialyzed!", zap.Bool("Loggin in file is", serv.InFileLog))

	serv.HandlersData = handlers.InitHandlersData(serv.StorageDir, reqCh)

	err = storage.NewStor(serv.StorageDir)
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Info("Error in creating storage", zap.Error(err))
		return err
	}

	DownloadReq(reqCh, &serv.HandlersData)

	srv := http.Server{
		Handler: serv.ChiRouter(),
		Addr:    serv.HostAddr,
	}

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Info("Error in ListenAndServe", zap.Error(err))
		return err
	}

	return nil
}

// ChiRouter создает роутер chi для хэндлеров.
func (s *Server) ChiRouter() chi.Router {

	rt := chi.NewRouter()

	rt.Use(logger.MiddlewareLogger)

	rt.Post("/task", handlers.PostTask(s.HandlersData))
	rt.Post("/task/{id}", handlers.PostUrl(s.HandlersData))
	// rt.Get("/{id}", handlers.GetAddr(s.hd))
	// rt.Get("/ping", handlers.Ping(s.hd))
	// rt.Post("/api/shorten", handlers.PostAddrJSON(s.hd))
	// rt.Post("/api/shorten/batch", handlers.PostBatch(s.hd))
	// rt.Get("/api/user/urls", handlers.GetAddrUser(s.hd))
	// rt.Delete("/api/user/urls", handlers.DeleteAddr(s.hd))

	return rt
}

func DownloadReq(inp chan models.ChanURLs, Hd *handlers.HandlersData) error {
	go func() {
		for URL := range inp {
			go func(URL models.ChanURLs) {
				// здесь написать логику запросов на скачивание
				fmt.Printf("TaskID: %s\t URL:%s\n", URL.TaskID, URL.URL)
			}(URL)
		}
	}()
	return nil
}
