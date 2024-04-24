package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

var re_typename = regexp.MustCompile(`\.(\w+)$`)

type BasicConfigOpts struct {
	LoggerName string
	LogLevel   LogLevel
}

func BasicConfig(o BasicConfigOpts) (logger Logger) {
	if o.LoggerName == "" {
		logger.Name = "root"
	} else {
		logger.Name = o.LoggerName
	}

	return logger.WithLevel(o.LogLevel)
}

// JSON file configuration

type LoggingParsedStorage map[string]map[string]any

func (lps LoggingParsedStorage) AddLogger(l Logger) {
	if _, ok := lps["loggers"]; !ok {
		lps["loggers"] = map[string]any{}
	}
	lps["loggers"][l.Name] = l
}

func (lps LoggingParsedStorage) GetLogger(name string) Logger {
	if logger, ok := lps["loggers"][name]; ok {
		switch T := logger.(type) {
		case Logger:
			return T
		default:
			loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['loggers']['%s']: Get a unknown logger object %#v", name, logger))
			return Logger{}
		}
	} else {
		loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['loggers']: Unknown logger name: '%s'\n", name))
		return Logger{}
	}
}

func (lps LoggingParsedStorage) AddFormatter(name string, f Formattable) {
	if _, ok := lps["formatters"]; !ok {
		lps["formatters"] = map[string]any{}
	}
	lps["formatters"][name] = f
}

func (lps LoggingParsedStorage) GetFormatter(name string) Formattable {
	if frmtr, ok := lps["formatters"][name]; ok {
		switch T := frmtr.(type) {
		case Formattable:
			return T
		default:
			loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['formatters']['%s']: Get a unknown formatter object %#v", name, frmtr))
			return nil
		}
	} else {
		loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['formatters']: Unknown formatter name: '%s'\n", name))
		return nil
	}
}

func (lps LoggingParsedStorage) AddHandler(name string, h Handlerable) {
	if _, ok := lps["handlers"]; !ok {
		lps["handlers"] = map[string]any{}
	}
	lps["handlers"][name] = h
}

func (lps LoggingParsedStorage) GetHandler(name string) Handlerable {
	if frmtr, ok := lps["handlers"][name]; ok {
		switch T := frmtr.(type) {
		case Handlerable:
			return T
		default:
			loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['handlers']['%s']: Get a unknown handler object %#v", name, frmtr))
			return nil
		}
	} else {
		loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['handlers']: Unknown handler name: '%s'\n", name))
		return nil
	}
}

var LPS LoggingParsedStorage

type JSONLoggerConfig struct {
	LogLevel string   `json:"level,omitempty"`
	Handlers []string `json:"handlers,omitempty"`
}

type JSONFormatterConfig struct {
	Class       string `json:"class,omitempty"`
	TimeFormat  string `json:"timefmt,omitempty"`
	LogFormat   string `json:"fmt,omitempty"`
	ColoredLog  bool   `json:"colored,omitempty"`
	UsingIndent bool   `json:"indent,omitempty"`
}

type JSONHandlerConfig struct {
	Class     string `json:"class,omitempty"`
	LogLevel  string `json:"level,omitempty"`
	Formatter string `json:"formatter,omitempty"`
	Stream    string `json:"stream,omitempty"`
	Filename  string `json:"filename,omitempty"`
}

type JSONConfigStruct struct {
	Loggers    map[string]JSONLoggerConfig    `json:"loggers,omitempty"`
	Formatters map[string]JSONFormatterConfig `json:"formatters,omitempty"`
	Handlers   map[string]JSONHandlerConfig   `json:"handlers,omitempty"`
}

