package chrono

import (
	"github.com/dop251/goja"
	"time"
)

// Chrono struct
type Chrono struct {
	VM   *goja.Runtime
	Fn   goja.Callable
	This goja.Value
}

// New - new chrono instance
func New() (chrono Chrono, err error) {
	bytes, err := Asset("chrono.out.js")
	if err != nil {
		return chrono, err
	}

	vm := goja.New()

	prg, err := goja.Compile("chrono.js", string(bytes), false)
	if err != nil {
		return chrono, err
	}

	_, err = vm.RunProgram(prg)
	if err != nil {
		return chrono, err
	}

	fn, ok := goja.AssertFunction(vm.Get("chrono"))
	if !ok {
		return chrono, err
	}

	return Chrono{
		This: vm.ToValue(map[string]interface{}{}),
		VM:   vm,
		Fn:   fn,
	}, nil
}

// ParseDate - naturally parse the time
func (c *Chrono) ParseDate(expr string, now time.Time) (t *time.Time, err error) {
	v, err := c.Fn(c.This, c.VM.ToValue(expr), c.VM.ToValue(now.Format(time.RFC3339)))
	if err != nil {
		return t, err
	}

	switch o := v.Export().(type) {
	case time.Time:
		return &o, nil
	case nil:
		return nil, nil
	default:
		return t, err
	}
}
