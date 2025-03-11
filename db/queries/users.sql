-- name: CreateUser :exec
insert into users (userid, username) values ($1, $2);

-- name: GetLanguage :one
select language from users where userid = $1;

-- name: SetLanguage :exec
update users set language = $2 where userid = $1;

-- name: GetUserReged :one
select registered from users where userid = $1;