package percent

// FromFloat64Ptr returns a Percent pointer from a float64 pointer
// If the float64 pointer is nil, it returns nil
func FromFloat64Ptr(f *float64) *Percent {
	if f == nil {
		return nil
	}
	p := FromFloat64(*f)
	return &p
}
