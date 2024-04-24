package logging

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	BasicTimeFormat     = "02 January 2006 15:04:05"
	BasicLogFormat      = "%(time)s :: %(level)s :: %(name)s :: %(msg)s"
	re_fmt              = regexp.MustCompile(`\%\(\w+\)s`)
	re_struct_rec_level = regexp.MustCompile(`\%\(levelname\)s`)
)

type Formattable interface {
	ParseLog(vals LogValues) string
}

type Formatter struct {
	Name       string
	TimeFormat string
	LogFormat  string
	ColoredLog bool
}

func (f Formatter) ParseLog(vals LogValues) string {
	if f.ColoredLog {
		vals.LevelName = coloredMsgLogs[vals.LevelName]
	}
	return re_fmt.ReplaceAllStringFunc(f.LogFormat, func(s string) string {
		switch s {
		case "%(name)s":
			return vals.LoggerName
		case "%(time)s":
			return vals.Time.Format(f.TimeFormat)
		case "%(level)s":
			return vals.LevelName
		case "%(msg)s":
			return fmt.Sprintln(vals.Msg...)
		default:
			return s
		}
	})
}

func (f Formatter) WithName(n string) Formatter {
	f.Name = n
	return f
}

func (f Formatter) WithTimeFormat(tf string) Formatter {
	f.TimeFormat = tf
	return f
}

func (f Formatter) WithLogFormat(lf string) Formatter {
	f.LogFormat = lf
	return f
}

func (f Formatter) WithColor(b bool) Formatter {
	f.ColoredLog = b
	return f
}

type StructuredLogRecord struct {
	Time       string `json:"time,omitempty"`
	LoggerName string `json:"name,omitempty"`
	LevelName  string `json:"level,omitempty"`
	Message    string `json:"message,omitempty"`
}

func parseStructuredLog(lf string, tf string, vals LogValues) (rec StructuredLogRecord) {
	opts := re_fmt.FindAllString(lf, -1)
	for _, opt := range opts {
		switch opt {
		case "%(time)s":
			rec.Time = vals.Time.Format(tf)
		case "%(level)s":
			rec.LevelName = "%(levelname)s"
		case "%(name)s":
			rec.LoggerName = vals.LoggerName
		case "%(msg)s":
			temp := fmt.Sprintln(vals.Msg...)
			if len(temp) >= 1 {
				temp = temp[:len(temp)-1]
			}
			rec.Message = temp
		}
	}
	return
}

type StructFormatter struct {
	UsingIndent bool
	Formatter
}

type JSONFormatter StructFormatter

func (f JSONFormatter) ParseLog(vals LogValues) (res string) {
	var byte_json []byte
	var err error

	rec := parseStructuredLog(f.LogFormat, f.TimeFormat, vals)

	if f.UsingIndent {
		byte_json, err = json.MarshalIndent(rec, "", "\t")
	} else {
		byte_json, err = json.Marshal(rec)
	}

	if err != nil {
		msg := fmt.Sprintf("Logger [ %s ]; JSONFormatter [ %s ]; Error: %s", vals.LoggerName, f.Name, err)
		loggingLogger.Error(msg)
		res = ""
	} else {
		res = string(byte_json) + "\n"
	}
	levelname := vals.LevelName
	if f.ColoredLog {
		levelname = coloredMsgLogs[vals.LevelName]
	}
	res = re_struct_rec_level.ReplaceAllString(res, levelname)

	return
}

func (f JSONFormatter) WithName(n string) JSONFormatter {
	f.Name = n
	return f
}

func (f JSONFormatter) WithTimeFormat(tf string) JSONFormatter {
	f.TimeFormat = tf
	return f
}

func (f JSONFormatter) WithLogFormat(lf string) JSONFormatter {
	f.LogFormat = lf
	return f
}

func (f JSONFormatter) WithColor(b bool) JSONFormatter {
	f.ColoredLog = b
	return f
}

func (f JSONFormatter) WithIndent(b bool) JSONFormatter {
	f.UsingIndent = b
	return f
}
