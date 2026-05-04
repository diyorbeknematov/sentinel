-- accounts jadvali 
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    api_key TEXT UNIQUE
);

-- agents jadvali
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    ip_address INET,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_agents_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE
);

-- metrics jadvali
CREATE TABLE IF NOT EXISTS metrics (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,
    cpu NUMERIC(5,2),
    ram NUMERIC(5,2),
    disk NUMERIC(5,2),

    log_time TIMESTAMP,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_metrics_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id)
        ON DELETE CASCADE
);

-- applogs jadvali 
CREATE TABLE IF NOT EXISTS applogs (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,

    user_id VARCHAR(50),
    event VARCHAR(50),
    level VARCHAR(20),
    message TEXT,

    log_time TIMESTAMP,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_applogs_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id)
        ON DELETE CASCADE
);

-- nginxlogs jadvali 
CREATE TABLE IF NOT EXISTS nginxlogs (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,

    ip_address INET,
    method VARCHAR(10),
    path TEXT,
    status INT,
    bytes INT,
    user_agent TEXT,

    log_time TIMESTAMP,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_nginxlogs_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id)
        ON DELETE CASCADE
);

-- alerts jadvali
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,
    type VARCHAR(50),
    message TEXT,
    severity VARCHAR(20),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_alerts_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id)
        ON DELETE CASCADE
);