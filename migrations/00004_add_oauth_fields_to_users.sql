-- +goose Up
-- +goose StatementBegin
-- Make password_hash nullable for OAuth users who don't have passwords
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Add OAuth-specific fields
ALTER TABLE users ADD COLUMN provider VARCHAR(50) DEFAULT 'local'; -- 'google' or 'local'
ALTER TABLE users ADD COLUMN provider_id VARCHAR(255); -- Google's user ID
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(512); -- Profile picture URL
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

-- Create unique constraint for provider + provider_id combination
-- This prevents duplicate OAuth users
CREATE UNIQUE INDEX idx_users_provider_id ON users(provider, provider_id)
WHERE provider IS NOT NULL AND provider != 'local';

-- Add comment for clarity
COMMENT ON COLUMN users.provider IS 'Authentication provider: local (email/password) or google (OAuth)';
COMMENT ON COLUMN users.provider_id IS 'Unique user ID from OAuth provider (NULL for local users)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_provider_id;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS provider_id;
ALTER TABLE users DROP COLUMN IF EXISTS provider;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
-- +goose StatementEnd
