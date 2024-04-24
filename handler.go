package logging

import (
	"fmt"
	"io"
	"os"
)

type Handlerable interface {
	WriteLog(vals LogValues)
}

type handlerABC struct {
	Name      string
	LogLevel  LogLevel
	Formatter Formattable
}

type StreamHandler struct {
	Writer io.Writer
	handlerABC
}

func (h StreamHandler) WriteLog(vals LogValues) {
	if h.LogLevel <= vals.Level {
		log := h.Formatter.ParseLog(vals)
		_, err := fmt.Fprint(h.Writer, log)
		if err != nil {
			loggingLogger.Error(err)
		}
	}
}

func (h StreamHandler) WithName(n string) StreamHandler {
	h.Name = n
	return h
}

func (h StreamHandler) WithWriter(w io.Writer) StreamHandler {
	h.Writer = w
	return h
}

func (h StreamHandler) WithLogLevel(l LogLevel) StreamHandler {
	h.LogLevel = l
	return h
}

func (h StreamHandler) WithFormatter(f Formattable) StreamHandler {
	h.Formatter = f
	return h
}

type FileHandler struct {
	FileName string
	handlerABC
}

func (h FileHandler) WriteLog(vals LogValues) {
	if h.FileName == "" {
		DefaultLogger.Fatal(fmt.Sprintf("Empty FileHandler (%s) FileName!", h.Name))
		return
	}
	if h.LogLevel <= vals.Level {
		file, err := os.OpenFile(h.FileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			loggingLogger.Error(fmt.Sprintf("FileHandler (%s): %s", h.Name, err))
			return
		}
		defer file.Close()

		log := h.Formatter.ParseLog(vals)
		_, err = fmt.Fprint(file, log)
		if err != nil {
			loggingLogger.Error(fmt.Sprintf("FileHandler (%s): %s", h.Name, err))
		}
	}
}

func (h FileHandler) WithName(n string) FileHandler {
	h.Name = n
	return h
}

func (h FileHandler) WithFileName(f string) FileHandler {
	h.FileName = f
	return h
}

func (h FileHandler) WithLogLevel(l LogLevel) FileHandler {
	h.LogLevel = l
	return h
}

func (h FileHandler) WithFormatter(f Formattable) FileHandler {
	h.Formatter = f
	return h
}
