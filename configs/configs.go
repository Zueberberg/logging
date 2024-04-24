package configs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"

	"github.com/Zueberberg/logging"
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

type LoggingClassStorage map[string]map[string]reflect.Type

func (lcs LoggingClassStorage) RegisterFormatter(fmt_cls Formattable) {
	if _, ok := lcs["formatters"]; !ok {
		lcs["formatters"] = map[string]reflect.Type{}
	}

	r := reflect.ValueOf(fmt_cls).Type()
	fmt_cls_name := string(re_typename.FindSubmatch([]byte(r.String()))[1])
	lcs["formatters"][fmt_cls_name] = r
}

func (lcs LoggingClassStorage) GetFormatter(cls_name string) Formattable {
	if fmt_cls, ok := lcs["formatters"][cls_name]; ok {
		fmt.Println(fmt_cls)
		fmt.Println("GetFormatter:", reflect.New(fmt_cls))
		fmt.Printf("GetFormatter: %T\n", reflect.New(fmt_cls))
		fmt.Println("GetFormatter:", reflect.New(fmt_cls).Type())
		r := reflect.New(fmt_cls)
		fmt.Println(r)
		return nil
		// return reflect.New(fmt_cls)
		// switch T := fmt_cls.(type) {
		// case Formattable:
		// 	return T
		// default:
		// 	return nil
		// }
	} else {
		loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['formatters']: Unknown Formatter Class name: '%s'\n", cls_name))
		return nil
	}
}

func (lcs LoggingClassStorage) RegisterHandler(hlr_cls Formattable) {
	if _, ok := lcs["handlers"]; !ok {
		lcs["handlers"] = map[string]reflect.Type{}
	}

	r := reflect.ValueOf(hlr_cls).Type()
	hlr_cls_name := string(re_typename.FindSubmatch([]byte(r.String()))[1])
	lcs["handlers"][hlr_cls_name] = r
}

func (lcs LoggingClassStorage) GetHandler(cls_name string) Handlerable {
	if fmt_cls, ok := lcs["handlers"][cls_name]; ok {
		switch T := fmt_cls.(type) {
		case Handlerable:
			return T
		default:
			return nil
		}
	} else {
		loggingLogger.Fatal(fmt.Sprintf("LoggingParsedStorage['handlers']: Unknown Formatter Class name: '%s'\n", cls_name))
		return nil
	}
}

var LCS LoggingClassStorage

type LoggingParsedStorage map[string]map[string]any

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

type JSONConfigStruct struct {
	Loggers    map[string]JSONLoggerConfig    `json:"loggers,omitempty"`
	Formatters map[string]JSONFormatterConfig `json:"formatters,omitempty"`
	// Handlers   map[string]Handlerable `json:"handlers,omitempty"`
}

func JSONConfig(fileName string) {
	LCS = LoggingClassStorage{}
	LCS.RegisterFormatter(Formatter{})
	r := LCS.GetFormatter("Formatter")
	fmt.Println(r)
	fmt.Printf("%#v: %T\n", r, r)
	fmt.Println(reflect.TypeOf(r))
	os.Exit(0)

	LPS = LoggingParsedStorage{
		"handlers": map[string]any{
			"DefaultConsoleHandler": DefaultConsoleHandler,
		},
	}
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

	l := LPS.GetLogger("mai")
	l = l.WithHandler(DefaultConsoleHandler)
	l.Debug("It works!")
	l.Error("It works!")
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
	// for formatterName, formatterConf := range c.Formatters {
	// 	// var new_formatter Formattable
	//         // var fmt_cls :=
	// }
	for loggerName, loggerConf := range c.Loggers {
		_, ok := LPS["loggers"]
		if !ok {
			LPS["loggers"] = map[string]any{}
		}

		logLevel := parseLogLevel(loggerConf.LogLevel)
		if logLevel == 666 {
			return fmt.Errorf("Logger: '%s', Incorrect level: '%s'\n", loggerName, loggerConf.LogLevel)
		}

		// TODO : Add handlers to new logger
		new_logger := Logger{}.WithName(loggerName).WithLevel(logLevel)
		LPS["loggers"][loggerName] = new_logger
	}
	return nil
}
