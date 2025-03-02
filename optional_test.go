package optional

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestDefaultIsEmpty(t *testing.T) {
	var opt Optional[string]

	if opt.IsPresent() {
		t.Error("default Optional should not be present")
	}
	if !opt.IsEmpty() {
		t.Error("defautl Optional should be empty")
	}
}

func TestEmpty(t *testing.T) {
	opt := Empty[string]()

	if opt.IsPresent() {
		t.Error("optional.Empty() should not be present")
	}
	if !opt.IsEmpty() {
		t.Error("optional.Empty() should be empty")
	}
}

func TestOf(t *testing.T) {
	opt := Of(1)

	if !opt.IsPresent() {
		t.Error("optional.Of(1) should be present")
	}
	if opt.IsEmpty() {
		t.Error("optional.Of(1) should not be empty")
	}
}

func TestOfNillableWithNil(t *testing.T) {
	opt := OfNillable[int](nil)

	if opt.IsPresent() {
		t.Error("optional.OfNillable(nil) should not be present")
	}
	if !opt.IsEmpty() {
		t.Error("optional.OfNillable(nil) should be empty")
	}
}

func TestOfNillableWithNonNil(t *testing.T) {
	v := 1
	opt := OfNillable(&v)

	if !opt.IsPresent() {
		t.Error("optional.OfNillable(*1) should be present")
	}
	if opt.IsEmpty() {
		t.Error("optional.OfNillable(*1) should not be empty")
	}
}

func TestIfPresentWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	action := capturingAction[string]{}

	opt.IfPresent(action.Invoke)

	if len(action.arguments) != 0 {
		t.Errorf("action given to optional.Empty().IfPresent should not be invoked, was invoked with %v", action.arguments)
	}
}

func TestIfPresentWhenPresent(t *testing.T) {
	opt := Of(1)

	action := capturingAction[int]{}

	opt.IfPresent(action.Invoke)

	if len(action.arguments) != 1 || action.arguments[0] != 1 {
		t.Errorf("action given to optional.Of(1).IfPresent should be invoked with [1], was %v", action.arguments)
	}
}

func TestIfPresentOrElseWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	action := capturingAction[string]{}
	emptyAction := capturingNoArgAction{}

	opt.IfPresentOrElse(action.Invoke, emptyAction.Invoke)

	if len(action.arguments) != 0 {
		t.Errorf("action given to optional.Empty().IfPresentOrElse should not be invoked, was invoked with %v", action.arguments)
	}
	if emptyAction.invocations != 1 {
		t.Errorf("emptyAction given to optional.Empty().IfPresentOrElse should be invoked once, #invocations: %v", emptyAction.invocations)
	}
}

func TestIfPresentOrElseWhenPresent(t *testing.T) {
	opt := Of(1)

	action := capturingAction[int]{}
	emptyAction := capturingNoArgAction{}

	opt.IfPresentOrElse(action.Invoke, emptyAction.Invoke)

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
			opt := Empty[string]()

			predicate := capturingPredicate[string]{result: parameters[i]}

			filtered := opt.Filter(predicate.Invoke)

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
	opt := Of(1)

	predicate := capturingPredicate[int]{result: true}

	filtered := opt.Filter(predicate.Invoke)

	if filtered.IsEmpty() {
		t.Errorf("optional.Of(1).Filter should not return an empty Optional")
	}
	if filtered.value != opt.value {
		t.Errorf("optional.Of(1).Filter should return an Optional with the same value (%v), was %v", opt.value, filtered.value)
	}
	if len(predicate.arguments) != 1 || predicate.arguments[0] != 1 {
		t.Errorf("predicate given to optional.Empty().Filter should be invoked with [1], was %v", predicate.arguments)
	}
}

func TestFilterWhenPresentAndPredicateReturnsFalse(t *testing.T) {
	opt := Of(1)

	predicate := capturingPredicate[int]{result: false}

	filtered := opt.Filter(predicate.Invoke)

	if !filtered.IsEmpty() {
		t.Errorf("optional.Empty().Filter should return an empty Optional, was %v", filtered)
	}
	if len(predicate.arguments) != 1 || predicate.arguments[0] != 1 {
		t.Errorf("predicate given to optional.Empty().Filter should be invoked with [1], was %v", predicate.arguments)
	}
}

func TestMapWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	mapper := capturingFunction[string, string]{result: "foo"}

	mapped := opt.Map(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Empty().Map should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to optional.Empty().Map should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestMapWhenPresent(t *testing.T) {
	opt := Of(1)

	mapper := capturingFunction[int, int]{result: 2}

	mapped := opt.Map(mapper.Invoke)

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
	opt := Empty[string]()

	mapper := capturingFunction[string, int]{result: 2}

	mapped := Map(opt, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("Map called with optional.Empty() should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to Map should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestGlobalMapWhenPresent(t *testing.T) {
	opt := Of(1)

	mapper := capturingFunction[int, string]{result: "foo"}

	mapped := Map(opt, mapper.Invoke)

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
			opt := Empty[string]()

			mapper := capturingFunction[string, *string]{result: result}

			mapped := opt.MapNillable(mapper.Invoke)

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
	opt := Of(1)

	mapper := capturingFunction[int, *int]{result: nil}

	mapped := opt.MapNillable(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Of(1).MapNillable should return an empty Optional if mapper returns nil, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).MapNillable should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestMapNillableWhenPresentReturningNonNil(t *testing.T) {
	opt := Of(1)

	result := 2
	mapper := capturingFunction[int, *int]{result: &result}

	mapped := opt.MapNillable(mapper.Invoke)

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
			opt := Empty[int]()

			mapper := capturingFunction[int, *string]{result: result}

			mapped := MapNillable(opt, mapper.Invoke)

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
	opt := Of(1)

	mapper := capturingFunction[int, *string]{result: nil}

	mapped := MapNillable(opt, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("MapNillable called with optional.Of(1) should return an empty Optional if mapper returns nil, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to MapNillable should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalMapNillableWhenPresentReturningNonNil(t *testing.T) {
	opt := Of(1)

	result := "foo"
	mapper := capturingFunction[int, *string]{result: &result}

	mapped := MapNillable(opt, mapper.Invoke)

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
	opt := Empty[string]()

	mapper := capturingFunction[string, Optional[string]]{result: Of("foo")}

	mapped := opt.FlatMap(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Empty().FlatMap should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to optional.Empty().FlatMap should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestFlatMapWhenPresentReturningEmpty(t *testing.T) {
	opt := Of(1)

	mapper := capturingFunction[int, Optional[int]]{result: Empty[int]()}

	mapped := opt.FlatMap(mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("optional.Of(1).FlatMap should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to optional.Of(2).FlatMap should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestFlatMapWhenPresentReturningPresent(t *testing.T) {
	opt := Of(1)

	mapper := capturingFunction[int, Optional[int]]{result: Of(2)}

	mapped := opt.FlatMap(mapper.Invoke)

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
	opt := Empty[string]()

	mapper := capturingFunction[string, Optional[int]]{result: Of(1)}

	mapped := FlatMap(opt, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("FlatMap called with optional.Empty() should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 0 {
		t.Errorf("mapper given to FlatMap should not be invoked, was invoked with %v", mapper.arguments)
	}
}

func TestGlobalFlatMapWhenPresentReturningEmpty(t *testing.T) {
	opt := Of(1)

	mapper := capturingFunction[int, Optional[string]]{result: Empty[string]()}

	mapped := FlatMap(opt, mapper.Invoke)

	if !mapped.IsEmpty() {
		t.Errorf("FlatMap called with optional.Of(1) should return an empty Optional, was %v", mapped)
	}
	if len(mapper.arguments) != 1 || mapper.arguments[0] != 1 {
		t.Errorf("mapper given to FlatMap should be invoked with [1], was %v", mapper.arguments)
	}
}

func TestGlobalFlatMapWhenPresentReturningPresent(t *testing.T) {
	opt := Of(1)

	mapper := capturingFunction[int, Optional[string]]{result: Of("foo")}

	mapped := FlatMap(opt, mapper.Invoke)

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
	opt := Empty[string]()

	supplier := capturingSupplier[Optional[string]]{result: Empty[string]()}

	opt2 := opt.Or(supplier.Invoke)

	if !opt2.IsEmpty() {
		t.Errorf("optional.Empty().Or should return an empty Optional, was %v", opt2)
	}
	if supplier.invocations == 0 {
		t.Errorf("supplier given to optional.Empty().Or should be invoked")
	}
}

func TestOrWhenEmptyReturningPresent(t *testing.T) {
	opt := Empty[string]()

	supplier := capturingSupplier[Optional[string]]{result: Of("foo")}

	result := opt.Or(supplier.Invoke)

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

		opt := Of(1)

		supplier := capturingSupplier[Optional[int]]{result: result}

		opt2 := opt.Or(supplier.Invoke)

		if opt2.IsEmpty() {
			t.Errorf("optional.Of(1).Or should not return an empty Optional")
		}
		if *opt2.value != 1 {
			t.Errorf("optional.Of(2).Or should return an Optional with value 2, was %v", *opt2.value)
		}
		if supplier.invocations != 0 {
			t.Errorf("supplier given to optional.Of(1).Or should not be invoked, #invocations %v", supplier.invocations)
		}
	}
}

func TestSliceWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	slice := opt.Slice()

	if len(slice) != 0 {
		t.Errorf("optional.Empty().Slice should return an empty slice, was %v", slice)
	}
}

func TestSliceWhenPresent(t *testing.T) {
	opt := Of(1)

	slice := opt.Slice()

	if len(slice) != 1 || slice[0] != 1 {
		t.Errorf("optional.Empty().Slice should return [1], was %v", slice)
	}
}

func TestOrElseWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	value := opt.OrElse("foo")

	if value != "foo" {
		t.Errorf("optional.Empty().OrElse('foo') should return 'foo', was %v", value)
	}
}

func TestOrElseWhenPresent(t *testing.T) {
	opt := Of(1)

	value := opt.OrElse(2)

	if value != 1 {
		t.Errorf("optional.Of(1).OrElse(2) should return 1, was %v", value)
	}
}

func TestOrElseGetWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	supplier := capturingSupplier[string]{result: "foo"}

	value := opt.OrElseGet(supplier.Invoke)

	if value != "foo" {
		t.Errorf("optional.Empty().OrElseGet(() => 'foo') should return 'foo', was %v", value)
	}
	if supplier.invocations != 1 {
		t.Errorf("supplier given to optional.Empty().OrElseGet should be invoked once, #invocations: %v", supplier.invocations)
	}
}

func TestOrElseGetWhenPresent(t *testing.T) {
	opt := Of(1)

	supplier := capturingSupplier[int]{result: 2}

	value := opt.OrElseGet(supplier.Invoke)

	if value != 1 {
		t.Errorf("optional.Of(1).OrElseGet(() => 2) should return 1, was %v", value)
	}
	if supplier.invocations != 0 {
		t.Errorf("supplier given to optional.Empty().OrElseGet should not be invoked, #invocations: %v", supplier.invocations)
	}
}

func TestOrElsePanicWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	defer func() {
		expectedMessage := "no value present"
		r := recover()
		if r == nil {
			t.Errorf("expected '%v', actual: nil", expectedMessage)
		} else if s, ok := r.(string); !ok || s != expectedMessage {
			t.Errorf("expected '%v', actual: %v", expectedMessage, r)
		}
	}()

	opt.OrElsePanic()

	t.Error("expected an error")
}

func TestOrElsePanicWhenPresent(t *testing.T) {
	opt := Of(1)

	value := opt.OrElsePanic()

	if value != 1 {
		t.Errorf("optional.Of(1).OrElsePanic should return 1, was %v", value)
	}
}

func TestOrElseErrorWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	value, err := opt.OrElseError()

	expectedMessage := "no value present"
	if value != "" {
		t.Errorf("optional.Empty().OrElseError should return an empty string, was %v", value)
	}
	if err == nil || err.Error() != expectedMessage {
		t.Errorf("optional.Empty().OrElseError should return an error with message '%v', was %v", expectedMessage, err)
	}
}

func TestOrElseErrorWhenPresent(t *testing.T) {
	opt := Of(1)

	value, err := opt.OrElseError()

	if value != 1 {
		t.Errorf("optional.Of(1).OrElseError should return 1, was %v", value)
	}
	if err != nil {
		t.Errorf("optional.Of(1).OrElseError should not return an error, was %v", err)
	}
}

func TestOrElseSupplyErrorWhenEmpty(t *testing.T) {
	opt := Empty[string]()

	supplier := capturingSupplier[error]{result: io.EOF}

	value, err := opt.OrElseSupplyError(supplier.Invoke)

	if value != "" {
		t.Errorf("optional.Empty().OrElseSupplyError should return an empty string, was %v", value)
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("optional.Empty().OrElseSupplyError should return io.EOF, was %v", err)
	}
	if supplier.invocations != 1 {
		t.Errorf("supplier given to optional.Empty().OrElseSupplyError should be invoked once, #invocations: %v", supplier.invocations)
	}
}

func TestOrElseSupplyErrorWhenPresent(t *testing.T) {
	opt := Of(1)

	supplier := capturingSupplier[error]{result: io.EOF}

	value, err := opt.OrElseSupplyError(supplier.Invoke)

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
	opt := Empty[string]()

	s := opt.String()

	expectedValue := "Optional.empty"
	if s != expectedValue {
		t.Errorf("optional.Empty().String should return '%s', was %v", expectedValue, s)
	}
}

func TestStringWhenPresent(t *testing.T) {
	opt := Of(1)

	s := opt.String()

	expectedValue := "Optional[1]"
	if s != expectedValue {
		t.Errorf("optional.Of(1).String should return '%s', was %v", expectedValue, s)
	}
}

func TestEqualWithEmptyAndEmpty(t *testing.T) {
	opt1 := Empty[string]()
	opt2 := Empty[string]()

	if !Equal(opt1, opt2) {
		t.Errorf("%v and %v should be equal", opt1, opt2)
	}
}

func TestEqualWithEmptyAndPresent(t *testing.T) {
	opt1 := Empty[string]()
	opt2 := Of("foo")

	if Equal(opt1, opt2) {
		t.Errorf("%v and %v should not be equal", opt1, opt2)
	}
}

func TestEqualWithPresentAndEmpty(t *testing.T) {
	opt1 := Of("foo")
	opt2 := Empty[string]()

	if Equal(opt1, opt2) {
		t.Errorf("%v and %v should not be equal", opt1, opt2)
	}
}

func TestEqualWithPresentWithEqualValues(t *testing.T) {
	opt1 := Of(1)
	opt2 := Of(1)

	if !Equal(opt1, opt2) {
		t.Errorf("%v and %v should be equal", opt1, opt2)
	}
}

func TestEqualWithPresentWithUnequalValues(t *testing.T) {
	opt1 := Of(1)
	opt2 := Of(2)

	if Equal(opt1, opt2) {
		t.Errorf("%v and %v should not be equal", opt1, opt2)
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
