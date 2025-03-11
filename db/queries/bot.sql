-- name: CreateBot :exec
insert into bot_settings (id) select 1 where not exists (select 1 from bot_settings);

-- name: Enable_bot_registration :exec
update bot_settings set registration = true where registration = false;

-- name: Diasble_bot_registration :exec
update bot_settings set registration = false where registration = true;

-- name: Getbot :one
select registration from bot_settings where id = 1;