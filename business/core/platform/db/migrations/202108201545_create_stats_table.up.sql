create table stats
(
    commit_id varchar(255) not null
        constraint fk_stat_commit
            references commits,
    file text not null,
    insertions integer,
    deletions integer,
    app_id uuid not null
        constraint fk_stat_app
            references apps
            on delete cascade,
    constraint unique_app_commit_file_stat
        unique (app_id, commit_id, file)
);

create index stats_app_id_idx
    on stats (app_id);

create index stats_app_file_idx
    on stats (app_id, file);

create index stats_file
    on stats (file);

