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
	"github.com/spf13/viper"
)

var omiseChargeTokenPath = ""
var omiseChargeTokenUsername = ""

type HttpClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

type Omiser struct {
	httpClient HttpClienter
}

type ChargeTokenResponse struct {
	Object       string `json:"object"`
	ID           string `json:"id"`
	Livemode     bool   `json:"livemode"`
	Location     string `json:"location"`
	Used         bool   `json:"used"`
	ChargeStatus string `json:"charge_status"`
	Card         struct {
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
	CreatedAt time.Time `json:"created_at"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
}

func New(httpClient HttpClienter) Omiser {
	omiseChargeTokenPath = viper.GetString("omise.vault_url")
	omiseChargeTokenUsername = viper.GetString("omise.vault_key")
	omiseChargeByTokenPath = viper.GetString("omise.api_url")
	omiseChargeByTokenUsername = viper.GetString("omise.api_key")
	return Omiser{
		httpClient: httpClient,
	}
}

func (o Omiser) GetChargeToken(p []csv.Transaction) {
	logrus.Info("GetChargeToken")
	for k, v := range p {
		payload := url.Values{}
		payload.Add("card[name]", v.Name)
		payload.Add("card[number]", v.CCNumber)
		payload.Add("card[security_code]", v.CVV)
		payload.Add("card[expiration_month]", v.ExpMonth)
		payload.Add("card[expiration_year]", v.ExpYear)

		req, err := http.NewRequest(http.MethodPost, omiseChargeTokenPath, strings.NewReader(payload.Encode()))
		if err != nil {
			logrus.Error(err)
			continue
		}

		req.SetBasicAuth(omiseChargeTokenUsername, "")
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

		response := &ChargeTokenResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			logrus.WithError(err).Info(string(body))
			continue
		}

		if response.ID != "" {
			p[k].ChargeToken = response.ID
			logrus.Info(fmt.Sprintf("%v success to get token %v", p[k].CCNumber, response.ID))
		} else {
			logrus.Error(fmt.Sprintf("%v fail to get token %v - %v", p[k].CCNumber, response.Code, response.Message))
		}
	}
}
