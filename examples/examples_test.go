package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/robtimus/go-optional"
)

func TestMapChaining(t *testing.T) {
	o1 := optional.Of(1)

	f1 := func(input int) string {
		return fmt.Sprintf("x%vx", input)
	}
	f2 := strings.ToUpper
	f3 := func(input string) int {
		return len(input)
	}

	o2 := optional.Map(optional.Map(o1, f1).Map(f2), f3)

	if o2.OrElse(0) != 3 {
		t.Errorf("o2 should be present with value 3")
	}
}
