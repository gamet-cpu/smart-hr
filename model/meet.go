package model

type Meet struct {
	RectuterID int    `json:"recruter_id"`
	UserID     int    `json:"user_id"`
	MeetDate   string `json:"meet_date"`
	StartTime  string `json:"start_time"`
	MeetLeight int    `json:"meet_leight"`
	MeetType   string `json:"meet_type"`
	Link       string `json:"link"`
	Additional string `json:"additional"`
}
