create table exclusions(
    app_id uuid not null constraint fk_exclusions_app
        references apps
        on delete cascade,
    exclusion varchar(255) not null,
    primary key(app_id, exclusion)
)