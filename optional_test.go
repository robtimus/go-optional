package optional

import (
	"fmt"
	"io"
	"testing"
)

func TestDefaultIsEmpty(t *testing.T) {
	var o Optional[string]

	if o.IsPresent() {
		t.Error("default Optional should not be present")
	}
	if !o.IsEmpty() {
		t.Error("defautl Optional should be empty")
	}
}

func TestEmpty(t *testing.T) {
	o := Empty[string]()

	if o.IsPresent() {
		t.Error("optional.Empty() should not be present")
	}
	if !o.IsEmpty() {
		t.Error("optional.Empty() should be empty")
	}
}

func TestOf(t *testing.T) {
	o := Of(1)

	if !o.IsPresent() {
		t.Error("optional.Of(1) should be present")
	}
	if o.IsEmpty() {
		t.Error("optional.Of(1) should not be empty")
	}
}

func TestOfNillableWithNil(t *testing.T) {
	o := OfNillable[int](nil)

	if o.IsPresent() {
		t.Error("optional.OfNillable(nil) should not be present")
	}
	if !o.IsEmpty() {
		t.Error("optional.OfNillable(nil) should be empty")
	}
}

func TestOfNillableWithNonNil(t *testing.T) {
	v := 1
	o := OfNillable(&v)

	if !o.IsPresent() {
		t.Error("optional.OfNillable(*1) should be present")
	}
	if o.IsEmpty() {
		t.Error("optional.OfNillable(*1) should not be empty")
	}
}

func TestIfPresentWhenEmpty(t *testing.T) {
	o := Empty[string]()

	action := capturingAction[string]{}

	o.IfPresent(action.Invoke)

	if len(action.arguments) != 0 {
		t.Errorf("action given to optional.Empty().IfPresent should not be invoked, was invoked with %v", action.arguments)
	}
}

func TestIfPresentWhenPresent(t *testing.T) {
	o := Of(1)

	action := capturingAction[int]{}

	o.IfPresent(action.Invoke)

	if len(action.arguments) != 1 || action.arguments[0] != 1 {
		t.Errorf("action given to optional.Of(1).IfPresent should be invoked with [1], was %v", action.arguments)
	}
}

func TestIfPresentOrElseWhenEmpty(t *testing.T) {
	o := Empty[string]()

	action := capturingAction[string]{}
	emptyAction := capturingNoArgAction{}

	o.IfPresentOrElse(action.Invoke, emptyAction.Invoke)

	if len(action.arguments) != 0 {
		t.Errorf("action given to optional.Empty().IfPresentOrElse should not be invoked, was invoked with %v", action.arguments)
	}
	if emptyAction.invocations != 1 {
		t.Errorf("emptyAction given to optional.Empty().IfPresentOrElse should be invoked once, #invocations: %v", emptyAction.invocations)
	}
}

func TestIfPresentOrElseWhenPresent(t *testing.T) {
	o := Of(1)

	action := capturingAction[int]{}
	emptyAction := capturingNoArgAction{}

	o.IfPresentOrElse(action.Invoke, emptyAction.Invoke)

	if len(action.arguments) != 1 || action.arguments[0] != 1 {
		t.Errorf("action given to optional.Of(1).IfPresentOrElse should be invoked with [1], was %v", action.arguments)
	}
	if emptyAction.invocations != 0 {
		t.Errorf("emptyAction given to optional.Empty().IfPresentOrElse should not be invoked, #invocations: %v", emptyAction.invocations)
	}
}

func TestFilterWhenEmpty(t *testing.T) {
	parameters := []bool{true, false}

	for i := range parameters {
		result := parameters[i]

		t.Run(fmt.Sprintf("Filter returns %v", result), func(t *testing.T) {
			o := Empty[string]()

			predicate := capturingPredicate[string]{result: parameters[i]}

			filtered := o.Filter(predicate.Invoke)

			if !filtered.IsEmpty() {
				t.Errorf("optional.Empty().Filter should return an empty Optional, was %v", filtered)
			}
			if len(predicate.arguments) != 0 {
				t.Errorf("predicate given to optional.Empty().Filter should not be invoked, was invoked with %v", predicate.arguments)
			}
		})
	}
}

