package rate

import (
	"fmt"
	"github.com/AltScore/money/pkg/money"
	"github.com/AltScore/money/pkg/percent"
)

const (
	Daily    = 1
	Weekly   = 7
	BiWeekly = 14
	Monthly  = 30
	Yearly   = 360
	FullYear = 365
)

type Periodic struct {
	// Period The period of the rate expressed in days.
	Period uint `json:"period" bson:"period"`

	// Value The value of the rate.
	Value percent.Percent `json:"rate" bson:"rate"`
}

// NewPeriodicRate creates a new Periodic rate.
func NewPeriodicRate(period uint, value percent.Percent) Periodic {
	return Periodic{Period: period, Value: value}
}

// NewPeriodicRateFromInt creates a new Periodic rate from an int64.
func NewPeriodicRateFromInt(period uint, value int64) Periodic {
	return Periodic{Period: period, Value: percent.New(value)}
}

// NewPeriodicRateFromFloat64 creates a new Periodic rate from a float64.
func NewPeriodicRateFromFloat64(period uint, value float64) Periodic {
	return Periodic{Period: period, Value: percent.FromFloat64(value)}
}

// NominalToPeriod converts the rate to the given period applying linear conversión (Nominal rates).
func (r Periodic) NominalToPeriod(period uint) Periodic {
	return NewPeriodicRate(period, r.Value.ChangePeriodLinearlyFrom(r.Period, period))
}

// By applies the rate to the amount of money.
func (r Periodic) By(m money.Money) money.Money {
	return r.Value.By(m)
}

// RoundedBy applies the rate to the amount of money with half-even rounding.
func (r Periodic) RoundedBy(m money.Money) money.Money {
	return r.Value.RoundedBy(m)
}

// RoundedByWithPeriod applies the rate to the amount of money in the given period.
// It is equivalent to the following but with less rounding errors:
//
//	periodicRate.NominalToPeriod(period).By(amount)
func (r Periodic) RoundedByWithPeriod(amount money.Money, period uint) money.Money {
	return amount.Mul(int64(r.Value) * int64(period)).RoundedDiv(percent.ScaledPercentToRate * int64(r.Period))
}

func (r Periodic) String() string {
	return fmt.Sprintf("%s %s", r.Value, r.PeriodAsString())
}

func (r Periodic) PeriodAsString() string {
	switch r.Period {
	case Daily:
		return "daily"
	case Weekly:
		return "weekly"
	case BiWeekly:
		return "biweekly"
	case Monthly:
		return "monthly"
	case Yearly:
		return "yearly"
	case FullYear:
		return "full year"
	default:
		return fmt.Sprintf("%d days", r.Period)
	}
}

// Equal returns true if the other rate is equal to this rate.
// Two rates are equal if they have the same period and value.
func (r Periodic) Equal(other Periodic) bool {
	return r.Period == other.Period && r.Value == other.Value
}

// IsZero returns true if the rate is zero.
func (r Periodic) IsZero() bool {
	return r.Value.IsZero()
}

// IsNonZero returns true if the rate is not zero.
func (r Periodic) IsNonZero() bool {
	return r.Value.IsNonZero()
}
