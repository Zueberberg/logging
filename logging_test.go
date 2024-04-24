package logging

import (
	"fmt"
	"os"
	"testing"
)

func Test_basic_logger_works(t *testing.T) {
	fmt.Println("====================")
	logger := DefaultLogger.WithHandler(DefaultConsoleHandler)
	logger.Success("It works!")
	fmt.Println("====================")
}

func Test_json_config_logger_working(t *testing.T) {
	fmt.Println("====================")
	cfg := `
    {
      "formatters": {
        "main_colored": {
          "class": "Formatter",
          "timefmt": "02 01 2006",
          "fmt": "%(time)s <-> %(level)s <-> %(name)s <-> %(msg)s",
          "colored": true
        },
        "main_colorless": {
          "class": "JSONFormatter",
          "timefmt": "02 01 2006",
          "fmt": "%(time)s <-> %(level)s <-> %(name)s <-> %(msg)s",
          "colored": false,
          "indent": true
        }
      },
      "handlers": {
        "main_console": {
          "class": "StreamHandler",
          "level": "INFO",
          "stream": "stdout",
          "formatter": "main_colored"
        },
        "main_file": {
          "class": "FileHandler",
          "level": "ERROR",
          "filename": "main_file.log",
          "formatter": "main_colorless"
        }
      },
      "loggers": {
        "main": {
          "level": "DEBUG",
          "handlers": [
            "main_console",
            "main_file"
          ]
        }
      }
    }
    `
	file, _ := os.Create("config.json")
	fmt.Fprint(file, cfg)
	file.Close()

	JSONConfig(`config.json`)
	logger := LPS.GetLogger("main")
	if logger.LogLevel != configLevelNames["DEBUG"] {
		t.Error("Logger level wrong configured")
	}
	if len(logger.Handlers) != 2 {
		t.Error("Logger handlers wrong configured!")
	}

	logger.Debug("This is a message")
	logger.Info("This is a message")
	logger.Warning("This is a message")
	logger.Error("This is a message")
	os.Remove("config.json")
	fmt.Println("====================")
}
