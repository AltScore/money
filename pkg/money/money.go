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

func (a Money) Equal(another Money) bool {
	return a.amount == another.amount && a.currency.Equals(another.currency)
}

func (a Money) TryEqual(another Money) (bool, error) {
	return a.Equal(another), nil
}

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

func (a Money) Sub(b Money) Money {
	if sub, err := a.TrySub(b); err != nil {
		panic(err)
	} else {
		return sub
	}
}

func (a Money) Mul(multiplier int64) Money {
	return Money{
		amount:   a.amount * multiplier,
		currency: a.currency,
	}
}

func (a Money) Div(divider int64) Money {
	return Money{
		amount:   a.amount / divider,
		currency: a.currency,
	}
}

func (a Money) RoundedDiv(divider int64) Money {
	return Money{
		amount:   utils.HalfEvenRounding(a.amount, divider),
		currency: a.currency,
	}
}

func (a Money) CurrencyCode() string {
	cur := a.currency
	if cur == nil {
		return ""
	}
	return cur.Code
}

// GetCurrencyCode required for Money to implement CommonTypeMoney
func (a Money) GetCurrencyCode() string {
	return a.CurrencyCode()
}

func (a Money) Cmp(b Money) int {
	if cmp, err := a.TryCmp(b); err != nil {
		panic(err)
	} else {
		return cmp
	}
}

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

func (a Money) String() string {
	return a.currency.Format(a.amount)

}

func (a Money) GoString() string {
	return fmt.Sprintf("money.FromFloat64(%v, %q)", a.Number(), a.CurrencyCode())
}

func (a Money) Amount() string {
	_, number := a.formatAsNumber()
	return number
}

func (a Money) IsZero() bool {
	return a.amount == 0
}

func (a Money) IsNegative() bool {
	return a.amount < 0
}

func (a Money) IsPositive() bool {
	return a.amount > 0
}

func (a Money) LessThan(amount Money) bool {

	calc := a.Cmp(amount) < 0
	return calc
}

// IsLessThanEqual is an alias for IsLessThanOrEqual
// Deprecated: Use IsLessThanOrEqual instead
func (a Money) IsLessThanEqual(amount Money) bool {
	return a.IsLessThanOrEqual(amount)
}

func (a Money) IsLessThanOrEqual(amount Money) bool {
	return a.Cmp(amount) <= 0
}

func (a Money) Number() float64 {
	if a.IsZero() {
		return 0
	}
	return float64(a.amount) / math.Pow10(a.currency.Fraction)
}

func (a Money) CheckSameCurrency(other Money) error {
	return a.assertSameCurrency(other)
}

func (a Money) IsGreaterThan(other Money) bool {
	return a.Cmp(other) > 0
}

func (a Money) IsGreaterThanOrEqual(other Money) bool {
	return a.Cmp(other) >= 0
}

func (a Money) Min(other Money) Money {
	if a.LessThan(other) {
		return a
	}
	return other
}

func (a Money) Negated() Money {
	return Money{
		amount:   -a.amount,
		currency: a.currency,
	}
}

func (a Money) Sign() int {
	if a.IsPositive() {
		return 1
	} else if a.IsZero() {
		return 0
	}
	return -1
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

	formatter := *c

	formatter.Grapheme = "" // Remove grapheme
	formatter.Thousand = "" // Remove thousand-separator
	amountStr := formatter.Format(m.amount)
	return c.Code, amountStr
}

func MustAdd(a, b Money) Money {
	return a.Add(b)
}

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
