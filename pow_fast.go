package udecimal

import (
	"math/big"
)

const (
	maxExponential = uint32(0x80000)
)

func (d Decimal) PowFastInt32(e int32) (Decimal, error) {
	if d.coef.IsZero() {
		if e < 0 {
			return Decimal{}, ErrZeroPowNegative
		}
		if e == 0 {
			return One, nil
		}
		return Zero, nil
	}

	if e == 0 {
		return One, nil
	}

	if e == 1 {
		return d, nil
	}

	if d.coef.overflow() {
		return Decimal{}, errOverflow
	}

	dTrim := d.trimTrailingZeros()

	if e < 0 {
		return dTrim.powFastInverse(int(-e))
	}

	return dTrim.powFastPositive(int(e))
}

func (d Decimal) powFastPositive(e int) (Decimal, error) {
	if e >= int(maxExponential) {
		return Decimal{}, errOverflow
	}

	if d.coef.u128.hi != 0 && e >= 4 {
		return Decimal{}, errOverflow
	}

	exponent := int(d.prec) * e
	if exponent > int(defaultPrec)+38 {
		return Decimal{}, errOverflow
	}

	base := d.coef.u128
	result := u128{lo: 1}

	for e > 0 {
		if e&1 == 1 {
			product := result.MulToU256(base)
			if !product.carry.IsZero() {
				return Decimal{}, errOverflow
			}
			result = u128{hi: product.hi, lo: product.lo}
		}

		e >>= 1
		if e > 0 {
			product := base.MulToU256(base)
			if !product.carry.IsZero() {
				return Decimal{}, errOverflow
			}
			base = u128{hi: product.hi, lo: product.lo}
		}
	}

	neg := d.neg && e%2 == 1

	if exponent <= int(defaultPrec) {
		return newDecimal(neg, bintFromU128(result), uint8(exponent)), nil
	}

	factor := exponent - int(defaultPrec)
	if factor > 38 {
		return Decimal{}, errOverflow
	}

	result256 := u256{hi: result.hi, lo: result.lo}
	q, _, err := result256.fastQuo(pow10[factor])
	if err != nil {
		return Decimal{}, err
	}

	return newDecimal(neg, bintFromU128(q), defaultPrec), nil
}

func (d Decimal) powFastInverse(e int) (Decimal, error) {
	if e >= int(maxExponential) {
		return Decimal{}, errOverflow
	}

	absD := Decimal{coef: d.coef, neg: false, prec: d.prec}
	posResult, err := absD.powFastPositive(e)
	if err != nil {
		return Decimal{}, err
	}

	if posResult.coef.IsZero() {
		return Decimal{}, errOverflow
	}

	one := newDecimal(false, bintFromU128(u128{lo: 1}), 0)
	result, err := one.Div(posResult)
	if err != nil {
		return Decimal{}, err
	}

	if d.neg && e%2 == 1 {
		result.neg = true
	}

	return result, nil
}

func PowFastU128(base u128, exp int32) (u128, error) {
	if exp == 0 {
		return u128{lo: 1}, nil
	}

	if exp == 1 {
		return base, nil
	}

	if exp < 0 {
		return u128{}, errOverflow
	}

	if uint32(exp) >= maxExponential {
		return u128{}, errOverflow
	}

	result := u128{lo: 1}
	e := int(exp)

	for e > 0 {
		if e&1 == 1 {
			product := result.MulToU256(base)
			if !product.carry.IsZero() {
				return u128{}, errOverflow
			}
			result = u128{hi: product.hi, lo: product.lo}
		}

		e >>= 1
		if e > 0 {
			product := base.MulToU256(base)
			if !product.carry.IsZero() {
				return u128{}, errOverflow
			}
			base = u128{hi: product.hi, lo: product.lo}
		}
	}

	return result, nil
}

func PowFastBig(base *big.Int, exp int32) *big.Int {
	if exp == 0 {
		return big.NewInt(1)
	}

	if exp == 1 {
		return new(big.Int).Set(base)
	}

	if exp < 0 {
		return nil
	}

	result := big.NewInt(1)
	b := new(big.Int).Set(base)
	e := int(exp)

	for e > 0 {
		if e&1 == 1 {
			result.Mul(result, b)
		}
		e >>= 1
		if e > 0 {
			b.Mul(b, b)
		}
	}

	return result
}
