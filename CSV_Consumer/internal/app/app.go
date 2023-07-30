package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/csv"
	"github.com/tumbleweedd/intership/CSV_Consumer/internal/repository"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/broker/rabbit"
	"github.com/tumbleweedd/intership/CSV_Consumer/pkg/database/mongo"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	rabbitDSN = "amqp://guest:guest@localhost:5672/"
	rootDir   = "../data"
)

func Run() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("ошибка инициализации конфига: %s", err.Error())
	}

	files, err := filePathWalkDir(rootDir)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	defer close(done)

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := mongo.NewMongoConnection(ctx, &mongo.Config{
		MongoHost:   viper.GetString("mongo_db.host"),
		MongoPort:   viper.GetString("mongo_db.port"),
		MongoDBName: viper.GetString("mongo_db.dbname"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewStorage(session.Client, viper.GetString("mongo_db.dbname"), viper.GetString("mongo_db.usrs_collection"))
	initRabbitCon, err := rabbit.NewRabbitMQConnection(rabbitDSN)
	defer initRabbitCon.Close()
	if err != nil {
		logrus.Infof("Ошибка при инициализации RabbitConnection: %s", err.Error())
		return
	}

	for _, file := range files {
		wg.Add(1)
		go csv.Process(file, initRabbitCon, done, repo)
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
