package controllers

import (
	"fmt"
	"loan_application/models"
	"loan_application/repositories"

	"path/filepath"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanApprovalController struct {
	LoanApprovalRepo *repositories.LoanApprovalRepository
	LoanRepo         *repositories.LoanRepository
	EmployeeRepo     *repositories.EmployeeRepository
}

func NewLoanApprovalController(repo *repositories.LoanApprovalRepository, loanRepo *repositories.LoanRepository, employeeRepo *repositories.EmployeeRepository) *LoanApprovalController {
	return &LoanApprovalController{
		LoanApprovalRepo: repo,
		LoanRepo:         loanRepo,
		EmployeeRepo:     employeeRepo,
	}
}

func (ctrl *LoanApprovalController) ApproveLoan(c *gin.Context) {
	var loanApproval models.LoanApproval
	if err := c.ShouldBind(&loanApproval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "err": err.Error()})
		return
	}

	approvalLoan, err := ctrl.LoanRepo.GetLoanByID(loanApproval.LoanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if approvalLoan == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "loan not found"})
		return
	}

	if approvalLoan.LoanState != "proposed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan state"})
		return
	}

	employee, err := ctrl.EmployeeRepo.GetEmployeeByID(loanApproval.EmployeeValidatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if employee == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "employee not found"})
		return
	}

	file, err := c.FormFile("visit_proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visit proof file is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PDF files are allowed"})
		return
	}

	savePath := fmt.Sprintf("uploads/visit_proof/%s", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	loanApproval.VisitProof = savePath
	err = ctrl.LoanApprovalRepo.ApproveLoan(loanApproval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan approved successfully", "visit_proof": savePath})
}
