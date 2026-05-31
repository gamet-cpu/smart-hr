package repository

import (
	"context"
	"smart-hr/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MeetRepository struct {
	db *pgxpool.Pool
}

func NewMeetRepository(db *pgxpool.Pool) *MeetRepository {
	return &MeetRepository{db: db}
}
func (m *MeetRepository) GetAllMeets(ctx context.Context, userId int, role string) ([]model.Meet, error) {

	rows, err := m.db.Query(ctx, `
		SELECT recruter_id, user_id, meet_date, start_time,
		       meet_leight, meet_type, link,
		       additional
		FROM meets
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meets []model.Meet

	for rows.Next() {
		var m model.Meet

		err := rows.Scan(
			&m.RectuterID,
			&m.UserID,
			&m.MeetDate,
			&m.StartTime,
			&m.MeetLeight,
			&m.MeetType,
			&m.Link,
			&m.Additional,
		)
		if err != nil {
			return nil, err
		}

		meets = append(meets, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return meets, nil

}
