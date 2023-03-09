package accesslog

import (
	"api_gateway/internal/gateway/config"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct {
	handlers.RecoveryHandlerLogger
}

func NewHandler(cfg config.AccessLog, next http.Handler) (http.Handler, error) {
	httpLogFd, err := openAccessLogFile(cfg.HttpLogPath)
	if err != nil {
		return nil, err
	}
	handle := handlers.CombinedLoggingHandler(httpLogFd, next)

	return handle, nil
}

func openAccessLogFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log path %s: %w", dir, err)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o664)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}

	return file, nil
}
