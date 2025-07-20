-- +migrate StatementBegin
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON subscriber_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS set_updated_at ON subscriber_profiles;