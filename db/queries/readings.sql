-- name: CreateReadingLog :exec
INSERT INTO reading_logs (userid, username, date, minutes_read)
VALUES (
           $1,$2, $3, $4
       );


-- name: UpdateReadingLog :exec
UPDATE reading_logs
SET minutes_read = $2
WHERE userid = $1;



-- name: GetReadingLogsByUser :many
SELECT date, minutes_read
FROM reading_logs
WHERE userid = $1
ORDER BY date DESC;




-- name: GetTopReaders :many
SELECT u.username, SUM(rl.minutes_read) AS total_minutes
FROM users u
         JOIN reading_logs rl ON u.userid = rl.userid
GROUP BY u.username
ORDER BY total_minutes DESC
    LIMIT 3;

-- name: GetTopReadersThisMonth :many
SELECT u.username, SUM(rl.minutes_read) AS total_minutes
FROM users u
         JOIN reading_logs rl ON u.userid = rl.userid
WHERE rl.date >= date_trunc('month', CURRENT_DATE)  -- Start of the current month
  AND rl.date < date_trunc('month', CURRENT_DATE + INTERVAL '1 month')  -- Start of next month
GROUP BY u.username
ORDER BY total_minutes DESC
    LIMIT 3;


-- name: GetTopStreaks :many
WITH consecutive_days AS (
    SELECT
        rl.userid,
        rl.date,
        ROW_NUMBER() OVER (PARTITION BY rl.userid ORDER BY rl.date)
            - EXTRACT(DAY FROM rl.date)::INT AS streak_group
    FROM reading_logs rl
),
     streaks AS (
         SELECT
             u.username,
             COUNT(*) AS streak_length,
             MAX(date) AS last_date
         FROM consecutive_days cd
                  JOIN users u ON cd.userid = u.userid
         GROUP BY u.username, streak_group
     )
SELECT
    username,
    streak_length
FROM streaks
WHERE last_date = CURRENT_DATE  -- Ensure the streak continues up to today
ORDER BY streak_length DESC
    LIMIT 3;

-- name: GetUserTopStreak :one
WITH consecutive_days AS (
    SELECT
        rl.userid,
        rl.date,
        ROW_NUMBER() OVER (PARTITION BY rl.userid ORDER BY rl.date)
            - EXTRACT(DAY FROM rl.date)::INT AS streak_group
    FROM reading_logs rl
    WHERE rl.userid = $1
),
     streaks AS (
         SELECT
             COUNT(*) AS streak_length,
             MAX(date) AS last_date
         FROM consecutive_days
         GROUP BY streak_group
     )
SELECT
    CASE
        WHEN MAX(streak_length) IS NULL THEN '0'
        ELSE CAST(MAX(streak_length) AS TEXT)
        END AS top_streak
FROM streaks
WHERE last_date = CURRENT_DATE;