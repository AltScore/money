package money

const NanoDecimals = 9

// CommonTypeMoney allows to use a Google Common Type Money without creating a dependency on that package.
type CommonTypeMoney interface {
	GetCurrencyCode() string
	GetUnits() int64
	GetNanos() int32
}

func FromCommonType(cm CommonTypeMoney) Money {
	currency := getCurrencyWithDefault(cm.GetCurrencyCode())

	scale := scales.Int(currency.Fraction)
	nanosScale := scales.Int(NanoDecimals - currency.Fraction)

	return fromEquivalentInt(
		cm.GetUnits()*scale+int64(cm.GetNanos())/nanosScale,
		cm.GetCurrencyCode(),
	)
}

func (m Money) Decimals() int {
	if m.currency == nil {
		return 0
	}
	return m.currency.Fraction
}

func (m Money) AsUnitsAndNanos() (int64, int32) {
	decimals := m.Decimals()

	scale := scales.Int(decimals)

	intAmount := m.amount

	nanosScale := scales.Int(NanoDecimals - decimals)

	units := intAmount / scale
	nanos := int32((intAmount - units*scale) * nanosScale)

	return units, nanos
}

func (m Money) GetUnits() int64 {
	units, _ := m.AsUnitsAndNanos()
	return units
}

func (m Money) GetNanos() int32 {
	_, nanos := m.AsUnitsAndNanos()
	return nanos
}
