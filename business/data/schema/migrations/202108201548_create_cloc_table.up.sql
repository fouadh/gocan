create table cloc
(
    app_id uuid not null
        constraint fk_cloc_app
            references apps
            on delete cascade,
    commit_id varchar(255) not null,
    file text not null,
    lines integer not null
);

create index cloc_app_id_commit_id_idx
    on cloc (app_id, commit_id);