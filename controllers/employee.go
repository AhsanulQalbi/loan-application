package controllers

import (
	"loan_application/models"
	"loan_application/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	EmployeeRepo *repositories.EmployeeRepository
}

func NewEmployeeController(repo *repositories.EmployeeRepository) *EmployeeController {
	return &EmployeeController{EmployeeRepo: repo}
}

func (e *EmployeeController) GetEmployees(ctx *gin.Context) {
	employees, err := e.EmployeeRepo.GetEmployees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting data"})
		return
	}
	ctx.JSON(http.StatusOK, employees)
}

func (e *EmployeeController) GetEmployeeByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
		return
	}

	employee, err := e.EmployeeRepo.GetEmployeeByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed GetEmployeeByID"})
		return
	}

	ctx.JSON(http.StatusOK, employee)

}

func (e *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var createEmployeeRequest models.Employee
	if err := ctx.ShouldBind(&createEmployeeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEmployeeID, err := e.EmployeeRepo.CreateEmployee(&createEmployeeRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed CreateEmployee"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Employee created successfully", "employee_id": newEmployeeID})
}
