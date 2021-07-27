package summary

import (
	"fmt"
	"go-tamboon/csv"
	"strconv"
)

type Summary struct {
	TotalReceived       string
	TotalDonated        string
	TotalFaultyDonation string
	AveragePerPerson    string
	TopDonors           []Ranking
}

type Ranking struct {
	Name   string
	Amount float64
}

func Get(p []csv.Transaction) Summary {
	totalReceived := 0.0
	totalDonated := 0.0
	totalFaultyDonation := 0.0
	averagePerPerson := 0.0
	r := Ranking{
		Name:   "",
		Amount: 0.0,
	}
	rank := []Ranking{r, r, r}

	cntSuccess := 0

	for _, v := range p {
		amount, err := strconv.ParseFloat(v.AmountSubunits, 64)
		if err != nil {
			continue
		}
		totalReceived = totalReceived + amount
		if v.ChargeID != "" {
			if amount > rank[0].Amount {
				rank[0].Name = v.Name
				rank[0].Amount = amount
			} else if amount > rank[1].Amount {
				rank[1].Name = v.Name
				rank[1].Amount = amount
			} else if amount > rank[2].Amount {
				rank[2].Name = v.Name
				rank[2].Amount = amount
			}
			totalDonated = totalDonated + amount
			cntSuccess++
		} else {
			totalFaultyDonation = totalFaultyDonation + amount
		}
	}

	if cntSuccess > 0 {
		averagePerPerson = totalDonated / float64(cntSuccess)
	}

	return Summary{
		TotalReceived:       fmt.Sprintf("%.2f", totalReceived),
		TotalDonated:        fmt.Sprintf("%.2f", totalDonated),
		TotalFaultyDonation: fmt.Sprintf("%.2f", totalFaultyDonation),
		AveragePerPerson:    fmt.Sprintf("%.2f", averagePerPerson),
		TopDonors:           rank,
	}
}
