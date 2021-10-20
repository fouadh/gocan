alter table revision_trends
    add column app_id uuid
        constraint fk_trends_app references apps on delete cascade;

update revision_trends
set app_id=(select app_id from boundaries where id=revision_trends.boundary_id);

alter table revision_trends
    add constraint unique_trends_app unique(app_id, name);

create index revision_trends_app_name_idx
    on revision_trends (name, app_id);

create index revision_trends_app_idx
    on revision_trends (app_id);