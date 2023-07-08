package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *Logger {
	infoLogFile, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Не удалось открыть файл информационных логов: ", err)
	}

	// Открытие файла для записи логов об ошибках
	errorLogFile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Не удалось открыть файл логов об ошибках: ", err)
	}

	// Создание логгеров
	infoLogger := log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{infoLogger: infoLogger, errorLogger: errorLogger}
}

func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *Logger) Infof(message string, s ...string) {
	l.infoLogger.Printf(message, s)
}

func (l *Logger) Error(err error) {
	l.errorLogger.Println(err)
}

func (l *Logger) Fatalf(message string, s ...string) {
	l.infoLogger.Fatalf(message, s)
}
