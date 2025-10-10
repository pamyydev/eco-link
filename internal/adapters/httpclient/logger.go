package httpclient

import (
	"log"
	"os"
	"time"
)

// Logger provides structured logging for HTTP clients
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewLogger creates a new structured logger
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
	}
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(method, url string, startTime time.Time) {
	duration := time.Since(startTime)
	l.infoLogger.Printf("HTTP %s %s - Duration: %v", method, url, duration)
}

// LogResponse logs HTTP response details
func (l *Logger) LogResponse(url string, statusCode int, duration time.Duration) {
	if statusCode >= 400 {
		l.errorLogger.Printf("HTTP Error %d for %s - Duration: %v", statusCode, url, duration)
	} else {
		l.infoLogger.Printf("HTTP %d for %s - Duration: %v", statusCode, url, duration)
	}
}

// LogError logs errors with context
func (l *Logger) LogError(operation, url string, err error) {
	l.errorLogger.Printf("Error in %s for URL '%s': %v", operation, url, err)
}

// LogRetry logs retry attempts
func (l *Logger) LogRetry(operation, url string, attempt int, maxRetries int, err error) {
	l.debugLogger.Printf("Retry attempt %d/%d for %s on URL '%s': %v", attempt, maxRetries, operation, url, err)
}

// LogValidationError logs validation errors
func (l *Logger) LogValidationError(field, value string, err error) {
	l.errorLogger.Printf("Validation error for field '%s' with value '%s': %v", field, value, err)
}

// LogAPICall logs API call details
func (l *Logger) LogAPICall(apiName, url string, startTime time.Time) {
	l.debugLogger.Printf("Calling %s API for URL: %s", apiName, url)
}

// LogAPISuccess logs successful API calls
func (l *Logger) LogAPISuccess(apiName, url string, duration time.Duration) {
	l.infoLogger.Printf("%s API call successful for URL '%s' - Duration: %v", apiName, url, duration)
}

// LogAPIFailure logs failed API calls
func (l *Logger) LogAPIFailure(apiName, url string, err error, duration time.Duration) {
	l.errorLogger.Printf("%s API call failed for URL '%s' - Error: %v - Duration: %v", apiName, url, err, duration)
}

// LogDataValidation logs data validation results
func (l *Logger) LogDataValidation(apiName string, data interface{}) {
	l.debugLogger.Printf("Data validation for %s API: %+v", apiName, data)
}

// LogInvalidData logs invalid data received from APIs
func (l *Logger) LogInvalidData(apiName, field string, value interface{}, reason string) {
	l.errorLogger.Printf("Invalid data from %s API - Field: %s, Value: %v, Reason: %s", apiName, field, value, reason)
}

// LogRetryExhausted logs when all retry attempts are exhausted
func (l *Logger) LogRetryExhausted(operation, url string, maxRetries int, finalError error) {
	l.errorLogger.Printf("All %d retry attempts exhausted for %s on URL '%s'. Final error: %v", maxRetries, operation, url, finalError)
}

// LogServiceOperation logs service-level operations
func (l *Logger) LogServiceOperation(operation, url string, startTime time.Time) {
	l.infoLogger.Printf("Starting %s for URL: %s", operation, url)
}

// LogServiceCompletion logs service operation completion
func (l *Logger) LogServiceCompletion(operation, url string, success bool, duration time.Duration) {
	if success {
		l.infoLogger.Printf("Completed %s for URL '%s' successfully - Duration: %v", operation, url, duration)
	} else {
		l.errorLogger.Printf("Failed %s for URL '%s' - Duration: %v", operation, url, duration)
	}
}
