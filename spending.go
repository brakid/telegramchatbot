package main

type Spending struct {
	date     int64
	name     string
	amount   float32
	currency string
}

type Spendings struct {
	spendings []*Spending
}

func New() *Spendings {
	return &Spendings{make([]*Spending, 0)}
}

func (s *Spendings) Add(spending *Spending) {
	s.spendings = append(s.spendings, spending)
}

func (s *Spendings) TotalAmount() float32 {
	var sum float32
	for _, spending := range s.spendings {
		sum += spending.amount
	}

	return sum
}
