create table revision_trends
(
    id          uuid not null
        constraint revision_trends_pkey
            primary key,
    name        varchar(255),
    boundary_id uuid not null
        constraint fk_trend_boundary
            references boundaries
            on delete cascade,
    created_at  timestamp without time zone default (now() at time zone 'utc')
);

create table revision_trend_entries
(
    id                uuid        not null
        constraint revision_trend_entries_pkey primary key,
    revision_trend_id uuid        not null
        constraint fk_revision_trend_entry_revision_trend
            references revision_trends
            on delete cascade,
    date              varchar(10) not null,
    created_at        timestamp without time zone default (now() at time zone 'utc')
);

create table revision_trend_entry_revisions
(
    entry_id            uuid not null
        constraint fk_revision_trend_entry_revision_revision_trend_entries
            references revision_trend_entries on delete cascade,
    entity              text not null,
    number_of_revisions int  not null,
    created_at          timestamp without time zone default (now() at time zone 'utc')
);

create index revision_trends_idx
    on revision_trends (name, boundary_id);