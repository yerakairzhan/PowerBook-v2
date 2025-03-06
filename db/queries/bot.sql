-- name: enable_bot_registration :exec
update bot_settings set regsitration = true where id > 0;