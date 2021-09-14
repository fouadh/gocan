create or replace function main_developers(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                entity     text,
                author     text,
                added      bigint,
                totalAdded bigint,
                ownership  decimal
            )
as
$body$
select entity,
       name                                as author,
       added,
       totalAdded,
       cast(added as decimal) / totalAdded as ownership
from (select entity,
             name,
             added,
             (select sum(insertions)
              from stats
                       inner join commits c on c.id = stats.commit_id
                  and c.date between $3 and $2
              where file = entity
                and stats.app_id = $1)                                   as totalAdded,
             row_number() over (partition by entity order by added desc) as row
      from (SELECT file            as entity,
                   c.author        as name,
                   sum(insertions) as added
            FROM stats
                     inner join commits c on c.id = stats.commit_id
                and c.date between $3 and $2
            WHERE stats.app_id = $1
              AND FILE NOT LIKE '%=>%'
            GROUP BY c.author, file
           ) a) b
WHERE added > 0
  AND row = 1
ORDER BY entity asc, ownership desc
$body$
    language sql;