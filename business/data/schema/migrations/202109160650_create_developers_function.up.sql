create or replace function developers(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                name            text,
                numberOfCommits bigint
            )
as
$body$
select author as name, count(1) as numberOfCommits
from commits
where app_id = $1
  and date between $3 and $2
group by author
order by name asc
$body$
    language sql;