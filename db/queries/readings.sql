-- name: CreateReadingLog :exec
INSERT INTO reading_logs (userid, username, date, minutes_read)
VALUES (
           $1,$2, $3, $4
       );


-- name: UpdateReadingLog :exec
UPDATE reading_logs
SET minutes_read = $2
WHERE userid = $1 and date = $3;

-- name: GetReadingLogsByUser :many
SELECT date, minutes_read
FROM reading_logs
WHERE userid = $1
ORDER BY date DESC;

-- name: GetReadingLeaderboard :many
SELECT
    r.userid,
    u.username,
    SUM(r.minutes_read) AS total_minutes
FROM reading_logs r
         JOIN users u ON r.userid = u.userid
WHERE r.date >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY r.userid, u.username
ORDER BY total_minutes DESC
    LIMIT 5;

-- name: GetSumReading :one
select sum(minutes_read) as Sum, username, userid from reading_logs where userid = $1 group by userid, username;

