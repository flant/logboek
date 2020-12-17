package logger

import (
	"fmt"
	"io"

	"github.com/gookit/color"

	"github.com/werf/logboek/internal/stream"
	"github.com/werf/logboek/pkg/level"
	stylePkg "github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"
)

type Manager struct {
	level  level.Level
	logger *Logger
	style  color.Style
}

func NewManager(logger *Logger, lvl level.Level) *Manager {
	return &Manager{
		logger: logger,
		level:  lvl,
	}
}

func (m *Manager) SetStyle(style color.Style) {
	m.style = style
}

func (m *Manager) Style() color.Style {
	return m.style
}

func (m *Manager) IsAccepted() bool {
	return m.level <= m.logger.acceptedLevel
}

func (m *Manager) LogLn(a ...interface{}) {
	m.logLnCustom(m.style, a...)
}

func (m *Manager) LogF(format string, a ...interface{}) {
	m.logFCustom(m.style, format, a...)
}

func (m *Manager) LogLnDetails(a ...interface{}) {
	m.LogLnWithCustomStyle(stylePkg.Details(), a...)
}

func (m *Manager) LogFDetails(format string, a ...interface{}) {
	m.LogFWithCustomStyle(stylePkg.Details(), format, a...)
}

func (m *Manager) LogLnHighlight(a ...interface{}) {
	m.LogLnWithCustomStyle(stylePkg.Highlight(), a...)
}

func (m *Manager) LogFHighlight(format string, a ...interface{}) {
	m.LogFWithCustomStyle(stylePkg.Highlight(), format, a...)
}

func (m *Manager) LogLnWithCustomStyle(style color.Style, a ...interface{}) {
	m.logLnCustom(style, a...)
}

func (m *Manager) LogFWithCustomStyle(style color.Style, format string, a ...interface{}) {
	m.logFCustom(style, format, a...)
}

func (m *Manager) LogOptionalLn() {
	if !m.IsAccepted() {
		return
	}

	m.getStream().EnableOptionalLn()
}

func (m *Manager) LogBlock(format string, a ...interface{}) types.LogBlockInterface {
	logBlock := m.getStream().NewLogBlock(m, format, a...)
	logBlock.Options(func(options types.LogBlockOptionsInterface) {
		options.Style(m.style)
	})
	return logBlock
}

func (m *Manager) LogProcessInline(format string, a ...interface{}) types.LogProcessInlineInterface {
	logProcessInline := m.getStream().NewLogProcessInline(m, format, a...)
	logProcessInline.Options(func(options types.LogProcessInlineOptionsInterface) {
		options.Style(m.style)
	})
	return logProcessInline
}

func (m *Manager) LogProcess(format string, a ...interface{}) types.LogProcessInterface {
	logProcess := m.getStream().NewLogProcess(m, format, a...)
	logProcess.Options(func(options types.LogProcessOptionsInterface) {
		options.Style(m.style)
	})
	return logProcess
}

func (m *Manager) logLnCustom(style color.Style, a ...interface{}) {
	m.logFCustom(style, "%s", fmt.Sprintln(a...))
}

func (m *Manager) logFCustom(style color.Style, format string, a ...interface{}) {
	m.formatAndLogF(style, false, format, a...)
}

func (m *Manager) formatAndLogF(style color.Style, cacheIncompleteLine bool, format string, a ...interface{}) {
	if !m.IsAccepted() {
		return
	}

	m.getStream().FormatAndLogF(style, cacheIncompleteLine, format, a...)
}

func (m *Manager) getStream() *stream.Stream {
	return m.logger.GetLevelStream(m.level)
}

func (m *Manager) Stream() io.Writer {
	return proxyStream{Manager: m}
}

type proxyStream struct {
	*Manager
}

func (s proxyStream) Write(data []byte) (int, error) {
	s.Manager.formatAndLogF(s.Manager.style, true, "%s", string(data))
	return len(data), nil
}
