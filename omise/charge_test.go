package omise

import (
	"encoding/json"
	"fmt"
	"go-tamboon/csv"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockHttpClientForChargeByTokenSuccess struct{}

func (m *MockHttpClientForChargeByTokenSuccess) Do(req *http.Request) (*http.Response, error) {
	resBody, _ := json.Marshal(&ChargeByTokenResponse{
		ID:     "chrg_test",
		Status: "successful",
	})

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, string(resBody))
	}))
	return http.Get(ts.URL)
}

func TestChargeByTokenShouldBeSuccess(t *testing.T) {
	t.Run("test charge by token should be success", func(t *testing.T) {
		expected := "chrg_test"
		p1 := csv.Transaction{
			Name:           "Mr. Polo A Boffin",
			AmountSubunits: "4290655",
			CCNumber:       "4539519508429157",
			CVV:            "949",
			ExpMonth:       "6",
			ExpYear:        "2021",
			Ccy:            "THB",
			ChargeToken:    "tokn_test",
		}

		p := []csv.Transaction{p1}

		client := &MockHttpClientForChargeByTokenSuccess{}
		omiser := New(client)
		omiser.GetChargeByToken(p)
		require.Equal(t, expected, p[0].ChargeID)
	})
}
