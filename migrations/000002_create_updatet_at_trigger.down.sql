DROP TRIGGER IF EXISTS set_updated_at_users_trigger ON auth.users;

DROP FUNCTION IF EXISTS auth.set_updated_at();
