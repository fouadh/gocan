create table complexity_analysis
(
    id         uuid         not null
        constraint complexity_analysis_pkey
            primary key,
    name       varchar(255) not null,
    entity     text         not null,
    app_id     uuid         not null
        constraint fk_complexity_analysis_app
            references apps
            on delete cascade,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    UNIQUE (app_id, name)
);

create index complexity_analysis_app_id_idx
    on complexity_analysis (app_id);

create table complexity_analysis_entries
(
    complexity_analysis_id uuid      not null
        references complexity_analysis
            on delete cascade,
    date                   timestamp not null,
    lines                  int       not null,
    indentations           int       not null,
    mean                   decimal   not null,
    stdev                  decimal   not null,
    max                    int       not null,
    created_at             timestamp without time zone default (now() at time zone 'utc')
);

