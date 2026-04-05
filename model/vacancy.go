package model

type Vacancy struct {
	ID             int
	Title          string
	DepartmentID   int
	CityID         int
	IsRemote       bool
	SalaryMin      int
	SalaryMax      int
	Currency       string
	EmploymentType string
	Description    string
	Status         string
	PublishedAt    string
}
