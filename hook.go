package gcfhook

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// This maps a logrus level to a Google severity.
var logrusToGoogleSeverityMap = map[logrus.Level]logging.Severity{
	logrus.PanicLevel: logging.Emergency,
	logrus.FatalLevel: logging.Alert,
	logrus.ErrorLevel: logging.Error,
	logrus.WarnLevel:  logging.Warning,
	logrus.InfoLevel:  logging.Info,
	logrus.DebugLevel: logging.Debug,
	logrus.TraceLevel: logging.Default,
}

// GoogleCloudFunctionHook is the logrus hook.
type GoogleCloudFunctionHook struct {
	logger      *logging.Logger
	executionID string
}

// New creates a new hook.
func New() (*GoogleCloudFunctionHook, error) {
	project := os.Getenv("GCP_PROJECT")
	function := os.Getenv("FUNCTION_NAME")
	region := os.Getenv("FUNCTION_REGION")

	if project == "" {
		return nil, fmt.Errorf("Failed to create logging client: GCP_PROJECT environment variable unset or missing")
	}
	if function == "" {
		return nil, fmt.Errorf("Failed to create logging client: FUNCTION_NAME environment variable unset or missing")
	}
	if region == "" {
		return nil, fmt.Errorf("Failed to create logging client: FUNCTION_REGION environment variable unset or missing")
	}

	ctx := context.Background()
	client, err := logging.NewClient(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("Failed to create logging client: %v", err)
	}

	res := monitoredres.MonitoredResource{
		Type:   "cloud_function",
		Labels: map[string]string{"region": region, "function_name": function},
	}

	hook := &GoogleCloudFunctionHook{
		logger: client.Logger("cloudfunctions.googleapis.com/cloud-functions", logging.CommonResource(&res)),
	}
	return hook, nil
}

// NewForRequest creates a new hook that will include the "execution ID" of the request
// as a label with each log message.
func NewForRequest(r *http.Request) (*GoogleCloudFunctionHook, error) {
	hook, err := New()
	if err != nil {
		return nil, err
	}

	// If we can get the execution ID, then use it when we fire the log messages.
	id := r.Header.Get("Function-Execution-Id")
	if id != "" {
		hook.executionID = id
	}

	return hook, nil
}

// Levels are the available logging levels.
func (hook *GoogleCloudFunctionHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

// Fire sends an entry.
func (hook *GoogleCloudFunctionHook) Fire(entry *logrus.Entry) error {
	severity := logging.Default
	if value, okay := logrusToGoogleSeverityMap[entry.Level]; okay {
		severity = value
	}

	labels := map[string]string{}
	if hook.executionID != "" {
		labels["execution_id"] = hook.executionID
	}

	hook.logger.Log(logging.Entry{Severity: severity, Payload: entry.Message, Labels: labels})
	return nil
}

// Flush flushes the logs.
//
// Call this before your Google Cloud Function ends.
func (hook *GoogleCloudFunctionHook) Flush() {
	hook.logger.Flush()
}
