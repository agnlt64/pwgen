
CREATE TABLE vault (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    salt VARCHAR(24) NOT NULL, -- base64(16 bytes) = 24 bytes
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE vault_entry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ciphertext VARCHAR(255) NOT NULL,
    nonce VARCHAR(255) NOT NULL,
    
    -- user fields
    website VARCHAR(255) NOT NULL,
    label VARCHAR(255),
    
    -- audit
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,

    vault_id UUID NOT NULL REFERENCES vault(id) ON DELETE CASCADE
);