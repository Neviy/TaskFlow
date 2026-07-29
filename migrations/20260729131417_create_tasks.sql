-- +goose Up
CREATE TABLE tasks (

    id BIGSERIAL PRIMARY KEY,


    project_id BIGINT NOT NULL,


    title VARCHAR(200) NOT NULL,


    description TEXT,


    status VARCHAR(20) NOT NULL DEFAULT 'todo',


    priority VARCHAR(20) NOT NULL DEFAULT 'medium',


    assignee_id BIGINT,


    deadline TIMESTAMP,


    created_at TIMESTAMP NOT NULL DEFAULT NOW(),


    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_task_project
        FOREIGN KEY(project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,


    CONSTRAINT fk_task_assignee
        FOREIGN KEY(assignee_id)
        REFERENCES users(id)
        ON DELETE SET NULL,


    CONSTRAINT check_task_status
        CHECK(
            status IN(
                'todo',
                'in_progress',
                'done'
            )
        ),


    CONSTRAINT check_task_priority
        CHECK(
            priority IN(
                'low',
                'medium',
                'high',
                'critical'
            )
        )
);

-- +goose Down
DROP TABLE IF EXISTS tasks;
