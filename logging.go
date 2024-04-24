package logging

import (
	"fmt"
	"os"
)

var (
	loggingLogger Logger
	DefaultLogger Logger

	DefaultFormatter     Formatter
	DefaultJSONFormatter JSONFormatter

	DefaultConsoleHandler StreamHandler
	DefaultFileHandler    FileHandler
)

func json_testing() Logger {
	logger := DefaultLogger.
		WithHandler(
			DefaultConsoleHandler.
				WithFormatter(
					DefaultJSONFormatter.
						WithIndent(true).
						WithColor(true))).
		WithHandler(
			DefaultFileHandler.
				WithFileName("root.log").
				WithFormatter(
					DefaultJSONFormatter.
						// WithIndent(true).
						WithColor(true)))
	return logger
}

func console_simple_test() Logger {
	logger := DefaultLogger.
		WithName("root").
		WithHandler(DefaultConsoleHandler)
	return logger
}

func console_file_simple_test() Logger {
	logger := DefaultLogger.
		WithName("root").
		WithHandler(DefaultConsoleHandler).
		WithHandler(DefaultFileHandler.WithFileName("root.log").WithName("root_file_handler")).
		WithHandler(DefaultFileHandler.
			WithFileName("main.log").
			WithName("main_file_handler").
			WithFormatter(DefaultFormatter.WithColor(true)).
			WithLogLevel(ERROR))
	return logger
}

func logger_test(logger Logger) {
	logger.Error("Simple", "error", "message")
	logger.Trace("Simple", "trace", "message")
	logger.Warning("Simple", "warning", "message")
	logger.Info("Simple", "info", "message")
	logger.Debug("Simple", "debug", "message")
	logger.Success("Simple", "success", "message")
	logger.Fail("Simple", "fail", "message")
	fmt.Println()

	logger.Success("Another success message")
	logger.Error("Another error message")
	fmt.Println()
	logger.Info("Programm will be terminated right now!")
	logger.Fatal("Oh noo...")
}

func main() {
	// logger := console_simple_test()
	// logger := console_file_simple_test()
	// logger := json_testing()
	//
	// logger_test(logger)

	JSONConfig("config.json")
}

func init() {
	DefaultFormatter = Formatter{}.
		WithName("DefaultFormatter").
		WithTimeFormat(BasicTimeFormat).
		WithLogFormat(BasicLogFormat)

	DefaultJSONFormatter = JSONFormatter{}.
		WithName("DefaultJSONFormatter").
		WithTimeFormat(BasicTimeFormat).
		WithLogFormat(BasicLogFormat)

	DefaultConsoleHandler = StreamHandler{}.
		WithName("DefaultConsoleHandler").
		WithLogLevel(DEBUG).
		WithWriter(os.Stdout).
		WithFormatter(DefaultFormatter.WithColor(true))

	DefaultFileHandler = FileHandler{}.
		WithName("DefaultFileHandler").
		WithLogLevel(DEBUG).
		WithFormatter(DefaultFormatter)

	DefaultLogger = Logger{}.
		WithName("root").
		WithLevel(DEBUG)

	loggingLogger = Logger{}.
		WithName("logging").
		WithLevel(DEBUG).
		WithHandler(DefaultConsoleHandler)
}
