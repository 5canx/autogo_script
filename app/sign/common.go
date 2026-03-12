package sign

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/motion"
)

// EnsureForeground: if already in target package, keep; else enable and launch.
func EnsureForeground(pkg string) bool {
	pkg = strings.TrimSpace(pkg)
	if pkg == "" {
		return false
	}
	if strings.TrimSpace(app.CurrentPackage()) == pkg {
		return true
	}
	app.Enable(pkg)
	if app.Launch(pkg, 0) {
		return true
	}
	fmt.Printf("[warn] 无法切到前台: %s\n", pkg)
	return false
}

// PageDetector describes a page and how to detect it.
// Detect should be fast and side-effect free.
type PageDetector struct {
	Name   string
	Detect func() bool
}

// SilentPoll polls fn until it returns true or timeout.
// It does not print logs; caller decides what to log.
func SilentPoll(timeout, interval time.Duration, fn func() bool) bool {
	if timeout <= 0 {
		return fn()
	}
	if interval <= 0 {
		interval = 300 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
		time.Sleep(interval)
	}
	return fn()
}

// BackUntil tries to reach any known page by pressing Back.
// Returns the matched page name and whether matched.
func BackUntil(pages []PageDetector, maxBack int, settle time.Duration) (string, bool) {
	if maxBack < 0 {
		maxBack = 0
	}
	if settle <= 0 {
		settle = 900 * time.Millisecond
	}

	match := func() (string, bool) {
		for _, p := range pages {
			if p.Detect != nil && p.Detect() {
				return p.Name, true
			}
		}
		return "", false
	}

	if name, ok := match(); ok {
		return name, true
	}

	for i := 0; i < maxBack; i++ {
		motion.Back(0)
		time.Sleep(settle)
		if name, ok := match(); ok {
			return name, true
		}
	}
	return "", false
}

// ClickCenter clicks the center of a rectangle.
func ClickCenter(x1, y1, x2, y2 int, delay time.Duration) {
	x := (x1 + x2) / 2
	y := (y1 + y2) / 2
	motion.Click(x, y, 0, 0)
	if delay > 0 {
		time.Sleep(delay)
	}
}
