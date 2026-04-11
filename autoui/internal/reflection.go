package internal

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/ebitenui/ebitenui/widget"
)

// GetRadioGroupElements returns the elements slice from a RadioGroup using reflection.
func GetRadioGroupElements(rg *widget.RadioGroup) []widget.RadioGroupElement {
	field := getPrivateField(rg, "elements")
	elements := make([]widget.RadioGroupElement, field.Len())
	for i := 0; i < field.Len(); i++ {
		elements[i] = field.Index(i).Interface().(widget.RadioGroupElement)
	}
	return elements
}

// GetTabBookTabs returns the tabs slice from a TabBook using reflection.
func GetTabBookTabs(tb *widget.TabBook) []*widget.TabBookTab {
	field := getPrivateField(tb, "tabs")
	tabs := make([]*widget.TabBookTab, field.Len())
	for i := 0; i < field.Len(); i++ {
		tabs[i] = field.Index(i).Interface().(*widget.TabBookTab)
	}
	return tabs
}

// GetTabBookTabLabel returns the label string from a TabBookTab using reflection.
func GetTabBookTabLabel(tab *widget.TabBookTab) string {
	field := getPrivateField(tab, "label")
	return field.String()
}

// getPrivateField returns a reflect.Value for a private field using unsafe.
// obj must be a pointer to the struct containing the field.
func getPrivateField(obj interface{}, fieldName string) reflect.Value {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("getPrivateField: obj must be a pointer, got %T", obj))
	}
	v = v.Elem()

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf("getPrivateField: field '%s' not found in %T", fieldName, obj))
	}

	// Use unsafe to bypass visibility
	// Create a new accessible value from the unsafe pointer
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
