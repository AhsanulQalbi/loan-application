package controllers

import (
	"fmt"
	"loan_application/models"
	"loan_application/repositories"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoanInvestmentController struct {
	LoanInvestmentRepo *repositories.LoanInvestmentRepository
	LoanRepo           *repositories.LoanRepository
	InvestorRepo       *repositories.InvestorRepository
}

func NewLoanInvestmentController(repo *repositories.LoanInvestmentRepository, loanRepo *repositories.LoanRepository, investorRepo *repositories.InvestorRepository) *LoanInvestmentController {
	return &LoanInvestmentController{
		LoanInvestmentRepo: repo,
		LoanRepo:           loanRepo,
		InvestorRepo:       investorRepo,
	}
}

func (ctrl *LoanInvestmentController) InvestInLoan(c *gin.Context) {
	var investment models.LoanInvestment
	if err := c.ShouldBind(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check required fields", "err": err.Error()})
		return
	}

	approvalLoan, err := ctrl.LoanRepo.GetLoanByID(investment.LoanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if approvalLoan == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan not found"})
		return
	}

	if approvalLoan.LoanState != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan state"})
		return
	}

	investor, err := ctrl.InvestorRepo.GetInvestorByID(investment.InvestorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if investor == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "investor not found"})
		return
	}

	file, err := c.FormFile("agreement_letter")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PDF files are allowed"})
		return
	}

	savePath := fmt.Sprintf("uploads/agreement_letter/%s", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	investment.AgreementLetter = savePath
	totalInvested, err := ctrl.LoanInvestmentRepo.InvestInLoan(investment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Investment added successfully, Current Total Invested : %.2f", totalInvested)})
}
