CREATE TYPE job_status AS ENUM ('pending', 'in_progress', 'completed', 'failed');

create table users
(
    id         bigserial
        primary key,
    name       varchar(255)                        not null,
    email      varchar(255)                        not null
        constraint uni_users_email
            unique,
    password   varchar(255)                        not null,
    created_at timestamp default CURRENT_TIMESTAMP not null
);


create table files
(
    id          bigserial
        primary key,
    name        varchar(255)                        not null,
    description text,
    uploaded_at timestamp default CURRENT_TIMESTAMP not null,
    embedding   text,
    url         varchar(512)                        not null,
    size        numeric,
    user_id     bigint                              not null
        constraint fk_users_files
            references public.users,
    metadata    jsonb
);

create table unique_faces
(
    id        bigserial
        primary key,
    name      text,
    embedding vector(128) not null,
    image_url text
);

create table faces
(
    id             bigserial
        primary key,
    file_id        bigint             not null
        constraint fk_files_faces
            references files,
    unique_face_id bigint             not null
        constraint fk_unique_faces_faces
            references unique_faces,
    coordinates    double precision[] not null
);


create table if not exists jobs
(
    id                        bigserial
        primary key,
    file_id                   bigint                               not null
        constraint fk_files_jobs
            references files,
    started_at                timestamp  default CURRENT_TIMESTAMP not null,
    ended_at                  timestamp,
    face_encoding_status      job_status    default 'pending'::job_status,
    universal_encoding_status job_status default 'pending'::job_status
);

create table if not exists devices
(
    id         bigserial
        primary key,
    user_id    bigint                              not null
        constraint fk_devices_user
            references users
            on delete cascade,
    device_id  varchar(255)                        not null,
    created_at timestamp default CURRENT_TIMESTAMP not null
);


