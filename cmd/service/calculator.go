package main

import "time"

// Calculator скажет как быстро расло число (в сек)
// по сравнению с предыдущим измерением.
type Calculator struct {
	t time.Time
	n int

	prevRate float64
}

func NewCalculator() *Calculator {
	return &Calculator{
		t: time.Now(),
	}
}

func (c *Calculator) Set(n int) (float64, bool) {
	dt := time.Since(c.t)
	dn := n - c.n

	rate := float64(dn) / dt.Seconds()
	grown := rate > c.prevRate

	c.t = time.Now()
	c.n = n
	c.prevRate = rate

	return rate, grown
}
