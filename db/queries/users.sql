-- name: CreateUser :exec
insert into users (userid, username) values ($1, $2);

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
