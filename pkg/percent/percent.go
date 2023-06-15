package percent

import (
	"errors"
	"fmt"
	"math"

	"github.com/AltScore/money/pkg/money"
	"github.com/AltScore/money/pkg/parsers"
)

// Percent is a 3 decimal percent value. Internally it is stored as an int64 with a 3 digits scale
//
//	interestRate := percent.FromStr("3.5")
//
// It represents a 3.5% (a 0.03500) factor, and stored as 3500
type Percent int64

const (
	Decimals            = 4
	Scale               = 10000
	ScaledPercentToRate = 100 * Scale
	Zero                = Percent(0)
	OneHundred          = Percent(ScaledPercentToRate)

	// InterestRateNormalizingPeriod is the number of days in a period for the normalized interest rate (MIR).
	InterestRateNormalizingPeriod = 30
)

var ErrDivisionByZero = errors.New("division by zero")

// New returns a new Percent from the integer value. No decimals.
func New(intPct int64) Percent {
	return Percent(intPct * Scale)
}

// FromFloat64 returns a Percent from the float64 value. The values is the percent, 1.0 == 1%
func FromFloat64(pct float64) Percent {
	return Percent(math.Round(pct * Scale))
}

// FromFactor returns a Percent from the float64 value. It is assumed 1 == 100%
func FromFactor(pct float64) Percent {
	return Percent(math.Round(pct * ScaledPercentToRate))
}

// FromFraction64 returns a Percent the partial value represents on the total. If total is 0, returns Zero.
// Example: FromFraction64(3, 10) returns 30%
func FromFraction64(partial, total float64) Percent {
	if total == 0 {
		return Zero
	}
	return Percent(math.Round(partial * ScaledPercentToRate / total))
}

// FromFraction returns a Percent the partial amount represents on the total. If total is 0, returns Zero.
// Example: FromFraction(money.MustParse("3.00", "MXN"), money.MustParse("10, "MXN")) returns 30%
func FromFraction(partial, total money.Money) (Percent, error) {
	if total.IsZero() {
		return Zero, ErrDivisionByZero
	}

	if partial.IsZero() {
		return Zero, nil
	}

	err := partial.CheckSameCurrency(total)

	if err != nil {
		return Zero, err
	}

	return FromFraction64(partial.Number(), total.Number()), nil
}

func MustFromFraction(partial, total money.Money) Percent {
	if pct, err := FromFraction(partial, total); err != nil {
		panic(err)
	} else {
		return pct
	}
}

func Parse(pctStr string) (Percent, error) {
	pct, err := parsers.ParseNumber(pctStr, Decimals)
	return Percent(pct), err
}

func MustParse(pctStr string) Percent {
	if pct, err := Parse(pctStr); err != nil {
		panic(err)
	} else {
		return pct
	}
}

// By multiplies the given amount by this percent and returns the result amount.
// It does not round the result.
func (p Percent) By(amount money.Money) money.Money {
	return amount.Mul(int64(p)).Div(ScaledPercentToRate)
}

// RoundedBy multiplies the given amount by this percent and returns the result amount.
// It rounds the result using Money.RoundedDiv()
func (p Percent) RoundedBy(amount money.Money) money.Money {
	return amount.Mul(int64(p)).RoundedDiv(ScaledPercentToRate)
}

// ExtractPercentFromTotal returns the original base value of an amount witch already has been applied a percent.
// Example: 1000 * 0.3 + 1000 = 1300, ExtractPercentFromTotal(1300) returns 300
// It is equivalent to: 1300 / (1 + 0.3) * 0.3 = 300
func (p Percent) ExtractPercentFromTotal(amount money.Money) money.Money {
	if p.IsZero() {
		return money.Zero(amount.CurrencyCode())
	}
	return amount.Mul(int64(p)).Div(ScaledPercentToRate + int64(p))
}

// ExtractRoundedPercentFromTotal returns the original base value of an amount witch already has been applied a percent.
// Example: 1000 * 0.3 + 1000 = 1300, ExtractPercentFromTotal(1300) returns 300
// It is equivalent to: 1300 / (1 + 0.3) * 0.3 = 300 with additional rounding
func (p Percent) ExtractRoundedPercentFromTotal(amount money.Money) money.Money {
	if p.IsZero() {
		return money.Zero(amount.CurrencyCode())
	}
	return amount.Mul(int64(p)).RoundedDiv(ScaledPercentToRate + int64(p))
}

// func (p Percent) By(amount Money) (computed Money, remainder Money) {

func (p Percent) IsZero() bool {
	return p == 0
}

func (p Percent) Equal(other Percent) bool {
	return p == other
}

func (p Percent) String() string {
	number := parsers.FormatNumber(int64(p), Decimals)
	return removeDecimals(number)
}

func (p Percent) GoString() string {
	return fmt.Sprintf("percent.MustParse(%q)", p.String())
}

func (p Percent) IsNegative() bool {
	return p < 0
}

func (p Percent) Number() float64 {
	return float64(p) / float64(Scale)
}

func (p Percent) Factor() float64 {
	return float64(p) / float64(ScaledPercentToRate)
}

// ChangePeriod converts rate from one unit of time to another applying the compound interest.
// Used to convert Nominal Interested Rate to Effective Interest Rates.
//
// The problem:
//
//	We have the monthly (30 days) rate
//	We need the rate for 60 days (2 months)
//	We need the rate for 90 days (3 months)
//	We need the rate for 15 days (half months)
//	We need the rate for 45 days (one and a half months)
func (p Percent) ChangePeriod(periodSize, newPeriodSize int) Percent {
	if periodSize == newPeriodSize {
		return p
	}

	rate := math.Pow(1+p.Factor(), float64(newPeriodSize)/float64(periodSize)) - 1

	return FromFactor(rate)
}

// ChangePeriodLinearly converts a nominal interest rate to an interest rate for a given period.
// It is calculated as a linear interpolation between the nominal rate and the
// nominal rate for a period of InterestRateNormalizingPeriod.
func (p Percent) ChangePeriodLinearly(period uint) Percent {
	if period == InterestRateNormalizingPeriod {
		return p
	}

	rate := p.Factor() * float64(period) / InterestRateNormalizingPeriod

	return FromFactor(rate)
}

// GreaterThan returns true if this percent is greater than the other percent
func (p Percent) GreaterThan(percent Percent) bool {
	return p > percent
}

// LessThan returns true if this percent is less than the other percent
func (p Percent) LessThan(percent Percent) bool {
	return p < percent
}

func removeDecimals(number string) string {
	l := len(number) - 1
	for number[l] == '0' {
		l--
	}

	if number[l] == '.' {
		l--
	}

	return number[:l+1]
}
