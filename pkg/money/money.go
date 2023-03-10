package money

import (
	"fmt"
	"github.com/AltScore/money/pkg/utils"
	"math"
	"sync"

	"github.com/AltScore/money/pkg/parsers"
	m "github.com/Rhymond/go-money"
)

type Money m.Money

const USD = "USD"

var defaultCurrency = m.GetCurrency(USD)

func Zero(currencyCode string) Money {
	return NewFromInt(0, currencyCode)
}

func NewFromInt(amount int64, currencyCode string) Money {
	currency := m.GetCurrency(currencyCode)

	return fromEquivalentInt(amount*scales.Int(currency.Fraction), currencyCode)
}

func fromEquivalentInt(amount int64, currencyCode string) Money {
	return (Money)(*m.New(amount, currencyCode))
}

func FromFloat64(amount float64, currencyCode string) Money {
	if currencyCode == "" {
		panic(fmt.Errorf("currencyCode is empty"))
	}

	currency := m.GetCurrency(currencyCode)

	if currency == nil {
		panic(fmt.Sprintf("Currency %s not found", currencyCode))
	}

	amountInt := float2EquivalentInt(amount, currency)

	return fromEquivalentInt(amountInt, currencyCode)
}

func float2EquivalentInt(amount float64, currency *m.Currency) int64 {
	return int64(math.Round(amount * scales.Float(currency.Fraction)))
}

func Parse(amount string, currencyCode string) (Money, error) {
	currency := m.GetCurrency(currencyCode)

	amountInt, err := parsers.ParseNumber(amount, currency.Fraction)

	if err != nil {
		return Money{}, err
	}

	return fromEquivalentInt(amountInt, currencyCode), nil
}

func MustParse(amount string, currencyCode string) Money {
	currency := m.GetCurrency(currencyCode)

	amountInt, err := parsers.ParseNumber(amount, currency.Fraction)

	if err != nil {
		panic(m.ErrInvalidJSONUnmarshal)
	}

	return fromEquivalentInt(amountInt, currencyCode)
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
	am := a.asMoney()

	if am.IsZero() {
		return b, nil
	}

	bm := b.asMoney()

	if bm.IsZero() {
		return a, nil
	}

	add, err := am.Add(bm)
	if err != nil {
		return a, err
	}
	return (Money)(*add), err
}

func (a Money) Equal(another Money) bool {
	if equal, err := a.TryEqual(another); err != nil {
		panic(err)
	} else {
		return equal
	}
}

func (a Money) TryEqual(another Money) (bool, error) {
	return a.asMoney().Equals((*m.Money)(&another))
}

func (a Money) TrySub(b Money) (Money, error) {
	if b.IsZero() {
		return a, nil
	} else if a.IsZero() {
		return b.Negated(), nil
	}

	subtract, err := a.asMoney().Subtract(b.asMoney())
	return (Money)(*subtract), err
}

func (a Money) Sub(b Money) Money {
	if sub, err := a.TrySub(b); err != nil {
		panic(err)
	} else {
		return sub
	}
}

func (a Money) Mul(multiplier int64) Money {
	ma := a.asMoney()
	mul := ma.Multiply(multiplier)
	return (Money)(*mul)
}

func (a Money) Div(divider int64) Money {
	ma := a.asMoney()

	mul := ma.Amount() / divider
	return fromEquivalentInt(mul, a.CurrencyCode())
}

func (a Money) RoundedDiv(divider int64) Money {
	ma := a.asMoney()

	div := utils.HalfEvenRounding(ma.Amount(), divider)
	return fromEquivalentInt(div, a.CurrencyCode())
}

func (a Money) CurrencyCode() string {
	currency := a.asMoney().Currency()
	if currency == nil {
		return ""
	}
	return currency.Code
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

	am := a.asMoney()
	bm := b.asMoney()

	if isLess, err := am.LessThan(bm); err != nil {
		return 0, err
	} else if isLess {
		return -1, nil
	} else if isEqual, _ := am.Equals(bm); isEqual {
		return 0, nil
	} else {
		return 1, nil
	}
}

func (a Money) String() string {
	return a.asMoney().Display()
}

func (a Money) GoString() string {
	return fmt.Sprintf("money.FromFloat64(%v, %q)", a.Number(), a.CurrencyCode())
}

func (a Money) asMoney() *m.Money {
	return (*m.Money)(&a)
}

func (a Money) Amount() string {
	_, number := formatAsNumber(a.asMoney())
	return number
}

func (a Money) IsZero() bool {
	return a.asMoney().IsZero()
}

func (a Money) IsNegative() bool {
	return a.asMoney().IsNegative()
}

func (a Money) IsPositive() bool {
	return a.asMoney().IsPositive()
}

func (a Money) LessThan(amount Money) bool {

	calc := a.Cmp(amount) < 0
	return calc
}

func (a Money) IsLessThanEqual(line Money) bool {
	return a.Cmp(line) <= 0
}

func (a Money) Number() float64 {
	return float64(a.asMoney().Amount()) / math.Pow10(a.asMoney().Currency().Fraction)
}

func (a Money) CheckSameCurrency(total Money) error {
	_, err := a.TryCmp(total)
	return err
}

func (a Money) IsGreaterThan(zero Money) bool {
	return a.Cmp(zero) > 0
}

func (a Money) Min(other Money) Money {
	if a.LessThan(other) {
		return a
	}
	return other
}

func (a Money) Negated() Money {
	return fromEquivalentInt(-a.asMoney().Amount(), a.CurrencyCode())
}

func (a Money) Sign() int {
	if a.IsPositive() {
		return 1
	} else if a.IsZero() {
		return 0
	}
	return -1
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
