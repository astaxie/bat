package main

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"
)

const (
	DEFAULT_REFRESH_RATE = time.Millisecond * 200
	FORMAT               = "[=>-]"
)

type ProgressBar struct {
	current                          int64
	Total                            int64
	RefreshRate                      time.Duration
	ShowPercent, ShowCounters        bool
	ShowSpeed, ShowTimeLeft, ShowBar bool
	ShowFinalTime                    bool

	isFinish     int32
	startTime    time.Time
	currentValue int64

	BarStart string
	BarEnd   string
	Empty    string
	Current  string
	CurrentN string
}

func NewProgressBar(total int64) (pb *ProgressBar) {
	pb = &ProgressBar{
		Total:         total,
		RefreshRate:   DEFAULT_REFRESH_RATE,
		ShowPercent:   true,
		ShowBar:       true,
		ShowCounters:  true,
		ShowFinalTime: true,
		ShowTimeLeft:  true,
		ShowSpeed:     true,
		BarStart:      "[",
		BarEnd:        "]",
		Empty:         "_",
		Current:       "=",
		CurrentN:      ">",
	}
	return
}

func (pb *ProgressBar) Start() {
	pb.startTime = time.Now()
	if pb.Total == 0 {
		pb.ShowBar = false
		pb.ShowTimeLeft = false
		pb.ShowPercent = false
	}
	go pb.writer()
}

// Write the current state of the progressbar
func (pb *ProgressBar) Update() {
	c := atomic.LoadInt64(&pb.current)
	if c != pb.currentValue {
		pb.write(c)
		pb.currentValue = c
	}
}

// Internal loop for writing progressbar
func (pb *ProgressBar) writer() {
	for {
		if atomic.LoadInt32(&pb.isFinish) != 0 {
			break
		}
		pb.Update()
		time.Sleep(pb.RefreshRate)
	}
}

// Increment current value
func (pb *ProgressBar) Increment() int {
	return pb.Add(1)
}

// Set current value
func (pb *ProgressBar) Set(current int) {
	atomic.StoreInt64(&pb.current, int64(current))
}

// Add to current value
func (pb *ProgressBar) Add(add int) int {
	return int(pb.Add64(int64(add)))
}

func (pb *ProgressBar) Add64(add int64) int64 {
	return atomic.AddInt64(&pb.current, add)
}

// End print
func (pb *ProgressBar) Finish() {
	atomic.StoreInt32(&pb.isFinish, 1)
	pb.write(atomic.LoadInt64(&pb.current))
}

// implement io.Writer
func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.Add(n)
	return
}

func (pb *ProgressBar) write(current int64) {
	width := 123

	var percentBox, countersBox, timeLeftBox, speedBox, barBox, end, out string

	// percents
	if pb.ShowPercent {
		percent := float64(current) / (float64(pb.Total) / float64(100))
		percentBox = fmt.Sprintf(" %#.02f %% ", percent)
	}

	// counters
	if pb.ShowCounters {
		if pb.Total > 0 {
			countersBox = fmt.Sprintf("%s / %s ", FormatBytes(current), FormatBytes(pb.Total))
		} else {
			countersBox = FormatBytes(current) + " "
		}
	}

	// time left
	fromStart := time.Now().Sub(pb.startTime)
	if atomic.LoadInt32(&pb.isFinish) != 0 {
		if pb.ShowFinalTime {
			left := (fromStart / time.Second) * time.Second
			timeLeftBox = left.String()
		}
	} else if pb.ShowTimeLeft && current > 0 {
		perEntry := fromStart / time.Duration(current)
		left := time.Duration(pb.Total-current) * perEntry
		left = (left / time.Second) * time.Second
		timeLeftBox = left.String()
	}

	// speed
	if pb.ShowSpeed && current > 0 {
		fromStart := time.Now().Sub(pb.startTime)
		speed := float64(current) / (float64(fromStart) / float64(time.Second))
		speedBox = FormatBytes(int64(speed)) + "/s "
	}

	// bar
	if pb.ShowBar {
		size := width - len(countersBox+pb.BarStart+pb.BarEnd+percentBox+timeLeftBox+speedBox)
		if size > 0 {
			curCount := int(math.Ceil((float64(current) / float64(pb.Total)) * float64(size)))
			emptCount := size - curCount
			barBox = pb.BarStart
			if emptCount < 0 {
				emptCount = 0
			}
			if curCount > size {
				curCount = size
			}
			if emptCount <= 0 {
				barBox += strings.Repeat(pb.Current, curCount)
			} else if curCount > 0 {
				barBox += strings.Repeat(pb.Current, curCount-1) + pb.CurrentN
			}

			barBox += strings.Repeat(pb.Empty, emptCount) + pb.BarEnd
		}
	}

	// check len
	out = countersBox + barBox + percentBox + speedBox + timeLeftBox
	if len(out) < width {
		end = strings.Repeat(" ", width-len(out))
	}

	// and print!
	fmt.Print("\r" + out + end)
}
