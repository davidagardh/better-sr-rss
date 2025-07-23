package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationFmt(t *testing.T) {
	d := time.Duration(30*time.Hour + 5*time.Second)
	res := (&(Episode{Duration: d})).DurationFmt()
	require.Equal(t, "30:00:05", res)
}
