package zconfig

import (
	"fmt"
	"reflect"
)

type Initializable interface {
	Init() error
}

// Used for type comparison.
var typeInitializable = reflect.TypeOf((*Initializable)(nil)).Elem()

func Initialize(field *Field) error {
	// Not initializable, nothing to do.
	if !field.Value.Type().Implements(typeInitializable) {
		return nil
	}

	// Initialize the element itself via the interface.
	err := field.Value.Interface().(Initializable).Init()
	if err != nil {
		return fmt.Errorf("initializing field: %s", err)
	}

	return nil
}
