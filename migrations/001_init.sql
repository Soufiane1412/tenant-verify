-- simple but powerful schema
CREATE TABLE IF NOT EXISTS verifications (

    id VARCHAR(50) PRIMARY KEY,
    tenant_name VARCHAR(255) NOT NULL,
    tenant_email VARCHAR(255) NOT NULL,
    income INTEGER NOT NULL,
    risk_score INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    verified_at TIMESTAMP NOT NULL,
    details TEXT[], -- Array of detail messages
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 

CREATE INDEX idx_status_date ON verifications(status, verified_at DESC);
CREATE INDEX idx_email ON verifications(tenant_email);