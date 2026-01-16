create table if not exists public.users
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
create table if not exists public.files
(
    id          bigserial
        primary key,
    name        varchar(255)                        not null,
    description text,
    uploaded_at timestamp default CURRENT_TIMESTAMP not null,
    embedding   vector(128),
    url         varchar(512)                        not null,
    size        numeric,
    user_id     bigint                              not null
        constraint fk_users_files
            references public.users
);


create table if not exists public.unique_faces
(
    id        bigserial
        primary key,
    name      text,
    embedding vector(128) not null
);


create table if not exists public.faces
(
    id             bigserial
        primary key,
    file_id        bigint             not null
        constraint fk_files_faces
            references public.files,
    unique_face_id bigint             not null
        constraint fk_unique_faces_faces
            references public.unique_faces,
    coordinates    double precision[] not null
);


create table if not exists public.jobs
(
    id         bigserial
        primary key,
    file_id    bigint                              not null
        constraint fk_files_jobs
            references public.files,
    status     varchar(50)                         not null,
    started_at timestamp default CURRENT_TIMESTAMP not null,
    ended_at   timestamp
);


create table if not exists public.devices
(
    id         bigserial
        primary key,
    user_id    bigint                              not null
        constraint fk_devices_user
            references public.users
            on delete cascade,
    device_id  varchar(255)                        not null,
    created_at timestamp default CURRENT_TIMESTAMP not null
);


