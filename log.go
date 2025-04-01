package httplog

import (
	"cmp"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func JSON(outLogger, errLogger *log.Logger) func(req *http.Request, elapsed time.Duration, status int) {
	return func(req *http.Request, elapsed time.Duration, status int) {
		if status >= 500 {
			errLogger.Printf(`{"type": "HTTP_REQUEST", "method": %q, "path": %q, "duration": %q, "status": %d}`+"\n", req.Method, req.URL.Path, elapsed, status)
		}
		outLogger.Printf(`{"type": "HTTP_REQUEST", "method": %q, "path": %q, "duration": %q, "status": %d}`+"\n", req.Method, req.URL.Path, elapsed, status)
	}
}

func Structured(logger *slog.Logger) func(req *http.Request, elapsed time.Duration, status int) {
	level := ParseStructuredLogLevel("", slog.LevelInfo)
	return func(req *http.Request, elapsed time.Duration, status int) {
		if status >= 500 {
			logger.ErrorContext(req.Context(), "request error", slog.String("method", req.Method), slog.String("path", req.URL.Path), slog.Int("status", status), slog.Duration("duration", elapsed))
			return
		}
		logger.Log(req.Context(), level, "request", slog.String("method", req.Method), slog.String("path", req.URL.Path), slog.Int("status", status), slog.Duration("duration", elapsed))
	}
}

const DefaultStructuredLevelEnvironmentVariableName = "HTTP_LOG_LEVEL"

func ParseStructuredLogLevel(varName string, defaultLevel slog.Level) slog.Level {
	env := cmp.Or(varName, DefaultStructuredLevelEnvironmentVariableName)
	val, ok := os.LookupEnv(env)
	if !ok {
		return defaultLevel
	}
	switch val {
	case slog.LevelDebug.String():
		return slog.LevelDebug
	case slog.LevelInfo.String():
		return slog.LevelInfo
	case slog.LevelWarn.String():
		return slog.LevelWarn
	case slog.LevelError.String():
		return slog.LevelError
	default:
		if n, err := strconv.ParseInt(val, 10, 64); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "httplog: invalid integer value for %s: %s\n", env, val)
			os.Exit(1)
			return 0
		} else {
			return slog.Level(n)
		}
	}
}

type Func func(req *http.Request, elapsed time.Duration, status int)

// logRecord has a response writer and a status code
type logRecord struct {
	http.ResponseWriter
	status int
}

func (r *logRecord) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *logRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

// WriteHeader implements ResponseWriter for logRecord
func (r *logRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Wrap(f http.Handler, logFns ...Func) http.HandlerFunc {
	outLogger := log.New(os.Stdout, "", 0)
	errLogger := log.New(os.Stderr, "", 0)

	var fn Func
	if len(logFns) == 0 {
		fn = JSON(outLogger, errLogger)
	} else if len(logFns) == 1 {
		fn = logFns[0]
	} else {
		fn = func(req *http.Request, elapsed time.Duration, status int) {
			for _, lg := range logFns {
				lg(req, elapsed, status)
			}
		}
	}
	//it's a func!
	return func(w http.ResponseWriter, r *http.Request) {
		record := &logRecord{
			ResponseWriter: w,
		}

		start := time.Now()
		f.ServeHTTP(record, r)

		fn(r, time.Since(start), record.status)
	}
}
