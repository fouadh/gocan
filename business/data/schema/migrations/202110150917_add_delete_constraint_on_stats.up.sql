alter table stats
    drop constraint fk_stat_commit,
    add constraint fk_stat_commit
        foreign key (commit_id)
            references commits (id)
            on delete cascade;