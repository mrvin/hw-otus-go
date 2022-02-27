package main

import (
	"bytes"
	"io"
	"testing"
	"time"
)

type saved struct {
	out        io.Writer
	getTimeNTP func(string) (time.Time, error)
	getTimeNow func() time.Time
}

func TestHelloNow(t *testing.T) {
	const layout = "2 Jan 2006 15:04:05"
	const expected = "current time: 1945-05-09 10:03:00 +0000 UTC\nexact time: 1945-05-09 10:03:02 +0000 UTC\n"

	// Save and restore original: output, getTimeNTP, getTimeNow.
	sevedAll := saved{out, getTimeNTP, getTimeNow}
	defer func() {
		out = sevedAll.out
		getTimeNTP = sevedAll.getTimeNTP
		getTimeNow = sevedAll.getTimeNow

	}()

	out = new(bytes.Buffer) // captured output

	// Install the test's fake getTimeNTP.
	getTimeNTP = func(srvNTP string) (time.Time, error) {
		return time.Parse(layout, "9 May 1945 10:03:02")
	}

	// Install the test's fake getTimeNow.
	getTimeNow = func() time.Time {
		currentTime, _ := time.Parse(layout, "9 May 1945 10:03:00")
		return currentTime
	}

	main()

	got := out.(*bytes.Buffer).String()
	if got != expected {
		t.Fatalf("invalid output:\n%s, expected:\n%s", got, expected)
	}
}
