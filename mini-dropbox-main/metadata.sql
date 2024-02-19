DROP TABLE IF EXISTS file_metadata;

CREATE TABLE file_metadata (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    filename VARCHAR(255) NOT NULL,
    size_in_bytes BIGINT NOT NULL,
    s3_object_key VARCHAR(255) NOT NULL,
    description TEXT,
    mime_type VARCHAR(255),
    status TINYINT(4) NOT NULL DEFAULT 1,
    prev_key VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX active_files on file_metadata (filename, status);
CREATE INDEX trash_files on file_metadata (status, updated_at);
