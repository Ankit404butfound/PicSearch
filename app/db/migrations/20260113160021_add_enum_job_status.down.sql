DROP TYPE job_status;

ALTER TABLE jobs
    DROP COLUMN face_encoding_status,
    DROP COLUMN universal_encoding_status;

ALTER TABLE jobs
    ADD COLUMN status text;