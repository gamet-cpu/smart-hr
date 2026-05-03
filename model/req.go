package model

type req struct {
	Title          string `json:"title"`
	DepartmentID   int    `json:"department_id"`
	CityID         int    `json:"city_id"`
	IsRemote       bool   `json:"is_remote"`
	SalaryMin      int    `json:"salary_min"`
	SalaryMax      int    `json:"salary_max"`
	Currency       string `json:"currency"`
	EmploymentType string `json:"employment_type"`
	Description    string `json:"description"`
	Status         string `json:"status"`
}
