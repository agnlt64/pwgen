
CREATE TABLE vault (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name VARCHAR(255) NOT NULL,
    salt VARCHAR(24) NOT NULL, -- base64(16 bytes) = 24 bytes
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE current_vault (
    singleton BOOLEAN PRIMARY KEY DEFAULT TRUE,
    current_vault_id UUID NOT NULL REFERENCES vault(id)
    CONSTRAINT singleton_check CHECK (singleton = TRUE)
);

CREATE TABLE vault_entry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ciphertext VARCHAR(255) NOT NULL,
    nonce VARCHAR(255) NOT NULL,
    
    -- user fields
    website VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    
    -- audit
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,

    vault_id UUID NOT NULL REFERENCES vault(id) ON DELETE CASCADE
);