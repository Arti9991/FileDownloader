package server

import (
	"downloader/internal/config"
	"downloader/internal/handlers"
	"downloader/internal/logger"
	"downloader/internal/storage"
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

	serv.Config, err = config.InitConfig()
	if err != nil {
		return err
	}

	logger.Initialize(serv.InFileLog)

	logger.Log.Info("Logger initialyzed!", zap.Bool("Loggin in file is", serv.InFileLog))

	serv.HandlersData = handlers.InitHandlersData(serv.StorageDir)

	err = storage.NewStor(serv.StorageDir)
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Info("Error in creating storage", zap.Error(err))
		return err
	}

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
	rt.Post("/task/add", handlers.PostUrl(s.HandlersData))
	// rt.Get("/{id}", handlers.GetAddr(s.hd))
	// rt.Get("/ping", handlers.Ping(s.hd))
	// rt.Post("/api/shorten", handlers.PostAddrJSON(s.hd))
	// rt.Post("/api/shorten/batch", handlers.PostBatch(s.hd))
	// rt.Get("/api/user/urls", handlers.GetAddrUser(s.hd))
	// rt.Delete("/api/user/urls", handlers.DeleteAddr(s.hd))

	return rt
}
