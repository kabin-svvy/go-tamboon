package omise

import (
	"encoding/json"
	"fmt"
	"go-tamboon/csv"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var omiseChargeByTokenPath = ""
var omiseChargeByTokenUsername = ""

type ChargeByTokenResponse struct {
	Object         string `json:"object"`
	ID             string `json:"id"`
	Location       string `json:"location"`
	Amount         int    `json:"amount"`
	Net            int    `json:"net"`
	Fee            int    `json:"fee"`
	FeeVat         int    `json:"fee_vat"`
	Interest       int    `json:"interest"`
	InterestVat    int    `json:"interest_vat"`
	FundingAmount  int    `json:"funding_amount"`
	RefundedAmount int    `json:"refunded_amount"`
	Authorized     bool   `json:"authorized"`
	Capturable     bool   `json:"capturable"`
	Capture        bool   `json:"capture"`
	Disputable     bool   `json:"disputable"`
	Livemode       bool   `json:"livemode"`
	Refundable     bool   `json:"refundable"`
	Reversed       bool   `json:"reversed"`
	Reversible     bool   `json:"reversible"`
	Voided         bool   `json:"voided"`
	Paid           bool   `json:"paid"`
	Expired        bool   `json:"expired"`
	PlatformFee    struct {
		Fixed      string `json:"fixed"`
		Amount     string `json:"amount"`
		Percentage string `json:"percentage"`
	} `json:"platform_fee"`
	Currency        string `json:"currency"`
	FundingCurrency string `json:"funding_currency"`
	IP              string `json:"ip"`
	Refunds         struct {
		Object   string    `json:"object"`
		Data     []string  `json:"data"`
		Limit    int       `json:"limit"`
		Offset   int       `json:"offset"`
		Total    int       `json:"total"`
		Location string    `json:"location"`
		Order    string    `json:"order"`
		From     time.Time `json:"from"`
		To       time.Time `json:"to"`
	} `json:"refunds"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Metadata    struct {
	} `json:"metadata"`
	Card struct {
		Object            string    `json:"object"`
		ID                string    `json:"id"`
		Livemode          bool      `json:"livemode"`
		Location          string    `json:"location"`
		Deleted           bool      `json:"deleted"`
		Street1           string    `json:"street1"`
		Street2           string    `json:"street2"`
		City              string    `json:"city"`
		State             string    `json:"state"`
		PhoneNumber       string    `json:"phone_number"`
		PostalCode        string    `json:"postal_code"`
		Country           string    `json:"country"`
		Financing         string    `json:"financing"`
		Bank              string    `json:"bank"`
		Brand             string    `json:"brand"`
		Fingerprint       string    `json:"fingerprint"`
		FirstDigits       string    `json:"first_digits"`
		LastDigits        string    `json:"last_digits"`
		Name              string    `json:"name"`
		ExpirationMonth   int       `json:"expiration_month"`
		ExpirationYear    int       `json:"expiration_year"`
		SecurityCodeCheck bool      `json:"security_code_check"`
		CreatedAt         time.Time `json:"created_at"`
	} `json:"card"`
	Source                   string    `json:"source"`
	Schedule                 string    `json:"schedule"`
	Customer                 string    `json:"customer"`
	Dispute                  string    `json:"dispute"`
	Transaction              string    `json:"transaction"`
	FailureCode              string    `json:"failure_code"`
	FailureMessage           string    `json:"failure_message"`
	Status                   string    `json:"status"`
	AuthorizeURI             string    `json:"authorize_uri"`
	ReturnURI                string    `json:"return_uri"`
	CreatedAt                time.Time `json:"created_at"`
	PaidAt                   time.Time `json:"paid_at"`
	ExpiresAt                time.Time `json:"expires_at"`
	ExpiredAt                string    `json:"expired_at"`
	ReversedAt               string    `json:"reversed_at"`
	ZeroInterestInstallments bool      `json:"zero_interest_installments"`
	Branch                   string    `json:"branch"`
	Terminal                 string    `json:"terminal"`
	Device                   string    `json:"device"`
	Code                     string    `json:"code"`
	Message                  string    `json:"message"`
}

func (o Omiser) GetChargeByToken(p []csv.Transaction) {
	logrus.Info("GetChargeByToken")
	for k, v := range p {
		if v.ChargeToken == "" {
			continue
		}
		payload := url.Values{}
		payload.Add("description", fmt.Sprintf("%v charge by omise amount %v", v.CCNumber, v.AmountSubunits))
		payload.Add("amount", v.AmountSubunits)
		payload.Add("currency", v.Ccy)
		payload.Add("card", v.ChargeToken)

		req, err := http.NewRequest(http.MethodPost, omiseChargeByTokenPath, strings.NewReader(payload.Encode()))
		if err != nil {
			logrus.Error(err)
			continue
		}

		req.SetBasicAuth(omiseChargeByTokenUsername, "")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := o.httpClient.Do(req)
		if err != nil {
			logrus.Error(err)
			continue
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logrus.Error(err)
			continue
		}

		response := &ChargeByTokenResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			logrus.WithError(err).Info(string(body))
			continue
		}

		if response.ID != "" {
			p[k].ChargeID = response.ID
			p[k].ChargeStatus = response.Status
			logrus.Info(fmt.Sprintf("%v %v to charge %v %v with id %v", p[k].CCNumber, response.Status, p[k].AmountSubunits, p[k].Ccy, response.ID))
		} else {
			logrus.Error(fmt.Sprintf("%v fail to charge %v %v %v - %v", p[k].CCNumber, p[k].AmountSubunits, p[k].Ccy, response.Code, response.Message))
		}
	}
}
