-- +goose Up
CREATE TABLE task_history (

    id BIGSERIAL PRIMARY KEY,


    task_id BIGINT NOT NULL,


    changed_by BIGINT NOT NULL,


    old_status VARCHAR(20),


    new_status VARCHAR(20),


    created_at TIMESTAMP NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_history_task
        FOREIGN KEY(task_id)
        REFERENCES tasks(id)
        ON DELETE CASCADE,


    CONSTRAINT fk_history_user
        FOREIGN KEY(changed_by)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS task_history;
