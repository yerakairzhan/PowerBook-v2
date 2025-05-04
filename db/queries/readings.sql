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
    SUM(r.minutes_read) AS total_minutes,
    COUNT(DISTINCT CASE WHEN r.minutes_read > 29 THEN r.date END) AS days_read_more_than_30
FROM reading_logs r
         JOIN users u ON r.userid = u.userid
WHERE EXTRACT(YEAR FROM r.date) = EXTRACT(YEAR FROM CURRENT_DATE)
  AND EXTRACT(MONTH FROM r.date) = EXTRACT(MONTH FROM CURRENT_DATE)
GROUP BY r.userid, u.username
ORDER BY days_read_more_than_30 DESC, total_minutes DESC
    LIMIT 5;

-- name: GetSumReading :one
SELECT
    SUM(minutes_read) AS Sum,
    username,
    userid,
    COUNT(DISTINCT CASE WHEN minutes_read > 29 THEN date END) AS days_read_more_than_30
FROM reading_logs
WHERE userid = '710606281'
  AND EXTRACT(YEAR FROM date) = EXTRACT(YEAR FROM CURRENT_DATE)
  AND EXTRACT(MONTH FROM date) = EXTRACT(MONTH FROM CURRENT_DATE)
GROUP BY userid, username;


