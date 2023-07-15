package app

import (
	"github.com/spf13/viper"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/csv"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/broker/rabbit"
	postgres2 "github.com/tumbleweedd/intership/CSV_Consumer/pkg/database/postgres"
	myLogger "github.com/tumbleweedd/intership/CSV_Consumer/pkg/logger"
	"os"
	"path/filepath"
	"sync"
)

const (
	rabbitDSN = "amqp://guest:guest@localhost:5672/"
	rootDir   = "../data"
)

func Run() {

	logger := myLogger.NewLogger()

	if err := initConfig(); err != nil {
		logger.Fatalf("Ошибка инициализации конфига: %", err.Error())
	}

	files, err := filePathWalkDir(rootDir)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	defer close(done)

	wg := &sync.WaitGroup{}

	db, err := postgres2.NewPostgresDB(&postgres2.Config{
		PgPort:    viper.GetString("db.port"),
		PgHost:    viper.GetString("db.host"),
		PgDBName:  viper.GetString("db.dbname"),
		PgUser:    viper.GetString("db.username"),
		PgPwd:     viper.GetString("db.password"),
		PgSSLMode: viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewStorage(db)
	initRabbitCon, err := rabbit.NewRabbitMQConnection(rabbitDSN)
	defer initRabbitCon.Close()
	if err != nil {
		logger.Infof("Ошибка при инициализации RabbitConnection: %s", err.Error())
		return
	}

	for _, file := range files {
		wg.Add(1)
		go csv.Process(file, initRabbitCon, done, repo, logger)
	}

	wg.Wait()
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func initConfig() error {
	viper.AddConfigPath("internal/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
