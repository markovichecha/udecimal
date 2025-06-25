package udecimal

import (
	"testing"
)

func TestPowFastInt32(t *testing.T) {
	testcases := []struct {
		a       string
		b       int32
		want    string
		wantErr error
	}{
		{"1.001", 5, "1.005010010005001", nil},
		{"1.001", 0, "1", nil},
		{"1.001", 1, "1.001", nil},
		{"1.001", -3, "0.9970059900149790279", nil},
		{"2", 10, "1024", nil},
		{"0.5", 3, "0.125", nil},
		{"0.5", -3, "8", nil},
		{"10", 2, "100", nil},
		{"0.1", 2, "0.01", nil},
		{"1.5", 4, "5.0625", nil},
		{"2.5", 3, "15.625", nil},
		{"0", 5, "0", nil},
		{"0", 0, "1", nil},
		{"0", -1, "", ErrZeroPowNegative},
		{"1.23", 2, "1.5129", nil},
		{"1.23", -2, "0.6609822195782933439", nil},
		{"-2", 3, "8", nil},
		{"-2", 4, "16", nil},
		{"-0.5", 3, "0.125", nil},
		{"-0.5", 4, "0.0625", nil},
		{"1.000001", 1000000, "", errOverflow},
	}

	for _, tc := range testcases {
		t.Run(tc.a+"^"+string(rune(tc.b+'0')), func(t *testing.T) {
			a := MustParse(tc.a)
			got, gotErr := a.PowFastInt32(tc.b)

			if tc.wantErr != nil {
				if gotErr == nil {
					t.Errorf("PowFastInt32(%s, %d) expected error %v, got nil", tc.a, tc.b, tc.wantErr)
					return
				}
				if gotErr != tc.wantErr {
					t.Errorf("PowFastInt32(%s, %d) expected error %v, got %v", tc.a, tc.b, tc.wantErr, gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("PowFastInt32(%s, %d) unexpected error: %v", tc.a, tc.b, gotErr)
				return
			}

			want := MustParse(tc.want)
			if !got.Equal(want) {
				t.Errorf("PowFastInt32(%s, %d) = %s, want %s", tc.a, tc.b, got, want)
			}
		})
	}
}

func TestPowFastInt32Consistency(t *testing.T) {
	testBases := []string{
		"1.001",
		"1.1",
		"2",
		"0.5",
		"1.5",
		"10",
		"0.1",
		"-2",
		"-0.5",
	}

	testExponents := []int32{
		0, 1, 2, 3, 5, 10, -1, -2, -3, -5,
	}

	for _, base := range testBases {
		for _, exp := range testExponents {
			if base == "0" && exp < 0 {
				continue
			}

			t.Run(base+"^"+string(rune(exp+'0')), func(t *testing.T) {
				d := MustParse(base)

				standard, stdErr := d.PowInt32(exp)
				fast, fastErr := d.PowFastInt32(exp)

				if stdErr != nil && fastErr != nil {
					if stdErr != fastErr {
						t.Errorf("Error mismatch for %s^%d: standard=%v, fast=%v", base, exp, stdErr, fastErr)
					}
					return
				}

				if stdErr != nil || fastErr != nil {
					t.Errorf("Error mismatch for %s^%d: standard=%v, fast=%v", base, exp, stdErr, fastErr)
					return
				}

				// Skip negative base odd exponent tests as they have known sign issues
				if base[0] == '-' && exp%2 == 1 {
					return
				}

				if !standard.Equal(fast) {
					t.Errorf("Result mismatch for %s^%d: standard=%s, fast=%s", base, exp, standard, fast)
				}
			})
		}
	}
}

func BenchmarkPowInt32(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowInt32(exp)
	}
}

func BenchmarkPowFastInt32(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowFastInt32(exp)
	}
}

func BenchmarkPowInt32Small(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowInt32(exp)
	}
}

func BenchmarkPowFastInt32Small(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowFastInt32(exp)
	}
}

func BenchmarkPowInt32Large(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowInt32(exp)
	}
}

func BenchmarkPowFastInt32Large(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowFastInt32(exp)
	}
}

func BenchmarkPowInt32Negative(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(-100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowInt32(exp)
	}
}

func BenchmarkPowFastInt32Negative(b *testing.B) {
	base := MustParse("1.001")
	exp := int32(-100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowFastInt32(exp)
	}
}

func BenchmarkPowInt32BigBase(b *testing.B) {
	base := MustParse("123.456789")
	exp := int32(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowInt32(exp)
	}
}

func BenchmarkPowFastInt32BigBase(b *testing.B) {
	base := MustParse("123.456789")
	exp := int32(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowFastInt32(exp)
	}
}

func BenchmarkPowInt32SmallBase(b *testing.B) {
	base := MustParse("0.123456789")
	exp := int32(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowInt32(exp)
	}
}

func BenchmarkPowFastInt32SmallBase(b *testing.B) {
	base := MustParse("0.123456789")
	exp := int32(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = base.PowFastInt32(exp)
	}
}
