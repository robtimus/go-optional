package optional

import (
	"errors"
	"fmt"
	"log"
)

const noValuePresentMessage = "no value present"

var errNoValuePresent = errors.New(noValuePresentMessage)

// Optional is a container object that may or may not contain a value.
// If no value is present, the object is considered empty.
type Optional[T any] struct {
	value *T
}

// Empty returns an empty Optional.
func Empty[T any]() Optional[T] {
	return Optional[T]{value: nil}
}

// Of returns a non-empty Optional describing the given value.
func Of[T any](value T) Optional[T] {
	return OfNillable(&value)
}

// OfNillable returns a non-empty Optional describing the target of the given pointer if it's not nil,
// or an empty Optional otherwise.
func OfNillable[T any](value *T) Optional[T] {
	return Optional[T]{value: value}
}

// Omit Get, require one of the OrElse methods instead

// IsPresent returns true if a value is present, or false otherwise.
func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

// IsEmpty returns true if no value is present, or false otherwise.
func (o Optional[T]) IsEmpty() bool {
	return o.value == nil
}

// IfPresent calls the given action with the value if present, or does nothing otherwise.
func (o Optional[T]) IfPresent(action func(value T)) {
	if o.value != nil {
		action(*o.value)
	}
}

// IfPresentOrElse calls the given action with the value if present, or calls the given empty-based action otherwise.
func (o Optional[T]) IfPresentOrElse(action func(value T), emptyAction func()) {
	if o.value != nil {
		action(*o.value)
	} else {
		emptyAction()
	}
}

// Filter returns a non-empty Optional if a value is present and it matches the given predicate, or an empty Optional otherwise.
func (o Optional[T]) Filter(predicate func(value T) bool) Optional[T] {
	if o.value == nil || predicate(*o.value) {
		return o
	}

	return Empty[T]()
}

// Map returns a non-empty Optional containing the result of calling the given mapper function on the value if present, or an empty Optional otherwise.
//
// Due to the limitations of generics in Go, the mapper function must return the Optional's generic type.
// The [Map] function can be used to map to different types.
func (o Optional[T]) Map(mapper func(value T) T) Optional[T] {
	if o.value == nil {
		return o
	}

	return Of(mapper(*o.value))
}

// Map returns a non-empty Optional containing the result of calling the given mapper function on the given Optional's value if present, or an empty Optional otherwise.
//
// This function can be used where the generic type of the Optional and the mapper function's return type do not match.
func Map[T any, U any](optional Optional[T], mapper func(value T) U) Optional[U] {
	if optional.value == nil {
		return Empty[U]()
	}

	return Of(mapper(*optional.value))
}

// MapNillable returns a possibly empty Optional (as if by [OfNillable]) based on the result of calling the given mapper function on the value if present,
// or an empty Optional otherwise.
//
// Due to the limitations of generics in Go, the mapper function must return a pointer to the Optional's generic type.
// The [MapNillable] function can be used to map to different types.
func (o Optional[T]) MapNillable(mapper func(value T) *T) Optional[T] {
	if o.value == nil {
		return o
	}

	return OfNillable(mapper(*o.value))
}

// MapNillable returns a possibly empty Optional (as if by [OfNillable]) based on the result of calling the given mapper function on the given Optional's value if present,
// or an empty Optional otherwise.
//
// This function can be used where the generic type of the Optional and the mapper function's return type do not match.
func MapNillable[T any, U any](optional Optional[T], mapper func(value T) *U) Optional[U] {
	if optional.value == nil {
		return Empty[U]()
	}

	return OfNillable(mapper(*optional.value))
}

// FlatMap returns the result of applying the given function if the value is present, or an empty Optional otherwise.
//
// Due to the limitations of generics in Go, the mapper function must return the Optional's exact type.
// The [FlatMap] function can be used to map to different types.
func (o Optional[T]) FlatMap(mapper func(value T) Optional[T]) Optional[T] {
	if o.value == nil {
		return o
	}

	return mapper(*o.value)
}

// FlatMap returns the result of applying the given function if the given Optional's value is present, or an empty Optional otherwise.
//
// This function can be used where the generic type of the Optional and the mapper function's return type do not match.
func FlatMap[T any, U any](optional Optional[T], mapper func(value T) Optional[U]) Optional[U] {
	if optional.value == nil {
		return Empty[U]()
	}

	return mapper(*optional.value)
}

// Or returns the Optional if the value is present, or the result of calling the given function otherwise.
func (o Optional[T]) Or(supplier func() Optional[T]) Optional[T] {
	if o.value != nil {
		return o
	}

	return supplier()
}

// Slice returns a slice containing the value if present, or an empty slice otherwise.
func (o Optional[T]) Slice() []T {
	var result []T
	if o.value != nil {
		result = append(result, *o.value)
	}

	return result
}

// OrElse returns the value if present, or the given other value otherwise.
func (o Optional[T]) OrElse(other T) T {
	if o.value != nil {
		return *o.value
	}

	return other
}

// OrElseGet returns the value if present, or the result of calling the given function otherwise.
func (o Optional[T]) OrElseGet(supplier func() T) T {
	if o.value != nil {
		return *o.value
	}

	return supplier()
}

// OrElsePanic returns the value if present, or panics otherwise.
func (o Optional[T]) OrElsePanic() T {
	if o.value == nil {
		log.Panic(noValuePresentMessage)
	}

	return *o.value
}

// OrElseError returns the value if present. If the Optional is empty it will return a non-nil error.
func (o Optional[T]) OrElseError() (T, error) {
	if o.value == nil {
		var zero T

		return zero, errNoValuePresent
	}

	return *o.value, nil
}

// OrElseSupplyError returns the value if present. If the Optional is empty it will return an error returned by the given supplier.
func (o Optional[T]) OrElseSupplyError(errorSupplier func() error) (T, error) {
	if o.value == nil {
		var zero T

		return zero, errorSupplier()
	}

	return *o.value, nil
}

// String implements the [fmt.Stringer] interface.
func (o Optional[T]) String() string {
	if o.value == nil {
		return "Optional.empty"
	}

	return fmt.Sprintf("Optional[%v]", *o.value)
}

// Equal compares two Optional objects. It will return true if both Optionals are empty, or if both Optionals have equal values.
func Equal[T comparable](opt Optional[T], other Optional[T]) bool {
	if opt.value == nil && other.value == nil {
		return true
	}

	if opt.value == nil || other.value == nil {
		return false
	}

	return *opt.value == *other.value
}
