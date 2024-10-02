CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    pem_file TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'provisioning' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);