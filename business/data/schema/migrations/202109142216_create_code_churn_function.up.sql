create or replace function code_churn(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                date    text,
                added   bigint,
                deleted bigint
            )
as
$body$
select date, sum(insertions) added, sum(deletions) deleted
from (
         SELECT to_char(date_trunc('day', date::date), 'YYYY-MM-DD') "date", insertions, deletions
         from stats s
                  inner join commits c on c.id = s.commit_id
             and s.app_id = $1
             and c.date between $3 and $2
     ) a
group by 1
order by 1 asc
$body$
    language sql;