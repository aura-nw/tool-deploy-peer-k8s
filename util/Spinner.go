package util

import (
	"time"

	"github.com/theckman/yacspin"
)

func CreateSpinner() *yacspin.Spinner {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		panic(err)
	}
	return spinner
}
