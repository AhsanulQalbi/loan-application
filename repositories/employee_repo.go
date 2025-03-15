package repositories

import (
	"database/sql"
	"loan_application/models"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) GetEmployees() ([]models.Employee, error) {
	rows, err := r.DB.Query("SELECT id, user_name, email, created_at FROM employees")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var employees []models.Employee
	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.ID, &employee.Username, &employee.Email, &employee.CreatedAt); err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetEmployeeByID(id int64) (*models.Employee, error) {
	var employee models.Employee
	query := "SELECT id, user_name, created_at FROM employees WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&employee.ID, &employee.Username, &employee.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeRepository) CreateEmployee(employee *models.Employee) (int64, error) {
	query := "INSERT INTO employees (user_name, email, created_at) VALUES ($1,$2,NOW()) RETURNING id"
	err := r.DB.QueryRow(query, employee.Username, employee.Email).Scan(&employee.ID)
	if err != nil {
		return 0, err
	}

	return employee.ID, nil
}
