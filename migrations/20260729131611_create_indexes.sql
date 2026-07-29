-- +goose Up

CREATE INDEX IF NOT EXISTS idx_projects_owner_id
ON projects(owner_id);


CREATE INDEX IF NOT EXISTS idx_project_members_project_id
ON project_members(project_id);


CREATE INDEX IF NOT EXISTS idx_project_members_user_id
ON project_members(user_id);


CREATE INDEX IF NOT EXISTS idx_tasks_project_id
ON tasks(project_id);


CREATE INDEX IF NOT EXISTS idx_tasks_status
ON tasks(status);


CREATE INDEX IF NOT EXISTS idx_tasks_priority
ON tasks(priority);


CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id
ON tasks(assignee_id);


CREATE INDEX IF NOT EXISTS idx_comments_task_id
ON comments(task_id);


CREATE INDEX IF NOT EXISTS idx_history_task_id
ON task_history(task_id);


-- +goose Down

DROP INDEX IF EXISTS idx_projects_owner_id;

DROP INDEX IF EXISTS idx_project_members_project_id;

DROP INDEX IF EXISTS idx_project_members_user_id;

DROP INDEX IF EXISTS idx_tasks_project_id;

DROP INDEX IF EXISTS idx_tasks_status;

DROP INDEX IF EXISTS idx_tasks_priority;

DROP INDEX IF EXISTS idx_tasks_assignee_id;

DROP INDEX IF EXISTS idx_comments_task_id;

DROP INDEX IF EXISTS idx_history_task_id;