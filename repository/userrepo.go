package repository

import (
	"context"
	"smart-hr/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// ================= CREATE =================

func (r *UserRepository) Create(ctx context.Context, name, email, password, role string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1,$2,$3,$4)`,
		name, email, password, role)

	return err
}

// ================= LOGIN =================

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (int, string, error) {
	var id int
	var hash string

	err := r.db.QueryRow(ctx,
		`SELECT id, password_hash FROM users WHERE email=$1`,
		email).Scan(&id, &hash)

	return id, hash, err
}

// ================= GET BY ID =================

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, `
		SELECT id, name, email, role,
		       company_name, phone, description,
		       created_at::varchar
		FROM users WHERE id=$1`, id).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CompanyName,
			&user.Phone,
			&user.Description,
			&user.CreatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= GET ALL =================

func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email, role,
		       company_name, phone, description,
		       created_at::varchar
		FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CompanyName,
			&user.Phone,
			&user.Description,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// ================= UPDATE =================

func (r *UserRepository) Update(ctx context.Context,
	id int,
	name, companyName, phone, description string,
) error {

	_, err := r.db.Exec(ctx,
		`UPDATE users
		 SET name=$1,
		     company_name=$2,
		     phone=$3,
		     description=$4
		 WHERE id=$5`,
		name, companyName, phone, description, id)

	return err
}

// ================= DELETE =================

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM users WHERE id=$1`, id)
	return err
}
