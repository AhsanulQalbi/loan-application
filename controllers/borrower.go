package controllers

import (
	"loan_application/models"
	"loan_application/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BorrowerController struct {
	BorrowerRepo *repositories.BorrowerRepository
}

func NewBorrowerController(repo *repositories.BorrowerRepository) *BorrowerController {
	return &BorrowerController{BorrowerRepo: repo}
}

func (e *BorrowerController) GetBorrowers(ctx *gin.Context) {
	borrowers, err := e.BorrowerRepo.GetBorrowers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting data"})
		return
	}

	ctx.JSON(http.StatusOK, borrowers)
}

func (e *BorrowerController) CreateBorrower(ctx *gin.Context) {
	var createBorrowerRequest models.Borrower
	if err := ctx.ShouldBind(&createBorrowerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBorrowerID, err := e.BorrowerRepo.CreateBorrower(&createBorrowerRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed CreateBorrower", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Borrower created successfully", "borrower_id": newBorrowerID})
}
