create table if not exists users
(
    id         bigserial
        primary key,
    name       text,
    email      text,
    password   text,
    created_at timestamp with time zone
);

create table if not exists files
(
    id          bigserial
        primary key,
    name        text,
    description text,
    uploaded_at timestamp with time zone,
    embedding   vector(128),
    size        numeric,
    user_id     bigint
        constraint fk_users_files
            references users
);

create table if not exists unique_faces
(
    id        bigserial
        primary key,
    name      text,
    embedding vector(128)
);

create table if not exists faces
(
    id             bigserial
        primary key,
    file_id        bigint
        constraint fk_files_faces
            references files,
    unique_face_id bigint
        constraint fk_unique_faces_faces
            references unique_faces,
    coordinates    double precision[]
);

create table if not exists jobs
(
    id         bigserial
        primary key,
    file_id    bigint
        constraint fk_files_jobs
            references files,
    status     text,
    started_at timestamp with time zone,
    ended_at   timestamp with time zone
);


create index if not exists idx_files_user_id
    on files (user_id);

create index if not exists idx_faces_file_id
    on faces (file_id);

create index if not exists idx_faces_unique_face_id
    on faces (unique_face_id);

create index if not exists idx_jobs_file_id
    on jobs (file_id);
