package shipping

import "fmt"

type Money int64

func EUR(euros int64) Money {
	return Money(euros * 100)
}

func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	return fmt.Sprintf("%s%d.%02d", sign, m/100, m%100)
}
