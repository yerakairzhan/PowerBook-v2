// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: readings.sql

package db

import (
	"context"
	"time"
)

const createReadingLog = `-- name: CreateReadingLog :exec
INSERT INTO reading_logs (userid, username, date, minutes_read)
VALUES (
           $1,$2, $3, $4
       )
`

type CreateReadingLogParams struct {
	Userid      string    `json:"userid"`
	Username    string    `json:"username"`
	Date        time.Time `json:"date"`
	MinutesRead int32     `json:"minutes_read"`
}

func (q *Queries) CreateReadingLog(ctx context.Context, arg CreateReadingLogParams) error {
	_, err := q.db.ExecContext(ctx, createReadingLog,
		arg.Userid,
		arg.Username,
		arg.Date,
		arg.MinutesRead,
	)
	return err
}

const getReadingLeaderboard = `-- name: GetReadingLeaderboard :many
SELECT
    r.userid,
    u.username,
    SUM(r.minutes_read) AS total_minutes,
    COUNT(DISTINCT CASE WHEN r.minutes_read > 30 THEN r.date END) AS days_read_more_than_30
FROM reading_logs r
         JOIN users u ON r.userid = u.userid
WHERE r.date >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY r.userid, u.username
ORDER BY days_read_more_than_30 DESC, total_minutes DESC
    LIMIT 5
`

type GetReadingLeaderboardRow struct {
	Userid             string `json:"userid"`
	Username           string `json:"username"`
	TotalMinutes       int64  `json:"total_minutes"`
	DaysReadMoreThan30 int64  `json:"days_read_more_than_30"`
}

func (q *Queries) GetReadingLeaderboard(ctx context.Context) ([]GetReadingLeaderboardRow, error) {
	rows, err := q.db.QueryContext(ctx, getReadingLeaderboard)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReadingLeaderboardRow
	for rows.Next() {
		var i GetReadingLeaderboardRow
		if err := rows.Scan(
			&i.Userid,
			&i.Username,
			&i.TotalMinutes,
			&i.DaysReadMoreThan30,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReadingLogsByUser = `-- name: GetReadingLogsByUser :many
SELECT date, minutes_read
FROM reading_logs
WHERE userid = $1
ORDER BY date DESC
`

type GetReadingLogsByUserRow struct {
	Date        time.Time `json:"date"`
	MinutesRead int32     `json:"minutes_read"`
}

func (q *Queries) GetReadingLogsByUser(ctx context.Context, userid string) ([]GetReadingLogsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getReadingLogsByUser, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReadingLogsByUserRow
	for rows.Next() {
		var i GetReadingLogsByUserRow
		if err := rows.Scan(&i.Date, &i.MinutesRead); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSumReading = `-- name: GetSumReading :one
select sum(minutes_read) as Sum, username, userid from reading_logs where userid = $1 group by userid, username
`

type GetSumReadingRow struct {
	Sum      int64  `json:"sum"`
	Username string `json:"username"`
	Userid   string `json:"userid"`
}

func (q *Queries) GetSumReading(ctx context.Context, userid string) (GetSumReadingRow, error) {
	row := q.db.QueryRowContext(ctx, getSumReading, userid)
	var i GetSumReadingRow
	err := row.Scan(&i.Sum, &i.Username, &i.Userid)
	return i, err
}

const updateReadingLog = `-- name: UpdateReadingLog :exec
UPDATE reading_logs
SET minutes_read = $2
WHERE userid = $1 and date = $3
`

type UpdateReadingLogParams struct {
	Userid      string    `json:"userid"`
	MinutesRead int32     `json:"minutes_read"`
	Date        time.Time `json:"date"`
}

func (q *Queries) UpdateReadingLog(ctx context.Context, arg UpdateReadingLogParams) error {
	_, err := q.db.ExecContext(ctx, updateReadingLog, arg.Userid, arg.MinutesRead, arg.Date)
	return err
}
