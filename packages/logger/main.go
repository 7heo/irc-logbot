package logger

import (
	"fmt"
	"os"
	"time"
)

func getFileName(channel string) string {
	return time.Now().Format(fmt.Sprintf("logs/%s-2006-01-02.txt", channel))
}

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type EventLogger struct {
	file *os.File
	timeFmt string
}

func CreateEventLogger(channel string, timeFmt string) (*EventLogger, error) {
	fileName := getFileName(channel)
	file, err := openFile(fileName)

	if err != nil {
		return nil, err
	}

	eventLogger := &EventLogger{
		file: file,
		timeFmt: timeFmt,
	}

	return eventLogger, nil
}

func (eventLogger *EventLogger) LogEvent(channel, message string) error {
	_, err := eventLogger.file.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format(eventLogger.timeFmt), message))
	if err != nil {
		return err
	}

	return nil
}
