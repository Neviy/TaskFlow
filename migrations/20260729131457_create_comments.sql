-- +goose Up
CREATE TABLE comments (

    id BIGSERIAL PRIMARY KEY,


    task_id BIGINT NOT NULL,


    author_id BIGINT NOT NULL,


    text TEXT NOT NULL,


    created_at TIMESTAMP NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_comment_task
        FOREIGN KEY(task_id)
        REFERENCES tasks(id)
        ON DELETE CASCADE,


    CONSTRAINT fk_comment_author
        FOREIGN KEY(author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS comments;