func TestFilterWhenPresentAndPredicateReturnsTrue(t *testing.T) {
	o := Of(1)

	predicate := capturingPredicate[int]{result: true}

	filtered := o.Filter(predicate.Invoke)

	if filtered.IsEmpty() {
		t.Errorf("optional.Of(1).Filter should not return an empty Optional")
	}
	if filtered.value != o.value {
		t.Errorf("optional.Of(1).Filter should return an Optional with the same value (%v), was %v", o.value, filtered.value)
	}
	if len(predicate.arguments) != 1 || predicate.arguments[0] != 1 {
		t.Errorf("predicate given to optional.Empty().Filter should be invoked with [1], was %v", predicate.arguments)
	}
}

func TestFilterWhenPresentAndPredicateReturnsFalse(t *testing.T) {
	o := Of(1)

	predicate := capturingPredicate[int]{result: false}

	filtered := o.Filter(predicate.Invoke)

	if !filtered.IsEmpty() {
		t.Errorf("optional.Empty().Filter should return an empty Optional, was %v", filtered)
	}
	if len(predicate.arguments) != 1 || predicate.arguments[0] != 1 {
		t.Errorf("predicate given to optional.Empty().Filter should be invoked with [1], was %v", predicate.arguments)
	}
}

func TestMapWhenEmpty(t *testing.T) {
	o := Empty[string]()

	mapper := capturingFunction[string, string]{result: "foo"}

	mapped := o.Map(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Empty().Map should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to optional.Empty().Map should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestMapWhenPresent(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, int]{result: 2}

	mapped := o.Map(mapper.Invoke)

	if mapped.IsEmpty() {
		t.Errorf("optional.Of(1).Map should not return an empty Optional")
	}
	if *mapped.value != 2 {
		t.Errorf("optional.Of(2).Map should return an Optional with value 2, was %v", *mapped.value)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).Map should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalMapWhenEmpty(t *testing.T) {
	o := Empty[string]()

	mapper := capturingFunction[string, int]{result: 2}

	mapped := Map(o, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("Map called with optional.Empty() should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to Map should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestGlobalMapWhenPresent(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, string]{result: "foo"}

	mapped := Map(o, mapper.Invoke)

	if mapped.IsEmpty() {
		t.Errorf("Map called with optional.Of(1) should not return an empty Optional")
	}
	if *mapped.value != "foo" {
		t.Errorf("Map called with optional.Of(2) should return an Optional with value 'foo', was %v", *mapped.value)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to Map should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestMapNillableWhenEmpty(t *testing.T) {
	foo := "foo"
	parameters := []*string{nil, &foo}

	for i := range parameters {
		result := parameters[i]

		t.Run(fmt.Sprintf("Mapper returns %v", result), func(t *testing.T) {
			o := Empty[string]()

			mapper := capturingFunction[string, *string]{result: result}

			mapped := o.MapNillable(mapper.Invoke)

			if !mapped.IsEmpty() {
				t.Errorf("optional.Empty().MapNillable should return an empty Optional, was %v", mapped)
			}
			if len(mapper.arguments) != 0 {
				t.Errorf("mapper given to optional.Empty().MapNillable should not be invoked, was invoked with %v", mapper.arguments)
			}
		})
	}
}

func TestMapNillableWhenPresentReturningNil(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, *int]{result: nil}

	mapped := o.MapNillable(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Of(1).MapNillable should return an empty Optional if mapper returns nil, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).MapNillable should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestMapNillableWhenPresentReturningNonNil(t *testing.T) {
	o := Of(1)

	result := 2
	mapper := capturingFunction[int, *int]{result: &result}

	mapped := o.MapNillable(mapper.Invoke)

	if mapped.IsEmpty() {
		t.Errorf("optional.Of(1).MapNillable should not return an empty Optional if mapper returns non-nil")
	}
	if *mapped.value != 2 {
		t.Errorf("optional.Of(2).MapNillable should return an Optional with value 2, was %v", *mapped.value)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).Map should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalMapNillableWhenEmpty(t *testing.T) {
	foo := "foo"
	parameters := []*string{nil, &foo}

	for i := range parameters {
		result := parameters[i]

		t.Run(fmt.Sprintf("Mapper returns %v", result), func(t *testing.T) {
			o := Empty[int]()

			mapper := capturingFunction[int, *string]{result: result}

			mapped := MapNillable(o, mapper.Invoke)

			if !mapped.IsEmpty() {
				t.Errorf("MapNillable called with optional.Empty() should return an empty Optional, was %v", mapped)
			}
			if len(mapper.arguments) != 0 {
				t.Errorf("mapper given to MapNillable should not be invoked, was invoked with %v", mapper.arguments)
			}
		})
	}
}

func TestGlobalMapNillableWhenPresentReturningNil(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, *string]{result: nil}

	mapped := MapNillable(o, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("MapNillable called with optional.Of(1) should return an empty Optional if mapper returns nil, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to MapNillable should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalMapNillableWhenPresentReturningNonNil(t *testing.T) {
	o := Of(1)

	result := "foo"
	mapper := capturingFunction[int, *string]{result: &result}

	mapped := MapNillable(o, mapper.Invoke)

	if mapped.IsEmpty() {
		t.Errorf("MapNillable called with optional.Of(1) should not return an empty Optional if mapper returns non-nil")
	}
	if *mapped.value != "foo" {
		t.Errorf("MapNillable called with optional.Of(2) should return an Optional with value 2, was %v", *mapped.value)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to Map should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestFlatMapWhenEmpty(t *testing.T) {
	o := Empty[string]()

	mapper := capturingFunction[string, Optional[string]]{result: Of("foo")}

	mapped := o.FlatMap(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Empty().FlatMap should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to optional.Empty().FlatMap should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestFlatMapWhenPresentReturningEmpty(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, Optional[int]]{result: Empty[int]()}

	mapped := o.FlatMap(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Of(1).FlatMap should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).FlatMap should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestFlatMapWhenPresentReturningPresent(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, Optional[int]]{result: Of(2)}

	mapped := o.FlatMap(mapper.Invoke)

	if mapped.IsEmpty() {
		t.Errorf("optional.Of(1).FlatMap should not return an empty Optional")
	}
	if *mapped.value != 2 {
		t.Errorf("optional.Of(2).FlatMap should return an Optional with value 2, was %v", *mapped.value)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).FlatMap should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalFlatMapWhenEmpty(t *testing.T) {
	o := Empty[string]()

	mapper := capturingFunction[string, Optional[int]]{result: Of(1)}

	mapped := FlatMap(o, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("FlatMap called with optional.Empty() should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to FlatMap should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestGlobalFlatMapWhenPresentReturningEmpty(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, Optional[string]]{result: Empty[string]()}

	mapped := FlatMap(o, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("FlatMap called with optional.Of(1) should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to FlatMap should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalFlatMapWhenPresentReturningPresent(t *testing.T) {
	o := Of(1)

	mapper := capturingFunction[int, Optional[string]]{result: Of("foo")}

	mapped := FlatMap(o, mapper.Invoke)

	if mapped.IsEmpty() {
		t.Errorf("FlatMap called with optional.Of(1) should not return an empty Optional")
	}
	if *mapped.value != "foo" {
		t.Errorf("FlatMap called with optional.Of(2) should return an Optional with value 'foo', was %v", *mapped.value)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to FlatMap should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestOrWhenEmptyReturningEmpty(t *testing.T) {
	o := Empty[string]()

	supplier := capturingSupplier[Optional[string]]{result: Empty[string]()}

	o2 := o.Or(supplier.Invoke)

	if !o2.IsEmpty() {
		t.Errorf("optional.Empty().Or should return an empty Optional, was %v", o2)
	}
	if supplier.invocations == 0 {
		t.Errorf("supplier given to optional.Empty().Or should be invoked")
	}
}

func TestOrWhenEmptyReturningPresent(t *testing.T) {
	o := Empty[string]()

	supplier := capturingSupplier[Optional[string]]{result: Of("foo")}

	result := o.Or(supplier.Invoke)

	if result.IsEmpty() {
		t.Errorf("optional.Empty().Or should not return an empty Optional")
	}
	if *result.value != "foo" {
		t.Errorf("optional.Empty().Or should return an Optional with value 'foo', was %v", *result.value)
	}
	if supplier.invocations == 0 {
		t.Errorf("supplier given to optional.Empty().Or should be invoked")
	}
}

func TestOrWhenPresent(t *testing.T) {
	parameters := []Optional[int]{Empty[int](), Of(2)}

	for i := range parameters {
		result := parameters[i]

		o := Of(1)

		supplier := capturingSupplier[Optional[int]]{result: result}

		o2 := o.Or(supplier.Invoke)

		if o2.IsEmpty() {
			t.Errorf("optional.Of(1).Or should not return an empty Optional")
		}
		if *o2.value != 1 {
			t.Errorf("optional.Of(2).Or should return an Optional with value 2, was %v", *o2.value)
		}
		if supplier.invocations != 0 {
			t.Errorf("supplier given to optional.Of(1).Or should not be invoked, #invocations %v", supplier.invocations)
		}
	}
}

func TestSliceWhenEmpty(t *testing.T) {
	o := Empty[string]()

	slice := o.Slice()

	if len(slice) != 0 {
		t.Errorf("optional.Empty().Slice should return an empty slice, was %v", slice)
	}
}

func TestSliceWhenPresent(t *testing.T) {
	o := Of(1)

	slice := o.Slice()

	if len(slice) != 1 || slice[0] != 1 {
		t.Errorf("optional.Empty().Slice should return [1], was %v", slice)
	}
}

func TestOrElseWhenEmpty(t *testing.T) {
	o := Empty[string]()

	value := o.OrElse("foo")

	if value != "foo" {
		t.Errorf("optional.Empty().OrElse('foo') should return 'foo', was %v", value)
	}
}

func TestOrElseWhenPresent(t *testing.T) {
	o := Of(1)

	value := o.OrElse(2)

	if value != 1 {
		t.Errorf("optional.Of(1).OrElse(2) should return 1, was %v", value)
	}
}

func TestOrElseGetWhenEmpty(t *testing.T) {
	o := Empty[string]()

	supplier := capturingSupplier[string]{result: "foo"}

	value := o.OrElseGet(supplier.Invoke)

	if value != "foo" {
		t.Errorf("optional.Empty().OrElseGet(() => 'foo') should return 'foo', was %v", value)
	}
	if supplier.invocations != 1 {
		t.Errorf("supplier given to optional.Empty().OrElseGet should be invoked once, #invocations: %v", supplier.invocations)
	}
}

func TestOrElseGetWhenPresent(t *testing.T) {
	o := Of(1)

	supplier := capturingSupplier[int]{result: 2}

	value := o.OrElseGet(supplier.Invoke)

	if value != 1 {
		t.Errorf("optional.Of(1).OrElseGet(() => 2) should return 1, was %v", value)
	}
	if supplier.invocations != 0 {
		t.Errorf("supplier given to optional.Empty().OrElseGet should not be invoked, #invocations: %v", supplier.invocations)
	}
}

func TestOrElsePanicWhenEmpty(t *testing.T) {
	o := Empty[string]()

	defer func() {
		expectedMessage := "no value present"
		r := recover()
		if r == nil {
			t.Errorf("expected '%v', actual: nil", expectedMessage)
		} else if s, ok := r.(string); !ok || s != expectedMessage {
			t.Errorf("expected '%v', actual: %v", expectedMessage, r)
		}
	}()

	o.OrElsePanic()

	t.Error("expected an error")
}

func TestOrElsePanicWhenPresent(t *testing.T) {
	o := Of(1)

	value := o.OrElsePanic()

	if value != 1 {
		t.Errorf("optional.Of(1).OrElsePanic should return 1, was %v", value)
	}
}

func TestOrElseErrorWhenEmpty(t *testing.T) {
	o := Empty[string]()

	value, err := o.OrElseError()

	expectedMessage := "no value present"
	if value != "" {
		t.Errorf("optional.Empty().OrElseError should return an empty string, was %v", value)
	}
	if err == nil || err.Error() != expectedMessage {
		t.Errorf("optional.Empty().OrElseError should return an error with message '%v', was %v", expectedMessage, err)
	}
}

func TestOrElseErrorWhenPresent(t *testing.T) {
	o := Of(1)

	value, err := o.OrElseError()

	if value != 1 {
		t.Errorf("optional.Of(1).OrElseError should return 1, was %v", value)
	}
	if err != nil {
		t.Errorf("optional.Of(1).OrElseError should not return an error, was %v", err)
	}
}

func TestOrElseSupplyErrorWhenEmpty(t *testing.T) {
	o := Empty[string]()

	supplier := capturingSupplier[error]{result: io.EOF}

	value, err := o.OrElseSupplyError(supplier.Invoke)

	if value != "" {
		t.Errorf("optional.Empty().OrElseSupplyError should return an empty string, was %v", value)
	}
	if err != io.EOF {
		t.Errorf("optional.Empty().OrElseSupplyError should return io.EOF, was %v", err)
	}
	if supplier.invocations != 1 {
		t.Errorf("supplier given to optional.Empty().OrElseSupplyError should be invoked once, #invocations: %v", supplier.invocations)
	}
}

func TestOrElseSupplyErrorWhenPresent(t *testing.T) {
	o := Of(1)

	supplier := capturingSupplier[error]{result: io.EOF}

	value, err := o.OrElseSupplyError(supplier.Invoke)

	if value != 1 {
		t.Errorf("optional.Of(1).OrElseSupplyError should return 1, was %v", value)
	}
	if err != nil {
		t.Errorf("optional.Of(1).OrElseSupplyError should not return an error, was %v", err)
	}
	if supplier.invocations != 0 {
		t.Errorf("supplier given to optional.Empty().OrElseSupplyError should not be invoked, #invocations: %v", supplier.invocations)
	}
}

func TestStringWhenEmpty(t *testing.T) {
	o := Empty[string]()

	s := o.String()

	expectedValue := "Optional.empty"
	if s != expectedValue {
		t.Errorf("optional.Empty().String should return '%s', was %v", expectedValue, s)
	}
}

func TestStringWhenPresent(t *testing.T) {
	o := Of(1)

	s := o.String()

	expectedValue := "Optional[1]"
	if s != expectedValue {
		t.Errorf("optional.Of(1).String should return '%s', was %v", expectedValue, s)
	}
}

func TestEqualWithEmptyAndEmpty(t *testing.T) {
	o1 := Empty[string]()
	o2 := Empty[string]()

	if !Equal(o1, o2) {
		t.Errorf("%v and %v should be equal", o1, o2)
	}
}

func TestEqualWithEmptyAndPresent(t *testing.T) {
	o1 := Empty[string]()
	o2 := Of("foo")

	if Equal(o1, o2) {
		t.Errorf("%v and %v should not be equal", o1, o2)
	}
}

func TestEqualWithPresentAndEmpty(t *testing.T) {
	o1 := Of("foo")
	o2 := Empty[string]()

	if Equal(o1, o2) {
		t.Errorf("%v and %v should not be equal", o1, o2)
	}
}

func TestEqualWithPresentWithEqualValues(t *testing.T) {
	o1 := Of(1)
	o2 := Of(1)

	if !Equal(o1, o2) {
		t.Errorf("%v and %v should be equal", o1, o2)
	}
}

func TestEqualWithPresentWithUnequalValues(t *testing.T) {
	o1 := Of(1)
	o2 := Of(2)

	if Equal(o1, o2) {
		t.Errorf("%v and %v should not be equal", o1, o2)
	}
}

type capturingAction[T any] struct {
	arguments []T
}

func (f *capturingAction[T]) Invoke(input T) {
	f.arguments = append(f.arguments, input)
}

type capturingNoArgAction struct {
	invocations int
}

func (f *capturingNoArgAction) Invoke() {
	f.invocations++
}

type capturingPredicate[T any] struct {
	arguments []T
	result    bool
}

func (f *capturingPredicate[T]) Invoke(input T) bool {
	f.arguments = append(f.arguments, input)
	return f.result
}

type capturingFunction[T any, R any] struct {
	arguments []T
	result    R
}

func (f *capturingFunction[T, R]) Invoke(input T) R {
	f.arguments = append(f.arguments, input)
	return f.result
}

type capturingSupplier[T any] struct {
	invocations int
	result      T
}

func (f *capturingSupplier[T]) Invoke() T {
	f.invocations++
	return f.result
}