func JSONConfig(fileName string) {
	// TODO: Нужно это или нет?
	// LPS = LoggingParsedStorage{
	// 	"handlers": map[string]any{
	// 		"DefaultConsoleHandler": DefaultConsoleHandler,
	// 	},
	// }
	//
	LPS = LoggingParsedStorage{}

	var configValues JSONConfigStruct

	file, err := os.Open(fileName)
	if err != nil {
		loggingLogger.Fatal("JSONConfig:", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		loggingLogger.Fatal("JSONConfig:", err)
	}

	err = json.Unmarshal(data, &configValues)
	if err != nil {
		loggingLogger.Fatal("JSONConfig:", err)
	}

	err = parseJsonConfig(configValues)
	if err != nil {
		loggingLogger.Fatal("JSONConfig:", err)
	}

	loggingLogger.Success("JSON Config file parsed:", fileName)
}

func parseLogLevel(levelName string) LogLevel {
	logLevel, ok := configLevelNames[levelName]
	if !ok {
		return 666
	} else {
		return logLevel
	}
}

func parseJsonConfig(c JSONConfigStruct) error {
	for formatterName, formatterConf := range c.Formatters {
		var new_formatter Formattable

		// TODO: С этим надо что то сделать, убрать
		// дублирование кода + динамическое распознование
		// классов (в том числе пользовательского)
		switch formatterConf.Class {
		case "Formatter":
			new_formatter = Formatter{
				Name:       formatterName,
				TimeFormat: formatterConf.TimeFormat,
				LogFormat:  formatterConf.LogFormat,
				ColoredLog: formatterConf.ColoredLog,
			}
		case "JSONFormatter":
			new_formatter = JSONFormatter{
				formatterConf.UsingIndent,
				Formatter{
					Name:       formatterName,
					TimeFormat: formatterConf.TimeFormat,
					LogFormat:  formatterConf.LogFormat,
					ColoredLog: formatterConf.ColoredLog,
				},
			}
		default:
			return fmt.Errorf("Unknown formatter class: %s\n", formatterConf.Class)
		}
		LPS.AddFormatter(formatterName, new_formatter)
	}
	for handlerName, handlerConf := range c.Handlers {
		var new_handler Handlerable

		logLevel := parseLogLevel(handlerConf.LogLevel)
		if logLevel == 666 {
			return fmt.Errorf("Handler: '%s', Incorrect level: '%s'\n", handlerName, handlerConf.LogLevel)
		}

		switch handlerConf.Class {
		case "StreamHandler":
			var stream io.Writer
			switch handlerConf.Stream {
			case "", "stdout":
				stream = os.Stdout
			case "stderr":
				stream = os.Stderr
			// HACK: А надо ли stdin?
			case "stdin":
				stream = os.Stdin
			default:
				return fmt.Errorf("Handler: '%s', Incorrect stream: '%s'\n", handlerName, handlerConf.Stream)
			}
			new_handler = StreamHandler{
				stream,
				handlerABC{
					Name:      handlerName,
					LogLevel:  logLevel,
					Formatter: LPS.GetFormatter(handlerConf.Formatter),
				},
			}
		case "FileHandler":
			new_handler = FileHandler{
				handlerConf.Filename,
				handlerABC{
					Name:      handlerName,
					LogLevel:  logLevel,
					Formatter: LPS.GetFormatter(handlerConf.Formatter),
				},
			}
		}
		LPS.AddHandler(handlerName, new_handler)
	}
	for loggerName, loggerConf := range c.Loggers {
		logLevel := parseLogLevel(loggerConf.LogLevel)
		if logLevel == 666 {
			return fmt.Errorf("Logger: '%s', Incorrect level: '%s'\n", loggerName, loggerConf.LogLevel)
		}

		// TODO: Add handlers to new logger
		new_logger := Logger{}.WithName(loggerName).WithLevel(logLevel)
		if len(loggerConf.Handlers) == 0 {
			loggingLogger.Warning("No handlers specified on logger:", loggerName)
		} else {
			for _, handlerName := range loggerConf.Handlers {
				new_logger = new_logger.WithHandler(LPS.GetHandler(handlerName))
			}
		}
		LPS.AddLogger(new_logger)
	}
	return nil
}
