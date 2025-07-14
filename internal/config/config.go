package config

import (
	"downloader/internal/logger"
	"io"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// структура с всеми данными для конфигурации
type Config struct {
	HostAddr   string `yaml:"host_address"`
	InFileLog  bool   `yaml:"save_log_to_file"`
	StorageDir string `yaml:"storage_dir"`
	LogLevel   string `yaml:"log_level"`
	LimitFiles int    `yaml:"lim_files"`
	LimitTasks int    `yaml:"lim_tasks"`
	FileType   []string
}

// функция для создания конфигурации и
// чтения значений из файла
func InitConfig() (Config, error) {
	//Conf := CreateConfig()
	Conf := new(Config)

	err := ReadConfig("./Config.cfg", Conf)
	if err != nil {
		return *Conf, err
	}
	return *Conf, nil
}

// чтение конфигурационного файла
func ReadConfig(cfgFilePath string, config *Config) error {
	file, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	buff, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buff, config)
	if err != nil {
		return err
	}
	return nil
}

// создание конфигурационного файла
func CreateConfig() *Config {
	Conf := new(Config)
	Conf.HostAddr = ":8080"
	Conf.InFileLog = true
	Conf.StorageDir = "./Storage"
	Conf.FileType = []string{".pdf", ".jpg"}
	Conf.LogLevel = "INFO"
	Conf.LimitFiles = 3
	Conf.LimitTasks = 3
	WriteConfig("./Config.cfg", Conf)
	return Conf
}

// запись данных в файл конфигурации
func WriteConfig(cfgFilePath string, config *Config) {

	file, err := os.OpenFile(cfgFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("Wrong config file!", zap.Error(err))
		return
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	err = enc.Encode(&config)
	if err != nil {
		logger.Log.Error("Bad unmarshall config file!", zap.Error(err))
	}
	logger.Log.Info("Config file created!", zap.Error(err))
}
