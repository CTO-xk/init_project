-- 链状态表：记录各链最后处理的区块 (MySQL)
CREATE TABLE IF NOT EXISTS chain_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chain_name VARCHAR(50) NOT NULL UNIQUE,
    last_processed_block BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 用户总余额表 (MySQL)
CREATE TABLE IF NOT EXISTS user_balances (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chain_name VARCHAR(50) NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    current_balance VARCHAR(100) NOT NULL DEFAULT '0',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_chain_user (chain_name, user_address)
);

-- 余额变动记录表 (MySQL)
CREATE TABLE IF NOT EXISTS balance_changes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chain_name VARCHAR(50) NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    event_type VARCHAR(20) NOT NULL,
    amount VARCHAR(100) NOT NULL,
    balance_after VARCHAR(100) NOT NULL,
    block_number BIGINT NOT NULL,
    event_time TIMESTAMP NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_chain_addr_time (chain_name, user_address, event_time)
);

-- 用户总积分表 (MySQL)
CREATE TABLE IF NOT EXISTS user_points (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chain_name VARCHAR(50) NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    total_points DECIMAL(30,6) NOT NULL DEFAULT 0,
    last_calculated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_chain_user (chain_name, user_address)
);

-- 积分计算历史表 (MySQL)
CREATE TABLE IF NOT EXISTS points_calculation_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chain_name VARCHAR(50) NOT NULL,
    user_address VARCHAR(42) NOT NULL,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    points_added DECIMAL(30,6) NOT NULL,
    total_points DECIMAL(30,6) NOT NULL,
    calculated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_chain_addr_period (chain_name, user_address, period_start)
);
