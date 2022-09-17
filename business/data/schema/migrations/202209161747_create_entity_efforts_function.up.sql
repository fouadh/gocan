drop function entity_efforts(app_id uuid, before timestamp, after timestamp);
create or replace function entity_efforts(app_id uuid, before timestamp, after timestamp)
    returns table
            (
                entity          text,
                author          text,
                team            text,
                authorRevisions bigint,
                totalRevisions  bigint
            )
as
$body$
select entity,
       author,
       coalesce((select name
        from teams t
                 inner join team_members tm on t.id = tm.team_id
            and tm.member_name = a.author
        where app_id = $1), '')       team,
       authorRevisions,
       (select count(*)
        from stats
                 inner join commits c on c.id = stats.commit_id
            and c.date between $3 and $2
        where file = entity
          and stats.app_id = $1) totalRevisions
from (SELECT file     as entity,
             c.author as author,
             count(*) as authorRevisions
      FROM stats
               inner join commits c on c.id = stats.commit_id
          and c.date between $3 and $2
      WHERE stats.app_id = $1
        AND FILE NOT LIKE '%=>%'
      GROUP BY c.author, file) a
order by entity asc
$body$
    language sql;