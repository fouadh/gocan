create table scenes
(
    id uuid not null
        constraint scenes_pkey
            primary key,
    name varchar(255)
        constraint unique_scene_name
            unique,
    created_at timestamp without time zone default (now() at time zone 'utc')
);
