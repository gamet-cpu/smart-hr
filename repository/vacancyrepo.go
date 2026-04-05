package repository

import (
	"context"
	"smart-hr/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VacancyRepository struct {
	db *pgxpool.Pool
}

func NewVacancyRepository(db *pgxpool.Pool) *VacancyRepository {
	return &VacancyRepository{db: db}
}
func (v *VacancyRepository) GetAllVacancy(ctx context.Context) ([]model.Vacancy, error) {
	rows, err := v.db.Query(ctx, `
		SELECT id, title, department_id, city_id,
		       is_remote, salary_min, salary_max,
		       currency, employment_type,
		       description, status,
		       published_at::varchar
		FROM vacancy
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []model.Vacancy

	for rows.Next() {
		var v model.Vacancy

		err := rows.Scan(
			&v.ID,
			&v.Title,
			&v.DepartmentID,
			&v.CityID,
			&v.IsRemote,
			&v.SalaryMin,
			&v.SalaryMax,
			&v.Currency,
			&v.EmploymentType,
			&v.Description,
			&v.Status,
			&v.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		vacancies = append(vacancies, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (v *VacancyRepository) GetVacancyByID(ctx context.Context, id int) (*model.Vacancy, error) {
	var vacancy model.Vacancy

	err := v.db.QueryRow(ctx, `
		SELECT id, title, department_id, city_id,
		       is_remote, salary_min, salary_max,
		       currency, employment_type,
		       description, status,
		       published_at::varchar
		FROM vacancy WHERE id=$1`, id).
		Scan(
			&vacancy.ID,
			&vacancy.Title,
			&vacancy.DepartmentID,
			&vacancy.CityID,
			&vacancy.IsRemote,
			&vacancy.SalaryMin,
			&vacancy.SalaryMax,
			&vacancy.Currency,
			&vacancy.EmploymentType,
			&vacancy.Description,
			&vacancy.Status,
			&vacancy.PublishedAt,
		)

	if err != nil {
		return nil, err
	}

	return &vacancy, nil
}

func (v *VacancyRepository) CreateVacancy(ctx context.Context) (err string) {

}
