drop table couplings;

create or replace function couplings(app_id uuid, before timestamp, after timestamp, minimalCoupling float,
                                     minimalAverageRevisions float)
    returns table
            (
                entity           text,
                coupled          text,
                degree           decimal,
                averageRevisions decimal
            )
as
$body$
select entity,
       coupled,
       degree,
       average_revs averageRevisions
from (select entity,
             coupled,
             t0,
             t1,
             t2,
             (cast(t0 as decimal) / (cast((t1 + t2) as decimal) / 2)) as degree,
             cast((t1 + t2) as decimal) / 2                           as average_revs
      from (
               select entity,
                      coupled,
                      count(*)              as t0,
                      (select count(*)
                       from stats s
                                inner join commits c on s.commit_id = c.id and c.date between $3 and $2
                       where file = entity
                         and s.app_id = $1) as t1,
                      (select count(*)
                       from stats s
                                inner join commits c on s.commit_id = c.id and c.date between $3 and $2
                       where file = coupled
                         and s.app_id = $1) as t2
               from (
                        select s1.file as entity, s2.file as coupled
                        from stats s1
                                 inner join commits c
                                            on s1.commit_id = c.id
                                                and c.date between $3 and $2
                                 inner join stats s2
                                            on s1.commit_id = s2.commit_id
                        where s1.file < s2.file
                          and s1.file not like '%%=>%%'
                          and s2.file not like '%%=>%%'
                          and s1.app_id = $1) a
               group by entity, coupled
           ) b
     ) c
where degree >= $4
  and average_revs >= $5
order by degree desc
$body$
    language sql;