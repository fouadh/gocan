create view boundaries_transformations as
select id,
       name,
       app_id,
       (select array_to_json(array_agg(row_to_json(transformationsList.*))) as array_to_json
        from (
                 select transformations.name, transformations.path
                 from transformations
                 where boundary_id = boundaries.id
             ) transformationsList
       ) as transformations
from boundaries;
