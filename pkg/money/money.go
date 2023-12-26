package money

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/AltScore/money/pkg/money/currency"
	"github.com/AltScore/money/pkg/utils"
	"go.uber.org/zap"

	"github.com/AltScore/money/pkg/parsers"
)

var (
	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")
	ErrorInvalidAmount      = errors.New("invalid amount")
	ErrorInvalidCurrency    = errors.New("invalid currency")
)

type Money struct {
	amount   int64
	currency *currency.Currency
}

func Zero(currencyCode string) Money {
	return NewFromInt(0, currencyCode)
}

func NewFromInt(amount int64, currencyCode string) Money {
	fraction := currency.GetOrDefault(currencyCode).Fraction
	return fromEquivalentInt(amount*scales.Int(fraction), currencyCode)
}

func fromEquivalentInt(amount int64, currencyCode string) Money {
	return Money{
		amount:   amount,
		currency: currency.GetOrDefault(currencyCode),
	}
}

func FromFloat64(amount float64, currencyCode string) Money {
	c := currency.GetOrDefault(currencyCode)

	amountInt := float2EquivalentInt(amount, c)

	return fromEquivalentInt(amountInt, currencyCode)
}

func float2EquivalentInt(amount float64, currency *currency.Currency) int64 {
	return int64(math.Round(amount * scales.Float(currency.Fraction)))
}

func Parse(amount string, currencyCode string) (Money, error) {
	fraction := currency.GetOrDefault(currencyCode).Fraction
	amountInt, err := parsers.ParseNumber(amount, fraction)

	if err != nil {
		return Money{}, err
	}

	return fromEquivalentInt(amountInt, currencyCode), nil
}

func MustParse(amount string, currencyCode string) Money {
	fraction := currency.GetOrDefault(currencyCode).Fraction
	amountInt, err := parsers.ParseNumber(amount, fraction)

	if err != nil {
		panic(ErrInvalidJSONUnmarshal)
	}

	return fromEquivalentInt(amountInt, currencyCode)
}

// SameCurrency check if given Money is equals by currency.
func (m Money) SameCurrency(om Money) bool {
	return m.currency.Equals(om.currency)
}

var ErrCurrencyMismatch = fmt.Errorf("currencies don't match")

