package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger(logsDirPath string, env string) (*logrus.Logger, error) {
	logsFileName := fmt.Sprintf("logs-%d.log", time.Now().Unix())
	logsFilePath := filepath.Join(logsDirPath, logsFileName)

	// Check if the logs directory exists. Create a new one if it does not.
	if _, err := os.Stat(logsDirPath); os.IsNotExist(err) {
		err = os.Mkdir(logsDirPath, os.ModeDir)
		if err != nil {
			fmt.Println("Unable to create logs directory. ERR: ", err)
		}
		fmt.Println("Logs directory was created because it did not exist.")
	}

	logFile, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Multiwritter that assures the logger outputs to stdout and the log file.
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)

	// Set debug level depending on the environment
	if env == "development" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger, nil
}
