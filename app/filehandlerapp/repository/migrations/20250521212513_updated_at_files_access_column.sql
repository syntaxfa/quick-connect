-- +migrate up
CREATE TRIGGER set_updated_at BEFORE
UPDATE
    ON files_access FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS set_updated_at ON files_access;