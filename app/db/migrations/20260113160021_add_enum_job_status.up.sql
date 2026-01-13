CREATE TYPE job_status AS ENUM ('pending', 'in_progress', 'completed', 'failed');

ALTER TABLE jobs
    DROP COLUMN status;

ALTER TABLE jobs
    ADD COLUMN face_encoding_status job_status DEFAULT 'pending';

ALTER TABLE jobs
    ADD COLUMN universal_encoding_status job_status DEFAULT 'pending';