package app

import (
	"github.com/spf13/viper"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/csv"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository/postgres"
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

	db, err := postgres.NewPostgresDB(&postgres.Config{
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

	repo := postgres.NewStorage(db)

	for _, file := range files {
		wg.Add(1)
		go csv.Process(file, rabbitDSN, done, repo, logger)
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
