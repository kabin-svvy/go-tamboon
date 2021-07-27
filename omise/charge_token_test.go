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

type MockHttpClientForCreateTokenSuccess struct{}

func (m *MockHttpClientForCreateTokenSuccess) Do(req *http.Request) (*http.Response, error) {
	resBody, _ := json.Marshal(&ChargeTokenResponse{
		ID: "tokn_test",
	})

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, string(resBody))
	}))
	return http.Get(ts.URL)
}

func TestCreateTokenShouldBeSuccess(t *testing.T) {
	t.Run("test create token should be success", func(t *testing.T) {
		expected := "tokn_test"
		p1 := csv.Transaction{
			Name:           "Mr. Polo A Boffin",
			AmountSubunits: "4290655",
			CCNumber:       "4539519508429157",
			CVV:            "949",
			ExpMonth:       "6",
			ExpYear:        "2021",
		}

		p := []csv.Transaction{p1}

		client := &MockHttpClientForCreateTokenSuccess{}
		omiser := New(client)
		omiser.GetChargeToken(p)
		require.Equal(t, expected, p[0].ChargeToken)
	})
}
