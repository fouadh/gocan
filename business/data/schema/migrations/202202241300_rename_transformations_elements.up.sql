alter table transformations
    rename to modules;

drop view boundaries_transformations;

create view boundaries_modules as
select id,
       name,
       app_id,
       (select array_to_json(array_agg(row_to_json(modulesList.*))) as array_to_json
        from (
                 select modules.name, modules.path
                 from modules
                 where boundary_id = boundaries.id
             ) modulesList
       ) as modules
from boundaries;
