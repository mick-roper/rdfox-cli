package logging

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func New(level string) *zap.Logger {
	rawJSON := fmt.Sprintf(`{
	  "level": "%s",
	  "encoding": "console",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`, level)

	var cfg zap.Config
	if err := json.Unmarshal([]byte(rawJSON), &cfg); err != nil {
		panic(err)
	}

	if strings.EqualFold(level, "debug") {
		cfg.EncoderConfig.StacktraceKey = "stack-trace"
	} else {
		cfg.EncoderConfig.StacktraceKey = ""
	}

	return zap.Must(cfg.Build())
}
