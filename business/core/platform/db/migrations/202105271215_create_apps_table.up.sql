create table apps
(
    id uuid not null
        constraint apps_pkey
            primary key,
    name varchar(255) not null,
    scene_id uuid not null
        constraint fk_app_scene
            references scenes
            on delete cascade
);

create index apps_scene_id_idx
    on apps (scene_id);

