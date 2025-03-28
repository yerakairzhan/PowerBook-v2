-- name: CreateUser :exec
insert into users (userid, username) values ($1, $2);

-- name: DeleteUserReged :exec
update users set registered = false where userid = $1;

-- name: GetLanguage :one
select language from users where userid = $1;

-- name: SetLanguage :exec
update users set language = $2 where userid = $1;

-- name: GetUserReged :one
select registered from users where userid = $1;

-- name: SetUserReged :exec
update users set registered = true where userid = $1;

-- name: SetUserState :exec
update users set state = $2 where userid = $1;

-- name: GetUserState :one
SELECT state FROM users WHERE userid = $1;

-- name: DeleteUserState :exec
update users set state = null where userid = $1;

-- name: GetUsersWithoutReadingToday :many
SELECT u.userid, u.language
FROM users u
         LEFT JOIN reading_logs r
                   ON u.userid = r.userid AND r.date = CURRENT_DATE
WHERE u.registered = TRUE AND r.userid IS NULL;


-- name: GetRegisteredUsers :many
select * from users where registered = true;