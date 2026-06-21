-- SQL schema for APMS Task Tracker

CREATE TABLE users (
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id)
);

CREATE TABLE projects (
    project_id VARCHAR(36) NOT NULL,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    manager_id varchar(36) REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(project_id)
);

CREATE TABLE skills (
    skill_id VARCHAR(36) NOT NULL ,
    name VARCHAR(50) UNIQUE NOT NULL,
    PRIMARY KEY (skill_id)
);

CREATE TABLE developer_skills (
    developer_id VARCHAR(36) REFERENCES users(user_id) ON DELETE CASCADE,
    skill_id VARCHAR(36) REFERENCES skills(skill_id) ON DELETE CASCADE,
    PRIMARY KEY (developer_id, skill_id)
);

CREATE TABLE tasks (
    task_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL REFERENCES projects(project_id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'backlog',
    priority INT DEFAULT 2,
    estimated_hours INT NOT NULL,
    required_skill_id varchar(36) REFERENCES skills(skill_id),
    executor_id varchar(36) REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (task_id)
);

CREATE TABLE task_dependencies (
    blocked_task_id VARCHAR(36) NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    dependent_on_id VARCHAR(36) NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    PRIMARY KEY (blocked_task_id, dependent_on_id)
);

CREATE TABLE roles (
    role_id VARCHAR(36) NOT NULL,
    title VARCHAR(20) NOT NULL,
    status VARCHAR(64) DEFAULT 'active',
    PRIMARY KEY (role_id)
);

create table role_permissions (
    role_id varchar(36) NOT NULL ,
    permission varchar(128) NOT NULL ,
    PRIMARY KEY (role_id, permission)
);
