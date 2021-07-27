package summary

import (
	"go-tamboon/csv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSummaryShouldBeSuccess(t *testing.T) {
	p1 := csv.Transaction{
		Name:           "A",
		AmountSubunits: "10",
		CCNumber:       "1",
		CVV:            "111",
		ExpMonth:       "1",
		ExpYear:        "2021",
		ChargeID:       "chrg_test",
		Ccy:            "THB",
		ChargeToken:    "tokn_test",
		ChargeStatus:   "successful",
	}
	p2 := csv.Transaction{
		Name:           "B",
		AmountSubunits: "20",
		CCNumber:       "2",
		CVV:            "222",
		ExpMonth:       "1",
		ExpYear:        "2021",
		ChargeID:       "",
		Ccy:            "THB",
		ChargeToken:    "",
		ChargeStatus:   "",
	}
	p3 := csv.Transaction{
		Name:           "C",
		AmountSubunits: "30",
		CCNumber:       "3",
		CVV:            "333",
		ExpMonth:       "1",
		ExpYear:        "2021",
		ChargeID:       "chrg_test",
		Ccy:            "THB",
		ChargeToken:    "tokn_test",
		ChargeStatus:   "successful",
	}
	p := []csv.Transaction{p1, p2, p3}
	t.Run("test get summary total received should be 60", func(t *testing.T) {
		expected := "60.00"
		sum := Get(p)
		require.Equal(t, expected, sum.TotalReceived)
	})
	t.Run("test get summary total donated should be 40", func(t *testing.T) {
		expected := "40.00"
		sum := Get(p)
		require.Equal(t, expected, sum.TotalDonated)
	})
	t.Run("test get summary total faulty should be 20", func(t *testing.T) {
		expected := "20.00"
		sum := Get(p)
		require.Equal(t, expected, sum.TotalFaultyDonation)
	})
	t.Run("test get summary average should be 20", func(t *testing.T) {
		expected := "20.00"
		sum := Get(p)
		require.Equal(t, expected, sum.AveragePerPerson)
	})
	t.Run("test get summary top 1 should be C", func(t *testing.T) {
		expected := "C"
		sum := Get(p)
		require.Equal(t, expected, sum.TopDonors[0].Name)
	})
}
