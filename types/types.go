package types

import (
	"sync"
	"time"
)

// Hook allows optional observability without adding dependencies.
// Provide your own implementation and register with SetHook.
type Hook interface {
	// OnLoad is called after loading from a file or source.
	OnLoad(source string, keys int)
	// OnGet is called on each read. Duration is the total time spent.
	OnGet(key string, ok bool, err error, dur time.Duration)
}

var (
	// hookMu protects hook.
	hookMu sync.RWMutex
	// hook is the global hook instance.
	hook Hook
)

// SetHook installs a global hook. It is safe to call at program init.
//
// Parameters:
//   - h: The hook to install.
func SetHook(h Hook) {
	hookMu.Lock()
	defer hookMu.Unlock()
	hook = h
}

// CallOnLoad calls the OnLoad hook.
func CallOnLoad(source string, keys int) {
	hookMu.RLock()
	defer hookMu.RUnlock()
	if hook != nil {
		hook.OnLoad(source, keys)
	}
}

// CallOnGet calls the OnGet hook.
func CallOnGet(key string, ok bool, err error, d time.Duration) {
	hookMu.RLock()
	defer hookMu.RUnlock()
	if hook != nil {
		hook.OnGet(key, ok, err, d)
	}
}
