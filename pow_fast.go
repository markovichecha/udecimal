package udecimal

const (
	maxExponential = uint32(0x80000)
)

func (d Decimal) PowFastInt32(e int32) (Decimal, error) {
	// Handle zero base cases
	if d.IsZero() {
		if e < 0 {
			return Decimal{}, ErrZeroPowNegative
		}
		if e == 0 {
			return One, nil
		}
		return Zero, nil
	}

	invert := e < 0

	if e == 0 {
		return One, nil
	}

	exp := uint32(e)
	if invert {
		exp = uint32(-e)
	}

	if exp >= maxExponential {
		return Zero, errOverflow
	}

	// Handle sign - negative base with odd exponent
	originalExp := e
	if invert {
		originalExp = -e
	}
	neg := d.neg && (originalExp%2 == 1)
	squaredBase := d.Abs()
	result := One

	// Optimize for large bases by inverting
	if squaredBase.GreaterThanOrEqual(One) {
		var err error
		squaredBase, err = One.Div(squaredBase)
		if err != nil {
			return Zero, err
		}
		invert = !invert
	}

	// Binary exponentiation with unrolled loop (19 bits)
	if exp&0x1 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x2 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x4 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x8 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x10 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x20 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x40 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x80 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x100 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x200 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x400 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x800 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x1000 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x2000 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x4000 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x8000 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x10000 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x20000 > 0 {
		result = result.Mul(squaredBase)
	}

	squaredBase = squaredBase.Mul(squaredBase)

	if exp&0x40000 > 0 {
		result = result.Mul(squaredBase)
	}

	if result.IsZero() {
		return Zero, errOverflow
	}

	if invert {
		var err error
		result, err = One.Div(result)
		if err != nil {
			return Zero, err
		}
	}

	// Apply sign
	if neg {
		result = result.Neg()
	}

	return result, nil
}
