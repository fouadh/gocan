create or replace function revisions(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                entity                      text,
                numberOfRevisions           bigint,
                numberOfAuthors             bigint,
                normalizedNumberOfRevisions decimal,
                code                        integer
            )
as
$body$
select entity,
       numberOfRevisions,
       numberOfAuthors,
       cast(numberOfRevisions as decimal) / MAX(numberOfRevisions) OVER () as normalizedNumberOfRevisions,
       coalesce((
                    select lines
                    from cloc
                             inner join commits c2 on cloc.commit_id = c2.id
                    where cloc.file = entity
                      and cloc.app_id = $1
                      and c2.date >= $2 - interval '1 DAY'
                    order by c2.date asc
                    limit 1
                ), 0)                                                      as code
FROM (
         select file                     as entity,
                count(distinct c.author) as numberOfAuthors,
                COUNT(*)                 as numberOfRevisions
         from stats
                  inner join commits c on c.id = stats.commit_id
             and c.app_id = stats.app_id
             and stats.app_id = $1
             and c.date between $3 and $2
         WHERE FILE NOT LIKE '%%=>%%'
         group by stats.app_id, file
     ) q
ORDER BY numberOfRevisions desc
$body$
    language sql;