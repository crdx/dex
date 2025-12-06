-- name: FindLogsSince :many
select *
from logs
where id > ?;
