package csv_parser

import (
	"encoding/csv"
	"github.com/tumbleweedd/intership/CSV_Parser/pkg/rabbitmq"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func ParseCSV(fileSource string, rabbitCon *rabbitmq.RabbitConnection, wg *sync.WaitGroup, done <-chan struct{}) {
	defer wg.Done()

	file, err := os.Open(fileSource)
	if err != nil {
		log.Printf("Ошибка os.Open: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	producer := rabbitmq.NewProducer(rabbitCon)
	if err != nil {
		log.Printf("Ошибка при инициализации RabbitConnection: %v", err)
		return
	}

	queueName := strings.Split(fileSource, "\\")
	log.Println(queueName)

	for {
		select {
		case <-done:
			return
		default:
			record, err := reader.Read()
			if err == io.EOF {
				log.Printf("Достигли конца файла: %v", err)
				return
			}
			if err != nil {
				log.Printf("Ошибка при чтении строки CSV файла: %v", err)
				continue
			}

			err = producer.SendToQueue(queueName[2], record)
			if err != nil {
				log.Printf("Ошибка %v при занесении в очередь файла: %s", err, fileSource)
			}
		}
	}

}
