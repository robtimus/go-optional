package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/robtimus/go-optional"
)

func TestMapChaining(t *testing.T) {
	opt1 := optional.Of(1)

	func1 := func(input int) string {
		return fmt.Sprintf("x%vx", input)
	}
	func2 := strings.ToUpper
	func3 := func(input string) int {
		return len(input)
	}

	opt2 := optional.Map(optional.Map(opt1, func1).Map(func2), func3)

	if opt2.OrElse(0) != 3 {
		t.Errorf("o2 should be present with value 3")
	}
}
