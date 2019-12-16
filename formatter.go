package gcfhook

import (
	"github.com/sirupsen/logrus"
)

// NullFormatter will format all entries as empty bytes so that nothing is ever printed
// to the screen.
type NullFormatter struct{}

// Format renders a single log entry.
func (f *NullFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte{}, nil
}
