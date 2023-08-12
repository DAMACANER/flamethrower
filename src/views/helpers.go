package views

import (
	"strings"

	"github.com/rivo/tview"
)

func notEmpty(v interface{}) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		return trimmed != "" && trimmed != "None"
	case *string:
		if v == nil {
			return false
		}
		trimmed := strings.TrimSpace(*v)
		return trimmed != "" && trimmed != "None"
	case *int, *int32, *int64:
		return v != nil
	case *float32, *float64:
		return v != nil
	default:
		return true
	}
}

// squashAndCenter takes a tview.Primitive component and returns a new tview.Flex 
// component that centers the given component both horizontally and vertically.
// The returned component has empty tview.Box components on the left, right, top, and bottom of the given component.
func squashAndCenter(component tview.Primitive) *tview.Flex {
	vf := tview.NewFlex().SetDirection(tview.FlexRow)
	vf.AddItem(tview.NewBox(), 0, 1, false)
	vf.AddItem(component, 0, 1, false)
	vf.AddItem(tview.NewBox(), 0, 1, false)
	hf := tview.NewFlex().SetDirection(tview.FlexColumn)
	hf.AddItem(tview.NewBox(), 0, 1, false)
	hf.AddItem(vf, 0, 1, false)
	hf.AddItem(tview.NewBox(), 0, 1, false)
	return hf
}