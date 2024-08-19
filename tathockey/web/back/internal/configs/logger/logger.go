package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/natefinch/lumberjack"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"time"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.GreenString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[02 Jan - 15:04:05]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}

func setupPrettySlog() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

// CustomJSONHandler кастомный обработчик для цветного JSON логирования в dev режиме
type CustomJSONHandler struct {
	writer *os.File
	level  slog.Level
}

func (h *CustomJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup возвращает новый обработчик с добавленной группой
func (h *CustomJSONHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled проверяет, нужно ли логгировать данный уровень
func (h *CustomJSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle логгирует сообщение в формате JSON с цветом для уровня
func (h *CustomJSONHandler) Handle(_ context.Context, r slog.Record) error {
	data := map[string]interface{}{
		"time":    r.Time.Format(time.RFC3339),
		"level":   r.Level.String(), // сохраняем уровень без цвета
		"message": r.Message,
	}

	r.Attrs(func(a slog.Attr) bool {
		data[a.Key] = a.Value.Any()
		return true
	})

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Выводим JSON с цветным уровнем логирования
	_, _ = fmt.Fprintf(h.writer, "%s", jsonData)

	return nil
}

// InitLogger инициализирует логгер для разных режимов (dev/prod)
func InitLogger() *slog.Logger {
	var handler slog.Handler
	var logger *slog.Logger
	env := os.Getenv("APP_ENV")
	switch env {
	case "local":
		logger = setupPrettySlog()
	case "prod":
		// Логирование в файл с ротацией в формате JSON
		logFile := &lumberjack.Logger{
			Filename:   "app.log",
			MaxSize:    10,   // MB
			MaxBackups: 3,    // Сколько старых файлов хранить
			MaxAge:     28,   // Сколько дней хранить файлы
			Compress:   true, // Сжимать старые файлы
		}
		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo, // Уровень логирования для production
		})
		logger = slog.New(handler)
	case "dev":
		// Используем кастомный обработчик для цветного JSON логирования
		handler = &CustomJSONHandler{writer: os.Stdout, level: slog.LevelDebug}
		logger = slog.New(handler)
	default:
		stdLog.Fatalf("Unknown environment: %s", env)
	}

	return logger
}
