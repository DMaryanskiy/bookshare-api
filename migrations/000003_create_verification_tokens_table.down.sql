DROP TABLE IF EXISTS auth.verification_tokens;

DROP TRIGGER IF EXISTS set_updated_at_tokens_trigger ON auth.verification_tokens;

DROP FUNCTION IF EXISTS auth.set_token_updated_at();