func (m Money) assertSameCurrency(om Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

// Add sums the values including Zero
func (a Money) Add(b Money) Money {
	if add, err := a.TryAdd(b); err != nil {
		panic(err)
	} else {
		return add
	}
}

// TryAdd sums the values including Zero
// Returns error if currencies are not the same
func (a Money) TryAdd(b Money) (Money, error) {
	if a.IsZero() {
		if b.currency == nil && b.amount == 0 {
			// If zero is added to empty, return zero to preserve currency
			return a, nil
		}
		return b, nil
	}

	if b.IsZero() {
		return a, nil
	}

	if err := a.assertSameCurrency(b); err != nil {
		return a, err
	}

	return Money{
		amount:   a.amount + b.amount,
		currency: a.currency,
	}, nil
}

// Equal compares two Money values.
// Returns true if a == b and false otherwise
// If values are zero, and at most one currency is specified, returns true
func (a Money) Equal(another Money) bool {
	if a.currency == nil || another.currency == nil {
		// If one has no currency, check if amounts are 0. This is needed to compare empty values
		return a.amount == 0 && another.amount == 0
	}

	return a.amount == another.amount && a.currency.Equals(another.currency)
}

// IsEqual compares two Money values.
// Returns true if a == b and false otherwise
// If values are zero, and at most one currency is specified, returns true
// This is a synonym for Equals
func (a Money) IsEqual(another Money) bool {
	return a.Equal(another)
}

// IsNotEqual compares two Money values.
// Returns true if a != b and false otherwise
// If values are zero, and both currencies are specified, and they are different, returns true
func (a Money) IsNotEqual(another Money) bool {
	return !a.Equal(another)
}

// Same compares two Money values.
// Returns true if a == b and false otherwise.
// If values are zero, both currency should be equal or both nil.
func (a Money) Same(another Money) bool {
	return a.amount == another.amount && a.currency.Equals(another.currency)
}

// TryEqual compares two Money values.
// Returns true if a == b and false otherwise
// Returns error if currencies are not the same
// Deprecated: Use TryEquals instead
func (a Money) TryEqual(another Money) (bool, error) {
	return a.Equal(another), nil
}

// TryEquals compares two Money values.
// Returns true if a == b and false otherwise
// Returns error if currencies are not the same
func (a Money) TryEquals(another Money) (bool, error) {
	return a.Equal(another), nil
}

// TrySub subtracts the values including Zero
// Returns error if currencies are not the same
func (a Money) TrySub(b Money) (Money, error) {
	if b.IsZero() {
		return a, nil
	} else if a.IsZero() {
		return b.Negated(), nil
	}

	return Money{
		amount:   a.amount - b.amount,
		currency: a.currency,
	}, a.assertSameCurrency(b)
}

// Sub subtracts the values including Zero
func (a Money) Sub(b Money) Money {
	if sub, err := a.TrySub(b); err != nil {
		panic(err)
	} else {
		return sub
	}
}

// Mul multiplies money and returns result
func (a Money) Mul(multiplier int64) Money {
	return Money{
		amount:   a.amount * multiplier,
		currency: a.currency,
	}
}

// Div divides money and returns result without rounding
func (a Money) Div(divider int64) Money {
	return Money{
		amount:   a.amount / divider,
		currency: a.currency,
	}
}

// RoundedDiv divides money and rounds result using HalfEvenRounding
func (a Money) RoundedDiv(divider int64) Money {
	return Money{
		amount:   utils.HalfEvenRounding(a.amount, divider),
		currency: a.currency,
	}
}

// CurrencyCode returns currency code of the Money
func (a Money) CurrencyCode() string {
	cur := a.currency
	if cur == nil {
		return ""
	}
	return cur.Code
}

// GetCurrencyCode required for Money to implement CommonTypeMoney
func (a Money) GetCurrencyCode() string { return a.CurrencyCode() }

// Cmp compares two Money values.
// Returns -1 if a < b, 0 if a == b and 1 if a > b
// Panics if currencies are not the same
func (a Money) Cmp(b Money) int {
	if cmp, err := a.TryCmp(b); err != nil {
		panic(err)
	} else {
		return cmp
	}
}

// TryCmp compares two Money values.
// Returns -1 if a < b, 0 if a == b and 1 if a > b
// Returns error if currencies are not the same
func (a Money) TryCmp(b Money) (int, error) {
	if b.IsZero() {
		return a.Sign(), nil
	}
	if a.IsZero() {
		return -b.Sign(), nil
	}

	err := a.assertSameCurrency(b)

	if err != nil {
		return 0, err
	} else if a.amount < b.amount {
		return -1, nil
	} else if a.amount == b.amount {
		return 0, nil
	} else {
		return 1, nil
	}
}

// String implements fmt.Stringer
func (a Money) String() string {
	return a.currency.Format(a.amount)
}

// GoString implements fmt.GoStringer
func (a Money) GoString() string {
	return fmt.Sprintf("money.FromFloat64(%v, %q)", a.Number(), a.CurrencyCode())
}

// Amount returns the amount as a string
func (a Money) Amount() string {
	_, number := a.formatAsNumber()
	return number
}

// IsZero returns true if the amount is zero
func (a Money) IsZero() bool { return a.amount == 0 }

// IsEmpty returns true if the amount is zero and the currency is nil
func (a Money) IsEmpty() bool { return a.amount == 0 && a.currency == nil }

// IsNegative returns true if the amount is less than zero
func (a Money) IsNegative() bool { return a.amount < 0 }

// IsPositive returns true if the amount is greater than zero
func (a Money) IsPositive() bool { return a.amount > 0 }

// LessThan is an alias for IsLessThan
// Deprecated: Use IsLessThan instead
func (a Money) LessThan(amount Money) bool { return a.IsLessThan(amount) }

// IsLessThan is an alias for IsLessThan
func (a Money) IsLessThan(amount Money) bool { return a.Cmp(amount) < 0 }

// IsLessThanEqual is an alias for IsLessThanOrEqual
// Deprecated: Use IsLessThanOrEqual instead
func (a Money) IsLessThanEqual(amount Money) bool { return a.IsLessThanOrEqual(amount) }

// IsLessThanOrEqual returns true if the amount is less than or equal to the other amount
func (a Money) IsLessThanOrEqual(amount Money) bool { return a.Cmp(amount) <= 0 }

// Number returns the amount as a float64
func (a Money) Number() float64 {
	if a.IsZero() {
		return 0
	}
	return float64(a.amount) / math.Pow10(a.currency.Fraction)
}

// CheckSameCurrency returns an error if the other money is not the same currency
func (a Money) CheckSameCurrency(other Money) error { return a.assertSameCurrency(other) }

// IsGreaterThan returns true if the amount is greater than the other amount
func (a Money) IsGreaterThan(other Money) bool { return a.Cmp(other) > 0 }

// IsGreaterThanOrEqual returns true if the amount is greater than or equal to the other amount
func (a Money) IsGreaterThanOrEqual(other Money) bool { return a.Cmp(other) >= 0 }

// Min returns the smaller of two Money values.
func (a Money) Min(other Money) Money {
	if a.IsLessThan(other) {
		return a
	}
	return other
}

// Max returns the larger of two Money values.
func (a Money) Max(other Money) Money {
	if a.IsGreaterThan(other) {
		return a
	}
	return other
}

// StepToZero returns zero if the amount is negative, otherwise returns the amount.
// It correspond to the Step function.
func (a Money) StepToZero() Money {
	if a.IsNegative() {
		return Money{
			amount:   0,
			currency: a.currency,
		}
	}
	return a
}

// Negated returns the negated value of the money
func (a Money) Negated() Money {
	return Money{
		amount:   -a.amount,
		currency: a.currency,
	}
}

// Sign returns:
//
//	 1 if the amount is positive
//	 0 if the amount is zero
//	-1 if the amount is negative
func (a Money) Sign() int {
	if a.IsPositive() {
		return 1
	} else if a.IsZero() {
		return 0
	}
	return -1
}

// Zero returns the zero-ed value of the money
func (m Money) Zero() Money {
	return Money{
		amount:   0,
		currency: m.currency,
	}
}

func (m Money) formatAsNumber() (string, string) { // make
	c := m.currency

	if c == nil {
		if m.amount != 0 {
			amount := strconv.FormatInt(m.amount, 10)
			zap.L().Warn("Currency is nil, amount is " + amount)
		}
		c = currency.GetOrDefault("")
	}

	decimals := c.Fraction

	isNegative := m.amount < 0
	var absAmount int64
	var s string
	if isNegative {
		absAmount = -m.amount
	} else {
		absAmount = m.amount
	}
	s = strconv.FormatInt(absAmount, 10)

	if len(s) <= decimals {
		s = "0.0000000000000000"[0:decimals-len(s)+2] + s // Add leading zeros
	} else if decimals > 0 {
		s = s[:len(s)-decimals] + "." + s[len(s)-decimals:]
	}

	if isNegative {
		return c.Code, "-" + s
	}
	return c.Code, s
}

// MustAdd panics if the two currencies are not the same currency
func MustAdd(a, b Money) Money { return a.Add(b) }

type scale struct {
	Int   int64
	Float float64
}

type scaleMap struct {
	scales map[int]scale
	lock   sync.RWMutex
}

var scales = scaleMap{
	scales: make(map[int]scale),
	lock:   sync.RWMutex{},
}

func (s *scaleMap) GetScale(fraction int) scale {
	s.lock.RLock()
	if value, ok := s.scales[fraction]; ok {
		s.lock.RUnlock()
		return value
	}

	s.lock.RUnlock()
	s.lock.Lock()
	defer s.lock.Unlock()

	pow10 := math.Pow10(fraction)
	newScale := scale{
		Int:   int64(pow10),
		Float: pow10,
	}

	s.scales[fraction] = newScale

	return newScale
}

func (s *scaleMap) Int(decimals int) int64 {
	return s.GetScale(decimals).Int
}

func (s *scaleMap) Float(decimals int) float64 {
	return s.GetScale(decimals).Float
}
