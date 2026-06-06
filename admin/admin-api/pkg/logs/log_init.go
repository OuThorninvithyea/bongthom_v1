package logs

import (

	// Community packages
	"fmt"
	"github.com/rs/zerolog"
)

func NewLog(logLevel string) {
	switch logLevel {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "tracing":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	fmt.Println("log level", zerolog.GlobalLevel())
}
