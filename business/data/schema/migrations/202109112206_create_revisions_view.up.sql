create or replace view revisions
as
select app_id,
       entity,
       numberOfRevisions,
       numberOfAuthors,
       cast(numberOfRevisions as decimal) / MAX(numberOfRevisions) OVER ()       as normalizedNumberOfRevisions,
       coalesce((select lines from cloc where file = entity and app_id = app_id), 0) as code
FROM (
         select stats.app_id             as app_id,
                file                     as entity,
                count(distinct c.author) as numberOfAuthors,
                COUNT(*)                 as numberOfRevisions
         from stats
                  inner join commits c on c.id = stats.commit_id
             and c.app_id = stats.app_id
         WHERE FILE NOT LIKE '%%=>%%'
         group by stats.app_id, file
     ) q
ORDER BY numberOfRevisions desc;