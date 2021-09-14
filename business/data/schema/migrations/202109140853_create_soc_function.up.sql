create or replace function soc(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                entity text,
                soc    bigint
            )
as
$body$
select f1 as entity, count(f1) as soc
from (
         select s1.file as f1
         from stats s1
                  inner join commits c on c.id = s1.commit_id
             and c.date between $3 and $2
                  inner join stats s2 on s1.commit_id = s2.commit_id
             and s1.file not like '%=>%'
             and s2.file not like '%=>%'
             and s1.file != s2.file
             and s1.app_id = $1
     ) a
group by f1
order by count(f1) desc
$body$
    language sql;