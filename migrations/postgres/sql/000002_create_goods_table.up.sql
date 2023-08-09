CREATE TABLE IF NOT EXISTS goods (
    id BIGSERIAL,
    project_id BIGINT,
    name VARCHAR(60),
    description VARCHAR(120),
    priority INT DEFAULT 1,
    removed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id, project_id),
    CONSTRAINT fk_project FOREIGN KEY(project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS goods_id_project_id_idx ON goods(id, project_id);

CREATE INDEX IF NOT EXISTS goods_name_idx ON goods(name);

CREATE OR REPLACE FUNCTION update_priorities() RETURNS TRIGGER AS $$
DECLARE
    priority_diff INTEGER;
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.priority := COALESCE((SELECT MAX(priority) FROM goods), 0) + 1;
    ELSIF TG_OP = 'UPDATE' AND NEW.priority <> OLD.priority THEN
        UPDATE goods
        SET priority = priority + priority_diff
        WHERE priority >= NEW.priority AND id != NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_priorities
BEFORE INSERT OR UPDATE OF priority ON goods
FOR EACH ROW
EXECUTE FUNCTION update_priorities();
