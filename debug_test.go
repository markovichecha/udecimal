package udecimal

import (
	"fmt"
	"testing"
)

func TestDebugPowFast(t *testing.T) {
	// Test simple case: 2^2 = 4
	a := MustParse("2")
	fmt.Printf("Input decimal: coef=%v, prec=%d, neg=%v\n", a.coef.u128, a.prec, a.neg)

	// Check if it overflows
	if a.coef.overflow() {
		t.Fatalf("Input overflows")
	}

	base := a.coef.u128
	fmt.Printf("Base u128: hi=%d, lo=%d\n", base.hi, base.lo)

	result := one128
	fmt.Printf("One128: hi=%d, lo=%d\n", result.hi, result.lo)

	// Check comparison
	cmp := base.Cmp(result)
	fmt.Printf("Base.Cmp(One) = %d\n", cmp)

	if base.Cmp(result) >= 0 {
		fmt.Println("Base >= result, need to invert")
		quotient, remainder, err := max128.QuoRem(base)
		fmt.Printf("max128.QuoRem(base) = q=%v, r=%v, err=%v\n", quotient, remainder, err)
		if err != nil {
			t.Fatalf("QuoRem failed: %v", err)
		}
	}

	// Try multiplication
	squared, err := base.Mul(base)
	fmt.Printf("base.Mul(base) = %v, err=%v\n", squared, err)
	if err != nil {
		t.Fatalf("Multiplication failed: %v", err)
	}

	// Test actual function
	got, gotErr := a.PowFastInt32(2)
	fmt.Printf("PowFastInt32(2, 2) = %v, err=%v\n", got, gotErr)

	// Compare with standard
	standard, stdErr := a.PowInt32(2)
	fmt.Printf("PowInt32(2, 2) = %v, err=%v\n", standard, stdErr)
}
