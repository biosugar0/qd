package main

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestGetFileName(t *testing.T) {
	now, _ := time.Parse("2006-01-02", "2022-07-10")
	expect := struct {
		name string
		err  error
	}{
		name: "2022-07-10.md",
		err:  nil,
	}
	t.Run("getFileName", func(t *testing.T) {
		fname, err := getDailyNoteName(now)
		assert.Check(t, expect.err == err)
		assert.Check(t, expect.name == fname)
	})
}
