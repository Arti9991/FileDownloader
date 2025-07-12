package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

var config = "INFO"

// инициализация zap логгера (уровень логгирования INFO)
func Initialize(FileLog bool) {
	var infile zapcore.WriteSyncer
	var core zapcore.Core
	var file *os.File
	var err error

	var level zap.AtomicLevel

	out := zapcore.AddSync(os.Stdout)
	if FileLog {
		file, err = os.Create("logger.log")
		if err != nil {
			fmt.Println(err)
			FileLog = false
		} else {
			infile = zapcore.AddSync(file)
		}
	}

	// создание новой конфигурацию логера
	// для консоли
	ConsoleCfg := zap.NewDevelopmentEncoderConfig()
	// для файлов
	FileCfg := zap.NewProductionEncoderConfig()
	// установка времени
	ConsoleCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	FileCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
	// устанавка уровня
	if config == "INFO" {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else if config == "ERROR" {
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	// цветовая индикацию для консоли
	ConsoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(ConsoleCfg)
	fileEncoder := zapcore.NewJSONEncoder(FileCfg)

	if FileLog {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, out, level),
			zapcore.NewCore(fileEncoder, infile, level),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, out, level),
		)
	}
	// установка синглтона
	Log = zap.New(core)
}

// переопределение методов write и WriteHeader для использования middleware
type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Переопределение функции для интерфейса.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// Переопределение функции для интерфейса.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// MiddlewareLogger обработчик для zap логгера с логированием полученных и отправленных запросов
func MiddlewareLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		reslog := loggingResponseWriter{
			ResponseWriter: res, //встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&reslog, req)
		duration := time.Since(start)
		Log.Info("got incoming HTTP request",
			zap.String("URI", req.RequestURI),
			zap.String("method", req.Method),
		)
		Log.Info("responce on request",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
			zap.Duration("duration", duration),
		)
	})
}

// // инициализация zap логгера (уровень логгирования INFO)
// func InitializeTest(FileLog bool) {
// 	var infile zapcore.WriteSyncer
// 	var core zapcore.Core
// 	var file *os.File
// 	var err error

// 	out := zapcore.AddSync(os.Stdout)
// 	if FileLog {
// 		file, err = os.Create("logger.log")
// 		if err != nil {
// 			fmt.Println(err)
// 			FileLog = false
// 		} else {
// 			infile = zapcore.AddSync(file)
// 		}
// 	}

// 	// создаём новую конфигурацию логера
// 	ConsoleCfg := zap.NewDevelopmentEncoderConfig()
// 	FileCfg := zap.NewProductionEncoderConfig()
// 	// устанавливаем время
// 	ConsoleCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
// 	FileCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC1123)
// 	// устанавливаем уровень
// 	level := zap.NewAtomicLevelAt(zap.FatalLevel)
// 	// включаем цветовую индикацию для консоли
// 	ConsoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

// 	consoleEncoder := zapcore.NewConsoleEncoder(ConsoleCfg)
// 	fileEncoder := zapcore.NewJSONEncoder(FileCfg)

// 	if FileLog {
// 		core = zapcore.NewTee(
// 			zapcore.NewCore(consoleEncoder, out, level),
// 			zapcore.NewCore(fileEncoder, infile, level),
// 		)
// 	} else {
// 		core = zapcore.NewTee(
// 			zapcore.NewCore(consoleEncoder, out, level),
// 		)
// 	}
// 	// установка синглтона
// 	Log = zap.New(core)
// }
