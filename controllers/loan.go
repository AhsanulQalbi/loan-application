package controllers

import (
	"loan_application/models"
	"loan_application/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	LoanRepo     *repositories.LoanRepository
	BorrowerRepo *repositories.BorrowerRepository
}

func NewLoanController(repo *repositories.LoanRepository, borrowerRepo *repositories.BorrowerRepository) *LoanController {
	return &LoanController{
		LoanRepo:     repo,
		BorrowerRepo: borrowerRepo,
	}
}

func (c *LoanController) GetLoans(ctx *gin.Context) {
	employees, err := c.LoanRepo.GetLoans()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data", "err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (c *LoanController) GetLoanByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	loan, err := c.LoanRepo.GetLoanByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	ctx.JSON(http.StatusOK, loan)
}

func (l *LoanController) CreateLoan(c *gin.Context) {
	var CreateLoanRequest models.Loan
	if err := c.ShouldBind(&CreateLoanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrower, err := l.BorrowerRepo.GetBorrowerByID(CreateLoanRequest.BorrowerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if borrower == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "borrower not found"})
		return
	}

	newLoanID, err := l.LoanRepo.CreateLoan(&CreateLoanRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan created successfully", "loan_id": newLoanID})

}

func (l *LoanController) UpdateLoan(c *gin.Context) {
	var updateLoanRequest models.UpdateLoanRequest
	if err := c.ShouldBind(&updateLoanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loanID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	loan, err := l.LoanRepo.GetLoanByID(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error get loan data", "err": err.Error()})
		return
	}

	if loan == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan not found"})
		return
	}

	if loan.LoanState != "proposed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan state, already approved / proceed"})
		return
	}

	err = l.LoanRepo.UpdateLoan(loanID, updateLoanRequest.PrincipalAmount, updateLoanRequest.Rate, updateLoanRequest.ROI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan updated successfully"})
}
