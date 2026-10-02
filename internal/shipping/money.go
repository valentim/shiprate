package shipping

import "fmt"

// Money is an amount in euro cents.
type Money int64

// EUR returns the given number of euros as Money.
func EUR(euros int64) Money {
	return Money(euros * 100)
}

// String formats the amount with two decimals, for example "18.00".
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	return fmt.Sprintf("%s%d.%02d", sign, m/100, m%100)
}
