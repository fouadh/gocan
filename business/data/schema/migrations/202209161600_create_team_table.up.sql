create table teams
(
    id     uuid         not null
        constraint teams_pkey primary key,
    name   varchar(255) not null,
    app_id uuid         not null
        constraint fk_team_app
            references apps
            on delete cascade
);

create table team_members
(
    member_name varchar(255) not null,
    team_id      uuid         not null
        constraint fk_team_member_team
            references teams
            on delete cascade
);

create index team_app_id_idx
    on teams (app_id);

create index team_member_team_id_idx
    on team_members (team_id);
