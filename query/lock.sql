-- name: TryAdvisoryLock :one
SELECT pg_try_advisory_xact_lock(hashtext(sqlc.arg(key)));
