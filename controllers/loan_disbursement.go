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

type LoanDisbursementController struct {
	LoanDisbursementRepo *repositories.LoanDisbursementRepository
	LoanRepo             *repositories.LoanRepository
	EmployeeRepo         *repositories.EmployeeRepository
}

func NewLoanDisbursementController(repo *repositories.LoanDisbursementRepository, loanRepo *repositories.LoanRepository, employeeRepo *repositories.EmployeeRepository) *LoanDisbursementController {
	return &LoanDisbursementController{
		LoanDisbursementRepo: repo,
		LoanRepo:             loanRepo,
		EmployeeRepo:         employeeRepo,
	}
}

func (ctrl *LoanDisbursementController) DisburseLoan(c *gin.Context) {
	var disbursement models.LoanDisbursement
	if err := c.ShouldBind(&disbursement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check required fields", "err": err.Error()})
		return
	}

	loan, err := ctrl.LoanRepo.GetLoanByID(disbursement.LoanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if loan == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan not found"})
		return
	}

	if loan.LoanState != "invested" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan state"})
		return
	}

	employee, err := ctrl.EmployeeRepo.GetEmployeeByID(disbursement.EmployeeOfficerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if employee == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee not found"})
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

	savePath := fmt.Sprintf("uploads/agreement_letter_disbursement/%s", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	disbursement.AgreementLetter = savePath
	err = ctrl.LoanDisbursementRepo.DisburseLoan(disbursement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "loan disbursement success"})
}
