package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/sirupsen/logrus"

	"github.com/amartine59/thpractice/gen-go/employer/calculator"
)

const (
	// NetworkAddress ...
	NetworkAddress = "127.0.0.1:9090"
)

type Employer struct{}

// CalculatePaycheck calculates a paycheck for a given employee
func (em *Employer) CalculatePaycheck(ctx context.Context, employee *calculator.Employee) (paycheck *calculator.Paycheck, err error) {
	paycheck.EmployeeFullName = fmt.Sprintf("%v %v", employee.FirstName, employee.LastName)
	paycheck.TotalBenefits = CalculateTotalBenefits(*employee)
	paycheck.TotalDiscounts = CalculateDiscounts(employee.RawIncome)
	paycheck.TotalIncome = employee.RawIncome + (paycheck.TotalBenefits - paycheck.TotalDiscounts)

	logrus.Infof("Paycheck for %v\n TotalBenefits: %.2f\n TotalDiscounts: %.2f\n TotalIncome: %.2f\n", paycheck.EmployeeFullName, paycheck.TotalBenefits, paycheck.TotalDiscounts, paycheck.TotalIncome)
	return
}

// CalculateTotalDiscountsForEmployee calculates a totalDiscounts for a given employee
func (em *Employer) CalculateTotalDiscountsForEmployee(ctx context.Context, employee *calculator.Employee) (discounts string, err error) {
	totalDiscounts := CalculateDiscounts(employee.RawIncome)
	discounts = fmt.Sprintf("Total discounts for %v %v :\t RawValue: %v\t WithTwoDecimalPlaces: %.2f", employee.FirstName, employee.LastName, totalDiscounts, totalDiscounts)
	return
}

// CalculateTotalBenefitsForEmployee calculates a totalBenefits for a given employee
func (em *Employer) CalculateTotalBenefitsForEmployee(ctx context.Context, employee *calculator.Employee) (benefits string, err error) {
	totalBenefits := CalculateTotalBenefits(*employee)
	benefits = fmt.Sprintf("Total benefits for %v %v :\t RawValue: %v\t WithTwoDecimalPlaces: %.2f", employee.FirstName, employee.LastName, totalBenefits, totalBenefits)
	return
}

// ReceivesPositionBenefit prints a message that says whether an employee receives or not a benefit from his role in the company
func (em *Employer) ReceivesPositionBenefit(ctx context.Context, employee *calculator.Employee) (positionBenefit string, err error) {
	if pbenefit := calculatePositionBenefit(employee.Position); pbenefit == 0.0 {
		positionBenefit = fmt.Sprintf("Employee %v %v does not receive benefit for his role in the company", employee.FirstName, employee.LastName)
		return
	}

	positionBenefit = fmt.Sprintf("Employee %v %v does receive benefit for his role in the company", employee.FirstName, employee.LastName)
	return
}

// CalculateTotalBenefits calculates the total benefits for an employee
func CalculateTotalBenefits(employee calculator.Employee) float64 {
	return calculateFamilyBenefit(employee.FamilySize) + calculatePositionBenefit(employee.Position) + calculateYearsBenefit(employee.YearsInTheCompany)
}

func calculateFamilyBenefit(familySize int32) float64 {
	if familySize == 0 {
		return 0.0
	}

	switch familySize {
	case 2:
		return float64(calculator.TwoFamilyMembersBenefit)
	case 3:
		return float64(calculator.ThreeFamilyMembersBenefit)
	default:
		return float64(calculator.FourFamilyMembersBenefit)
	}
}

func calculatePositionBenefit(position string) float64 {
	switch position {
	case "Janitor", "Security", "Apprentice":
		return float64(calculator.PositionBenefit)
	default:
		return 0.0
	}
}

func calculateYearsBenefit(yearsInTheCompany int32) float64 {
	if yearsInTheCompany == 0 {
		return 0.0
	}

	if yearsInTheCompany >= 2 && yearsInTheCompany < 5 {
		return float64(calculator.BenefitOverTwoYears)
	}

	if yearsInTheCompany >= 5 && yearsInTheCompany < 10 {
		return float64(calculator.BenefitOverFiveYears)
	}

	return float64(calculator.BenefitOverTenYears)
}

// CalculateDiscounts calculates the total discount for an employee
func CalculateDiscounts(rawIncome float64) float64 {
	return calculateEmployerDiscount(rawIncome) + calculateExpenseDiscount(rawIncome)
}

func calculateEmployerDiscount(rawIncome float64) float64 {
	return rawIncome * calculator.EmployerDiscount
}

func calculateExpenseDiscount(rawIncome float64) float64 {
	return rawIncome * calculator.ExpenseDiscount
}

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTJSONProtocolFactory()
	serverTransport, err := thrift.NewTServerSocket(NetworkAddress)
	if err != nil {
		fmt.Println("Error! ", err)
		os.Exit(1)
	}

	handler := &Employer{}
	processor := calculator.NewEmployerProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("thrift server in", NetworkAddress)
	server.Serve()
}
