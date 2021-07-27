package app

import (
	"crypto/tls"
	"fmt"
	"go-tamboon/config"
	"go-tamboon/csv"
	"go-tamboon/omise"
	"go-tamboon/summary"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func Run() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	err := config.Get()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get config")
	}

	logrus.Info("performing donations...")
	args := os.Args[1:]

	filename := ""

	for _, v := range args {
		filename = v
	}

	logrus.Info("Read file")
	transaction, err := csv.Read(filename)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read file")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	omiser := omise.New(client)
	omiser.GetChargeToken(transaction)
	omiser.GetChargeByToken(transaction)

	sum := summary.Get(transaction)

	logrus.Info("done.")

	fmt.Printf("	      total received: THB	%v\n", sum.TotalReceived)
	fmt.Printf("	successfully donated: THB	%v\n", sum.TotalDonated)
	fmt.Printf("	     faulty donation: THB	%v\n", sum.TotalFaultyDonation)
	fmt.Printf("\n")
	fmt.Printf("	  average per person: THB	%v\n", sum.AveragePerPerson)
	fmt.Printf("	          top donors: 		%v\n", sum.TopDonors[0].Name)
	fmt.Printf("					%v\n", sum.TopDonors[1].Name)
	fmt.Printf("					%v\n", sum.TopDonors[2].Name)
}
