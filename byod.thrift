namespace go employer.calculator
/*Benefits for employees that have or are over
  * 2 - 5 - 10 Years
  * 2: 150USD
  * 5: 400USD
  * 10: 1200USD
 */
const double BenefitOverTwoYears = 150.00
const double BenefitOverFiveYears = 400.00
const double BenefitOverTenYears = 1200.00

/* 
 * Benefits for family size(excluding employee)
 * 2 members: 20USD
 * 3 members: 40USD
 * 4 or more: 60USD
*/ 
const double TwoFamilyMembersBenefit = 20.00
const double ThreeFamilyMembersBenefit = 40.00
const double FourFamilyMembersBenefit = 60.00

/* Benefits for position/role in the company
 * This is a static number that only covers three roles in particular
 * Security
 * Janitor
 * Apprentice
*/
const double PositionBenefit = 100.00

/* Discounts are applied to the Raw Income and so are the benefits. 
 * The total income corresponds to the raw income after applying all the benefits and discounts
*/
const double ExpenseDiscount = 0.03
const double EmployerDiscount = 0.01

struct Employee {
  1: required string FirstName
  2: required string LastName
  3: required double RawIncome
  5: required string Position 
  6: required i32 YearsInTheCompany
  7: required i32 FamilySize
}

struct Paycheck {
  1: required string EmployeeFullName
  2: required double TotalIncome
  3: required double TotalDiscounts
  4: required double TotalBenefits
}

service Employer {
  Paycheck CalculatePaycheck(1: Employee employee),
  string CalculateTotalDiscountsForEmployee(1: Employee employee),
  string CalculateTotalBenefitsForEmployee(1: Employee employee),
  string ReceivesPositionBenefit(1: Employee employee),
}