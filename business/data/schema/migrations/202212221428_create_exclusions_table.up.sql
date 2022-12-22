create table exclusions(
    app_id uuid not null,
    exclusion varchar(255) not null,
    primary key(app_id, exclusion)
)