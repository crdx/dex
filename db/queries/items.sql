-- name: FindItemsForList :many
select
    i.id,
    b.content
from items i
join blobs b on b.hash = i.blob_hash
where i.deleted_at is null
order by i.last_hit_at asc;

-- name: FindItemByRef :one
select *
from items
where deleted_at is null
and (label = ? or uuid = ?);

-- name: IncrementHits :exec
update items
set hits = hits + 1
where id = ?;
