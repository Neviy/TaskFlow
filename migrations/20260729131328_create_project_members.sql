-- +goose Up
CREATE TABLE project_members (

    id BIGSERIAL PRIMARY KEY,

    project_id BIGINT NOT NULL,

    user_id BIGINT NOT NULL,

    role VARCHAR(20) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_project_member_project
        FOREIGN KEY(project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,


    CONSTRAINT fk_project_member_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,


    CONSTRAINT unique_project_user
        UNIQUE(project_id, user_id),


    CONSTRAINT check_member_role
        CHECK(
            role IN(
                'owner',
                'admin',
                'member',
                'viewer'
            )
        )
);

-- +goose Down
DROP TABLE IF EXISTS project_members;
