package server

import (
	"context"
	"downloader/internal/config"
	"downloader/internal/handlers"
	"downloader/internal/logger"
	"downloader/internal/models"
	"downloader/internal/storage"
	"errors"
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

// Функция инициализации и запуска сервера
func StartServer() error {
	var err error
	serv := new(Server)
	// канал для отправки URL с запросами на скачивание файлов
	reqCh := make(chan models.ChanURLs, 3)
	// канал для ожидания сообщения о шатдауне
	shutCh := make(chan struct{})

	// wait группы для ожидания завершения всех запросов и отправок
	var HandWG sync.WaitGroup
	var ServWG sync.WaitGroup

	// мьютекс для карт с информацией о задачах
	var MapMU sync.Mutex

	serv.ServWG = &ServWG

	serv.Config, err = config.InitConfig()
	if err != nil {
		return err
	}

	// контекст для ожидания системного сигнала на завершение работы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Initialize(serv.InFileLog, serv.LogLevel)

	logger.Log.Info("Logger initialyzed!", zap.Bool("Loggin in file is", serv.InFileLog))

	serv.HandlersData = handlers.InitHandlersData(serv.StorageDir, serv.LimitTasks, serv.LimitFiles,
		reqCh, &HandWG, &MapMU)

	err = storage.NewStor(serv.StorageDir)
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Error("Error in creating storage", zap.Error(err))
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
	// ожидание конца работы всех отправителей
	// из хэндлеров и обработчиков
	serv.HdWG.Wait()
	serv.ServWG.Wait()
	// закрытие канала запросов на скачивание
	close(reqCh)
	// удаление папки хранилища со всеми файлами
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

// функция с горутиной для приема ссылок и скачивания файлов
func DownloadReq(inp chan models.ChanURLs, Srv *Server) error {
	go func() {
		// в основной горутине принимаются ссылки из хэндлера
		for val := range inp {
			// для каждой ссылки в отдельной горутине выполняется
			// проверка файла (расширения) и запрос по ссылке
			// на скачивание
			Srv.ServWG.Add(1)
			go func(val models.ChanURLs) {
				defer Srv.ServWG.Done()
				// при помощи регулярного выражения получаем имя файла
				regName := regexp.MustCompile(`[^/]+$`)
				name := regName.FindString(val.URL)

				// проверяем расширение файла
				AceptedExt := false
				for _, ext := range Srv.FileType {
					if strings.HasSuffix(name, ext) {
						AceptedExt = true
					}
				}

				Srv.Mu.Lock()
				defer Srv.Mu.Unlock()
				// если расширение не подходит ставим статус ERROR
				if !AceptedExt {
					logger.Log.Info("Bad extention!", zap.String("URL:", val.URL))
					Srv.Tasks[val.TaskID][name] = "ERROR"
					return
				}
				// если расширение подходит, ставим статус PROCESS
				Srv.Tasks[val.TaskID][name] = "PROCESS"

				// выполняем запрос
				body, err := Request(val.URL)
				if err != nil {
					Srv.Tasks[val.TaskID][name] = "ERROR"
					return
				}

				// сохраняем файл с полученным выше именем
				filePath := Srv.StorageAddr + "/" + val.TaskID + "/" + name
				err = os.WriteFile(filePath, body, 0644)
				if err != nil {
					logger.Log.Info("Error in write body for URL!",
						zap.Error(err), zap.String("URL:", val.URL))
					Srv.Tasks[val.TaskID][name] = "ERROR"
					return
				}
				// ставим статус "DONE"
				Srv.Tasks[val.TaskID][name] = "DONE"
			}(val)
		}
	}()
	return nil
}

// функция с запросом на скачивание файла по ссылке
func Request(URL string) ([]byte, error) {
	// устанавливаем таймаут 10 секунд
	clientReq := &http.Client{
		Timeout: 10 * time.Second,
	}
	// создание запроса с методом GET
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
	defer response.Body.Close()
	// чтение потока из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Log.Info("Error in read body for URL!",
			zap.Error(err), zap.String("URL:", URL))
		return nil, err
	}

	// проверка размера полученных данных
	if len(body) > 20971520 {
		logger.Log.Info("Recieved body is too big!")
		return nil, errors.New("body is too big")
	}
	return body, nil
}

// функция с горутиной ожидания сообщения о выключении
func RunWaitShutdown(ctx context.Context, shutCh chan struct{}, server *http.Server) {
	go func() {
		<-ctx.Done()
		// получен сигнал os.Interrupt, запуск процедуры graceful shutdown
		logger.Log.Info("Graceful shutdown...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Log.Info("Error in HTTP server Shutdown", zap.Error(err))
			return
		}

		// сообщение о Shutdown в основную горутину
		close(shutCh)
	}()
}
