alter table cloc
    add constraint unique_cloc_app_commit unique(app_id, commit_id, file);