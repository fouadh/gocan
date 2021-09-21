create table boundaries
(
    id uuid not null
        constraint boundaries_pkey
            primary key,
    name varchar(255) not null,
    app_id uuid not null
        constraint fk_boundaries_app
            references apps,
    constraint unique_boundaries_name
        unique (app_id, name)
);

create table transformations
(
    name varchar(255) not null,
    path varchar(255) not null,
    boundary_id uuid not null
        constraint fk_transformations_boundary
            references boundaries
);


