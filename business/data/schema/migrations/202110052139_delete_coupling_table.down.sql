create table couplings
(
    entity           text not null,
    coupled          text not null,
    degree           decimal,
    averageRevisions decimal,
    app_id           uuid not null
        constraint fk_couplings_app
            references apps
            on delete cascade
);

create index couplings_app_degree_revisions_idx
    on couplings (app_id, degree, averageRevisions);