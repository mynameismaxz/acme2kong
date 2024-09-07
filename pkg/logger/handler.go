package logger

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"
)

// JSONHandler is a custom handler that formats logs in JSON and writes to multiple outputs.
type JSONHandler struct {
	mu     sync.Mutex  // Ensures thread safety for concurrent logging
	writer *os.File    // Write to stdout
	attrs  []slog.Attr // Attributes for structured logging
	groups []string    // Groups for nested attributes
}

// NewJSONHandler initializes and returns a new JSONHandler instance.
func NewJSONHandler(output *os.File) *JSONHandler {
	return &JSONHandler{
		writer: output,
		attrs:  make([]slog.Attr, 0),
		groups: make([]string, 0),
	}
}

// Handle processes the log record and outputs it in JSON format.
func (h *JSONHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Build the log message structure
	logMsg := map[string]interface{}{
		"level": r.Level.String(),
		"date":  time.Now().Format("2006-01-02 15:04:05"),
		"msg":   r.Message,
	}

	// Add attributes and groups to the log message
	for _, attr := range h.attrs {
		logMsg[attr.Key] = attr.Value
	}

	// Convert the log message to JSON
	jsonData, err := json.Marshal(logMsg)
	if err != nil {
		return err
	}

	// Write JSON data to stdout
	_, err = h.writer.Write(jsonData)
	if err != nil {
		return err
	}

	// Write a newline after the JSON object
	_, err = h.writer.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}

// Enabled returns whether the handler is enabled or not.
func (h *JSONHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

// WithAttrs returns a new handler with additional attributes.
func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	newHandler := &JSONHandler{
		writer: h.writer,
		attrs:  append(h.attrs, attrs...), // Merge existing and new attributes
		groups: h.groups,
	}

	return newHandler
}

// WithGroup returns a new handler with a new group context.
func (h *JSONHandler) WithGroup(name string) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	newHandler := &JSONHandler{
		writer: h.writer,
		attrs:  h.attrs,
		groups: append(h.groups, name), // Add the new group to the existing groups
	}

	return newHandler
}

// Close closes the JSONHandler output.
func (h *JSONHandler) Close() error {
	return h.writer.Close()
}
