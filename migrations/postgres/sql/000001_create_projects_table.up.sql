CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY, 
    name VARCHAR(60) NOT NULL, 
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS projects_id_idx ON projects(id);

INSERT INTO projects(name) VALUES ('Первая запись');
