// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bot.sql

package db

import (
	"context"
)

const enable_bot_registration = `-- name: enable_bot_registration :exec
update bot_settings set regsitration = true where id > 0
`

func (q *Queries) enable_bot_registration(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, enable_bot_registration)
	return err
}
