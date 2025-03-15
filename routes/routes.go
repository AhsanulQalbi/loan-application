package routes

import (
	"database/sql"
	"loan_application/controllers"
	"loan_application/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	loanStateRepo := repositories.NewLoanStateRepository(db)

	borrowerRepo := repositories.NewBorrowerRepository(db)
	borrowerController := controllers.NewBorrowerController(borrowerRepo)

	loanRepo := repositories.NewLoanRepository(db, loanStateRepo)
	loanController := controllers.NewLoanController(loanRepo, borrowerRepo)

	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeController := controllers.NewEmployeeController(employeeRepo)

	investorRepo := repositories.NewInvestorRepository(db)
	investorController := controllers.NewInvestorController(investorRepo)

	loanApprovalRepo := repositories.NewLoanApprovalRepository(db, loanStateRepo)
	loanApprovalController := controllers.NewLoanApprovalController(loanApprovalRepo, loanRepo, employeeRepo)

	loanInvestmentRepo := repositories.NewLoanInvestmentRepository(db, loanStateRepo, investorRepo)
	loanInvestmentController := controllers.NewLoanInvestmentController(loanInvestmentRepo, loanRepo, investorRepo)

	loanDisbursementRepo := repositories.NewLoanDisbirsementRepository(db, loanStateRepo)
	loanDisbursementController := controllers.NewLoanDisbursementController(loanDisbursementRepo, loanRepo, employeeRepo)

	loanStateController := controllers.NewLoanStateController(loanStateRepo)

	r.GET("/loans", loanController.GetLoans)
	r.POST("/loans", loanController.CreateLoan)
	r.PUT("/loans/:id", loanController.UpdateLoan)

	r.POST("/loan-approval", loanApprovalController.ApproveLoan)
	r.POST("/loan-invest", loanInvestmentController.InvestInLoan)
	r.POST("/loan-disburse", loanDisbursementController.DisburseLoan)

	r.POST("/employees", employeeController.CreateEmployee)
	r.GET("/employees", employeeController.GetEmployees)

	r.POST("/investors", investorController.CreateInvestor)
	r.GET("/investors", investorController.GetInvestors)

	r.POST("/borrowers", borrowerController.CreateBorrower)
	r.GET("/borrowers", borrowerController.GetBorrowers)

	r.GET("/loan-states", loanStateController.GetLoanStates)

	return r
}
