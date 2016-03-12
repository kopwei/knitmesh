package common

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
)

type textFormatter struct {
}

// Based off logrus.TextFormatter, which behaves completely
// differently when you don't want colored output
func (f *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	levelText := strings.ToUpper(entry.Level.String())[0:4]
	timeStamp := entry.Time.Format("2006/01/02 15:04:05.000000")
	if len(entry.Data) > 0 {
		fmt.Fprintf(b, "%s: %s %-44s ", levelText, timeStamp, entry.Message)
		for k, v := range entry.Data {
			fmt.Fprintf(b, " %s=%v", k, v)
		}
	} else {
		// No padding when there's no fields
		fmt.Fprintf(b, "%s: %s %s", levelText, timeStamp, entry.Message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

var (
	standardTextFormatter = &textFormatter{}
)

// Log is the logging
var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.Formatter = standardTextFormatter
}

// SetLogLevel is used to set the log level
func SetLogLevel(levelname string) {
	level, err := logrus.ParseLevel(levelname)
	if err != nil {
		Log.Fatal(err)
	}
	Log.Level = level
}

// CheckFatal is used to check the fatal
func CheckFatal(e error) {
	if e != nil {
		Log.Fatal(e)
	}
}

// CheckWarn is used to check warn
func CheckWarn(e error) {
	if e != nil {
		Log.Warnln(e)
	}
}
