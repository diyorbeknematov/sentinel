-- agents jadvali
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- metrics jadvali
CREATE TABLE IF NOT EXISTS metrics (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,
    cpu NUMERIC(5,2),
    ram NUMERIC(5,2),
    disk NUMERIC(5,2),
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_metrics_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id)
        ON DELETE CASCADE
);

-- applogs jadvali 
CREATE TABLE applogs (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,

    user_id VARCHAR(50),
    type VARCHAR(50),     -- login_failed, login_success
    level VARCHAR(20),     -- INFO, WARN, ERROR
    message TEXT,
    ip_address INET,   
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (agent_id)
    REFERENCES agents(id)
    ON DELETE CASCADE
);

-- nginxlogs jadvali 
CREATE TABLE nginxlogs (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL,
    ip_address INET,
    method VARCHAR(10),     -- GET, POST
    path TEXT,
    status INT,             -- 200, 404, 500
    response_time INT,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

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