package repository

import (
	"context"
	"smart-hr/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MeetRepository struct {
	db *pgxpool.Pool
}

func NewMeetRepository(db *pgxpool.Pool) *MeetRepository {
	return &MeetRepository{db: db}
}
func (m *MeetRepository) GetAllMeets(ctx context.Context, userId int, role string) ([]model.Meet, error) {
	var rows pgx.Rows
	var err error

	if role == "company" {
		rows, err = m.db.Query(ctx, `
		SELECT recruter_id, user_id, CAST(meet_date AS text), CAST(start_time AS text),
		       meet_leight, meet_type, link,
		       additional
		FROM meets WHERE recruter_id=$1
	`, userId)
	} else {
		rows, err = m.db.Query(ctx, `
		SELECT recruter_id, user_id, CAST(meet_date AS text), CAST(start_time AS text),
		       meet_leight, meet_type, link,
		       additional
		FROM meets WHERE user_id=$1
	`, userId)
	}
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

func (m *MeetRepository) GetMeetById(ctx context.Context, id int) (*model.Meet, error) {
	var meet model.Meet
	err := m.db.QueryRow(ctx, `
	SELECT recruter_id, user_id, CAST(meet_date AS text), CAST(start_time AS text),meet_leight, meet_type, link,additional FROM meets WHERE id=$1`, id).Scan(
		&meet.RectuterID,
		&meet.UserID,
		&meet.MeetDate,
		&meet.StartTime,
		&meet.MeetLeight,
		&meet.MeetType,
		&meet.Link,
		&meet.Additional,
	)
	if err != nil {
		return nil, err
	}

	return &meet, err
}

func (m *MeetRepository) CreateMeet(ctx context.Context, recruterId int, userId int, meetDate string, startTime string, meetLeight int, meetType string, link string, additional string) error {
	_, err := m.db.Exec(ctx,
		`INSERT INTO meets (recruter_id, user_id, meet_date, start_time, meet_leight, meet_type, link, additional)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, recruterId, userId, meetDate, startTime, meetLeight, meetType, link, additional)

	return err
}

func (m *MeetRepository) UpdateMeet(ctx context.Context, id int, meetDate string, startTime string, meetLeight int, meetType string, link string, additional string) error {
	_, err := m.db.Exec(ctx, `UPDATE meets SET
    meet_date = COALESCE(NULLIF($1, '')::date, meet_date),
    start_time = COALESCE(NULLIF($2, '')::time, start_time),
    meet_leight = COALESCE(NULLIF($3, 0), meet_leight),
    meet_type = COALESCE(NULLIF($4, ''), meet_type),
    link = COALESCE(NULLIF($5, ''), link),
    additional = COALESCE(NULLIF($6, ''), additional)
	WHERE id = $7;
	`, meetDate, startTime, meetLeight, meetType, link, additional, id)
	return err
}

func (m *MeetRepository) DeleteMeet(ctx context.Context, id int) error {
	_, err := m.db.Exec(ctx,
		`DELETE FROM meets WHERE id=$1`, id)
	return err
}
