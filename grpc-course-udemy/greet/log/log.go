package log

import (
	"os"

	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
)

func init() {

}

// Here you can configure different levels of logs and also output formats for the logs
//(example: plain text, json, etc ...)
// Logger ...
func Logger(file string) *apexLog.Entry {

	apexLog.SetHandler(cli.Default)
	apexLog.SetLevel(apexLog.InfoLevel)
	apexLog.SetHandler(json.New(os.Stderr))

	ctx := apexLog.WithFields(apexLog.Fields{
		"file": file,
	})

	return ctx

}

// LoggerJSON ...
func LoggerJSON() *apexLog.Logger {

	l := &apexLog.Logger{
		Handler: json.New(os.Stderr),
		Level:   apexLog.DebugLevel,
	}

	return l
}
