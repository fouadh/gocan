create or replace function app_summary(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                name                    text,
                id                      uuid,
                numberOfCommits         bigint,
                numberOfEntities        bigint,
                numberOfEntitiesChanged bigint,
                numberOfAuthors         bigint
            )
as
$body$
select a.name,
       (select $1)            as id,
       count(distinct c.*)       numberOfCommits,
       count(distinct s.file) as numberOfEntities,
       count(s.*)                numberOfEntitiesChanged,
       count(distinct c.author)  numberOfAuthors
from commits c
         inner join stats s
                    on c.id = s.commit_id
                        and c.app_id = s.app_id
                        and c.date between $3 and $2
         inner join apps a on a.id = s.app_id
where a.id = $1
  and c.date between $3 and $2
  and s.file not like '%=>%'
group by a.name
$body$
    language sql;