package controllers

import (
	"loan_application/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanStateController struct {
	LoanStateRepo *repositories.LoanStateRepository
}

func NewLoanStateController(repo *repositories.LoanStateRepository) *LoanStateController {
	return &LoanStateController{LoanStateRepo: repo}
}

func (e *LoanStateController) GetLoanStates(ctx *gin.Context) {
	loanStates, err := e.LoanStateRepo.GetLoanStates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting data"})
		return
	}

	ctx.JSON(http.StatusOK, loanStates)
}
