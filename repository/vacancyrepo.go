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

func (v *VacancyRepository) CreateVacancy(ctx context.Context,
	title string,
	department_id int,
	city_id int,
	is_remote bool,
	salary_min int,
	salary_max int,
	currency string,
	employment_type string,
	description string,
	status string) error {
	_, err := v.db.Exec(ctx,
		`INSERT INTO vacancy (title,department_id, city_id,is_remote,salary_min,salary_max, currency,employment_type,description,status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		title, department_id, city_id, is_remote, salary_min, salary_max, currency, employment_type, description, status)

	return err
}

func (v *VacancyRepository) UpdateVacancy(
	ctx context.Context,
	id int,
	title string,
	department_id int,
	city_id int,
	is_remote bool,
	salary_min int,
	salary_max int,
	currency string,
	employment_type string,
	description string,
	status string,

) error {
	_, err := v.db.Exec(ctx, `
	UPDATE vacancy SET
		title = COALESCE(NULLIF($2,''), title),
		department_id = COALESCE(NULLIF($3,0), department_id),
		city_id = COALESCE(NULLIF($4,0), city_id),
		is_remote = COALESCE($5, is_remote),
		salary_min = COALESCE(NULLIF($6,0), salary_min),
		salary_max = COALESCE(NULLIF($7,0), salary_max),
		currency = COALESCE(NULLIF($8,''), currency),
		employment_type = COALESCE(NULLIF($9,'')::employment_enum, employment_type),
		description = COALESCE(NULLIF($10,''), description),
		status = COALESCE(NULLIF($11,'')::status_enum, status)
	WHERE id = $1
	`,
		id,
		title, department_id, city_id, is_remote, salary_min, salary_max, currency, employment_type, description, status,
	)

	return err
}
func (v *VacancyRepository) DeleteVacancy(ctx context.Context, id int) error {
	_, err := v.db.Exec(ctx,
		`DELETE FROM vacancy WHERE id=$1`, id)
	return err
}
