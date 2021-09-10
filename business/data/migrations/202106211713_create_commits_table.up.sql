create table commits
(
    id varchar(255) not null
        constraint commits_pkey
            primary key,
    author varchar(255) not null,
    message text,
    app_id uuid not null
        constraint fk_commit_app
            references apps
            on delete cascade,
    date timestamp
);

create index commits_app_id_idx
    on commits (app_id);

