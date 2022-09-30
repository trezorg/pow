package wordofwisdom

import "math/rand"

var quotations = []string{
	"The best way out is always through.",
	"Carpe Diem",
	"Always Do What You Are Afraid To Do",
	"Believe and act as if it were impossible to fail",
	"Keep steadily before you the fact that all true success depends at last upon yourself.",
}

func Quote() string {
	return quotations[rand.Intn(len(quotations))]
}
