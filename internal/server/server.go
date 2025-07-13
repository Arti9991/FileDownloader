package server

import (
	"context"
	"downloader/internal/config"
	"downloader/internal/handlers"
	"downloader/internal/logger"
	"downloader/internal/models"
	"downloader/internal/storage"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	config.Config
	handlers.HandlersData
	ServWG *sync.WaitGroup
}

func StartServer() error {
	var err error
	serv := new(Server)
	reqCh := make(chan models.ChanURLs, 3)
	shutCh := make(chan struct{})

	var HandWG sync.WaitGroup
	var ServWG sync.WaitGroup

	var MapMU sync.Mutex

	serv.ServWG = &ServWG

	serv.Config, err = config.InitConfig()
	if err != nil {
		return err
	}
	// контекст для ожидания системного сигнала на завершение работы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Initialize(serv.InFileLog)

	logger.Log.Info("Logger initialyzed!", zap.Bool("Loggin in file is", serv.InFileLog))

	serv.HandlersData = handlers.InitHandlersData(serv.StorageDir, reqCh, &HandWG, &MapMU)

	err = storage.NewStor(serv.StorageDir)
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Info("Error in creating storage", zap.Error(err))
		return err
	}

	DownloadReq(reqCh, serv)

	srv := http.Server{
		Handler: serv.ChiRouter(),
		Addr:    serv.HostAddr,
	}

	RunWaitShutdown(ctx, shutCh, &srv)

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Info("Error in ListenAndServe", zap.Error(err))
		return err
	}
	// ожидание сообщения о Shutdown
	<-shutCh

	serv.HdWG.Wait()
	serv.ServWG.Wait()

	close(reqCh)

	os.RemoveAll(serv.StorageAddr)

	logger.Log.Info("Server shutted down!")

	return nil
}

// ChiRouter создает роутер chi для хэндлеров.
func (s *Server) ChiRouter() chi.Router {

	rt := chi.NewRouter()

	rt.Use(logger.MiddlewareLogger)

	rt.Post("/task", handlers.PostTask(s.HandlersData))
	rt.Post("/task/{id}", handlers.PostUrl(s.HandlersData))
	rt.Get("/info/{id}", handlers.GetStatus(s.HandlersData))
	rt.Get("/download/{id}", handlers.GetDownload(s.HandlersData))
	// rt.Get("/{id}", handlers.GetAddr(s.hd))
	// rt.Get("/ping", handlers.Ping(s.hd))
	// rt.Post("/api/shorten", handlers.PostAddrJSON(s.hd))
	// rt.Post("/api/shorten/batch", handlers.PostBatch(s.hd))
	// rt.Get("/api/user/urls", handlers.GetAddrUser(s.hd))
	// rt.Delete("/api/user/urls", handlers.DeleteAddr(s.hd))

	return rt
}

func DownloadReq(inp chan models.ChanURLs, Srv *Server) error {
	go func() {
		for val := range inp {
			Srv.ServWG.Add(1)
			go func(val models.ChanURLs) {

				defer Srv.ServWG.Done()
				regName := regexp.MustCompile(`[^/]+$`)
				name := regName.FindString(val.URL)

				//regExt := regexp.MustCompile(`[^.]+$`)
				// extension := regExt.FindString(name)
				// fmt.Println(extension)

				fmt.Println("Есть ли суффикс:", strings.HasSuffix(name, ".png"))
				AceptedExt := false
				for _, ext := range Srv.FileType {
					if strings.HasSuffix(name, ext) {
						AceptedExt = true
					}
				}

				Srv.Mu.Lock()
				defer Srv.Mu.Unlock()
				if !AceptedExt {
					logger.Log.Info("Bad extention!", zap.String("URL:", val.URL))
					Srv.Tasks[val.TaskID][name] = "ERROR"
					return
				}

				Srv.Tasks[val.TaskID][name] = "PROCESS"
				body, err := Request(val.URL)
				if err != nil {
					Srv.Tasks[val.TaskID][name] = "ERROR"
					return
				}

				filePath := Srv.StorageAddr + "/" + val.TaskID + "/" + name
				err = os.WriteFile(filePath, body, 0644)
				if err != nil {
					logger.Log.Info("Error in write body for URL!",
						zap.Error(err), zap.String("URL:", val.URL))
					Srv.Tasks[val.TaskID][name] = "ERROR"
					return
				}
				Srv.Tasks[val.TaskID][name] = "DONE"
			}(val)
		}
	}()
	return nil
}

func Request(URL string) ([]byte, error) {
	clientReq := &http.Client{
		Timeout: 10 * time.Second,
	}
	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		logger.Log.Info("Error in creating request for URL!",
			zap.Error(err), zap.String("URL:", URL))
		return nil, err
	}
	// отправляем запрос и получаем ответ
	response, err := clientReq.Do(request)
	if err != nil {
		logger.Log.Info("Error in applying request for URL!",
			zap.Error(err), zap.String("URL:", URL))
		return nil, err
	}
	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Log.Info("Error in read body for URL!",
			zap.Error(err), zap.String("URL:", URL))
		return nil, err
	}
	return body, nil
}

func RunWaitShutdown(ctx context.Context, shutCh chan struct{}, server *http.Server) {
	go func() {
		<-ctx.Done()
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		logger.Log.Info("Graceful shutdown...")
		if err := server.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			logger.Log.Info("Error in HTTP server Shutdown", zap.Error(err))
			return
		}

		// сообщение о Shutdown
		close(shutCh)
	}()
}
