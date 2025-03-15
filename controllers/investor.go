package controllers

import (
	"loan_application/models"
	"loan_application/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvestorController struct {
	InvestorRepo *repositories.InvestorRepository
}

func NewInvestorController(repo *repositories.InvestorRepository) *InvestorController {
	return &InvestorController{InvestorRepo: repo}
}

func (i *InvestorController) GetInvestors(ctx *gin.Context) {
	investors, err := i.InvestorRepo.GetInvestors()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting data"})
		return
	}
	ctx.JSON(http.StatusOK, investors)
}

func (i *InvestorController) CreateInvestor(ctx *gin.Context) {
	var createInvestorRequest models.Investor
	if err := ctx.ShouldBind(&createInvestorRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newInvestorID, err := i.InvestorRepo.CreateInvestors(&createInvestorRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed CreateInvestors", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Investor created successfully", "investor_id": newInvestorID})
}
