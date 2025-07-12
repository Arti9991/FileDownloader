package config

import (
	"downloader/internal/logger"
	"io"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HostAddr   string `env:"HOST_ADDRESS" yaml:"host_address"`
	InFileLog  bool   `yaml:"save_log_to_file"`
	StorageDir string `env:"STORAGE_DIR" yaml:"storage_dir"`
	FileType   []string
}

func InitConfig() (Config, error) {
	//Conf := CreateConfig()
	Conf := new(Config)

	err := ReadConfig("./Config.cfg", Conf)
	if err != nil {
		return *Conf, err
	}
	return *Conf, nil
}

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

func CreateConfig() *Config {
	Conf := new(Config)
	Conf.HostAddr = ":8080"
	Conf.InFileLog = true
	Conf.StorageDir = "./"
	Conf.FileType = []string{".pdf", ".jpg"}

	WriteConfig("./Config.cfg", Conf)
	return Conf
}

// создание файла конфигурации с данными переданными через флаги или переменными окружения
func WriteConfig(cfgFilePath string, config *Config) {
	//var config Config

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
