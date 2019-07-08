package main

import (
	"context"
	"fmt"
	"os"

	"github.com/amartine59/thpractice/gen-go/employer/calculator"
	"github.com/apache/thrift/lib/go/thrift"
)

const (
	// NetworkAddress ...
	NetworkAddress = "127.0.0.1:9090"
)

func main() {
	ctx := context.Background()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTJSONProtocolFactory()

	transport, err := thrift.NewTSocket(NetworkAddress)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address: ", err)
		os.Exit(1)
	}

	useTransport, err := transportFactory.GetTransport(transport)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error getting transport: ", err)
		os.Exit(1)
	}

	client := calculator.NewEmployerClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", NetworkAddress, " ", err)
		os.Exit(1)
	}

	defer transport.Close()
	genericEmployee := &calculator.Employee{
		FirstName:         "Peter",
		LastName:          "Harbor",
		FamilySize:        2,
		Position:          "Apprentice",
		YearsInTheCompany: 2,
		RawIncome:         900.00,
	}

	resultantPaycheck, err := client.CalculatePaycheck(ctx, genericEmployee)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while calculating the paycheck", err)
		os.Exit(1)
	}
	printPaycheckToTerminal(*resultantPaycheck)

	totalBenefits, err := client.CalculateTotalBenefitsForEmployee(ctx, genericEmployee)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while calculating the benefits", err)
		os.Exit(1)
	}

	fmt.Println(totalBenefits)

	totalDiscounts, err := client.CalculateTotalDiscountsForEmployee(ctx, genericEmployee)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while calculating the discounts", err)
		os.Exit(1)
	}

	fmt.Println(totalDiscounts)

	positionBenefit, err := client.ReceivesPositionBenefit(ctx, genericEmployee)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while calculating the position benefit", err)
		os.Exit(1)
	}

	fmt.Println(positionBenefit)
}

func printPaycheckToTerminal(paycheck calculator.Paycheck) {
	fmt.Printf("Paycheck\n Employee Full Name: %v\n TotalBenefits: %.2f\n TotalDiscounts: %.2f\n TotalIncome: %.2f\n", paycheck.EmployeeFullName, paycheck.TotalBenefits, paycheck.TotalDiscounts, paycheck.TotalIncome)
}
