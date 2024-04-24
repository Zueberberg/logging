package logging

import (
	"github.com/charmbracelet/lipgloss"
)

// TODO:
// Дополнить цвета, для лог левелов - hex коды +
// полноценные палитры для белого / черного цвета терминала

var (
	BlackColor = "0"

	FatalColor = "1"
	ErrorColor = "1"
	TraceColor = "#c4eb17"

	WarningColor = "3"
	InfoColor    = "4"

	DebugColor   = "7"
	SuccessColor = "2"
	FailColor    = "1"
)

var (
	COLOR_LOG_MSG_PAD_X = 0
	COLOR_LOG_MSG_PAD_Y = 1
)

const (
	logMsgFatal   = "FATAL"
	logMsgError   = "ERROR"
	logMsgTrace   = "TRACE"
	logMsgWarning = "WARNING"
	logMsgInfo    = "INFO"
	logMsgDebug   = "DEBUG"
	logMsgSuccess = "SUCCESS"
	logMsgFail    = "FAIL"
)

var coloredMsgLogs = map[string]string{
	logMsgFatal:   ColoredFatalMessage,
	logMsgError:   ColoredErrorMessage,
	logMsgTrace:   ColoredTraceMessage,
	logMsgWarning: ColoredWarningMessage,
	logMsgInfo:    ColoredInfoMessage,
	logMsgDebug:   ColoredDebugMessage,
	logMsgSuccess: ColoredSuccessMessage,
	logMsgFail:    ColoredFailMessage,
}

var ColoredFatalMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(FatalColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgFatal)

var ColoredErrorMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(ErrorColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgError)

var ColoredTraceMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(TraceColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgTrace)

var ColoredWarningMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(WarningColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgWarning)

var ColoredInfoMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(InfoColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgInfo)

var ColoredDebugMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(DebugColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgDebug)

var ColoredSuccessMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(SuccessColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgSuccess)

var ColoredFailMessage = lipgloss.NewStyle().
	Background(lipgloss.Color(FailColor)).
	Foreground(lipgloss.Color(BlackColor)).
	Padding(COLOR_LOG_MSG_PAD_X, COLOR_LOG_MSG_PAD_Y).
	Align(lipgloss.Center).
	Render(logMsgFail)
