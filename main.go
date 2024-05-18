package main

import (
	"fmt"
	"io"
	"mongo-fundamential/api"
	"mongo-fundamential/driver/mongodb"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Dir      string `env:"CONFIG_DIR" envDefault:"config.json"`
	Port     string `json:"port"`
	LogType  string
	LogLevel string
	LogFile  string
	MongoDB  string `json:"mongodb"`
}

var config Config

func init() {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc

	if err := env.Parse(&config); err != nil {
		log.Error("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}

	cfg := Config{
		Dir:      config.Dir,
		Port:     viper.GetString(`main.port`),
		LogType:  viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		LogFile:  viper.GetString(`main.log_file`),
		MongoDB:  viper.GetString(`main.mongodb`),
	}
	config = cfg

	if config.MongoDB == "enabled" {
		mongoConfig := mongodb.MongoConfg{
			Host:     viper.GetString(`mongodb.host`),
			Port:     viper.GetInt(`mongodb.port`),
			Database: viper.GetString(`mongodb.db`),
		}
		_, err := mongodb.NewMongoClient(mongoConfig)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func main() {
	fmt.Println("ok")
	server := api.NewServer()
	server.Start(config.Port)
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}

func createNewLogFile(logDir string) error {
	files, err := os.ReadDir("tmp")
	if err != nil {
		return err
	}
	last10dayUnix := time.Now().Add(-1 * 24 * time.Hour).Unix()
	for _, f := range files {
		tmp := strings.Split(f.Name(), ".")
		if len(tmp) > 2 {
			fileUnix, err := strconv.Atoi(tmp[2])
			if err != nil {
				return err
			} else if int64(fileUnix) < last10dayUnix {
				if err := os.Remove("tmp/" + f.Name()); err != nil {
					return err
				}
			}
		}
	}
	_, err = os.Stat(logDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err := os.Rename(logDir, fmt.Sprintf(logDir+".%d", time.Now().Unix())); err != nil {
		return err
	}
	return nil
}
