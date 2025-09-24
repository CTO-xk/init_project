package db

import (
	"database/sql"
	"time"
)

// BalanceChange 余额变动记录
type BalanceChange struct {
	ChainName    string
	UserAddress  string
	EventType    string
	Amount       string
	BalanceAfter string
	BlockNumber  uint64
	EventTime    time.Time
	TxHash       string
}

// GetLastProcessedBlock 获取链最后处理的区块
func GetLastProcessedBlock(chainName string) (uint64, error) {
	var block uint64
	err := QueryRow(
		"SELECT last_processed_block FROM chain_status WHERE chain_name = ?",
		chainName,
	).Scan(&block)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return block, err
}

// UpdateLastProcessedBlock 更新链最后处理的区块
func UpdateLastProcessedBlock(chainName string, block uint64) error {
	_, err := Exec(
		"UPDATE chain_status SET last_processed_block = ?, updated_at = CURRENT_TIMESTAMP WHERE chain_name = ?",
		block, chainName,
	)
	return err
}

// GetUserCurrentBalance 获取用户当前余额
func GetUserCurrentBalance(chainName, userAddr string) (string, error) {
	var balance string
	err := QueryRow(
		"SELECT current_balance FROM user_balances WHERE chain_name = ? AND user_address = ?",
		chainName, userAddr,
	).Scan(&balance)
	if err == sql.ErrNoRows {
		return "0", nil
	}
	return balance, err
}

// RecordBalanceChange 记录余额变动
func RecordBalanceChange(change BalanceChange) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入变动记录
	_, err = TxExec(tx, `
        INSERT INTO balance_changes (
            chain_name, user_address, event_type, amount, balance_after,
            block_number, event_time, tx_hash
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `,
		change.ChainName, change.UserAddress, change.EventType,
		change.Amount, change.BalanceAfter, change.BlockNumber,
		change.EventTime, change.TxHash,
	)
	if err != nil {
		return err
	}
	// 更新用户余额
	_, err = TxExec(tx, `
        INSERT INTO user_balances (chain_name, user_address, current_balance)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE current_balance = VALUES(current_balance), updated_at = CURRENT_TIMESTAMP
    `,
		change.ChainName, change.UserAddress, change.BalanceAfter,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// GetBalanceChangesInPeriod 获取指定时间段的余额变动
func GetBalanceChangesInPeriod(chainName, userAddr string, start, end time.Time) ([]BalanceChange, error) {
	rows, err := Query(`
        SELECT chain_name, user_address, event_type, amount, balance_after,
               block_number, event_time, tx_hash
        FROM balance_changes
        WHERE chain_name = ?
          AND user_address = ?
          AND event_time BETWEEN ? AND ?
        ORDER BY event_time ASC
    `, chainName, userAddr, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var changes []BalanceChange
	for rows.Next() {
		var c BalanceChange
		if err := rows.Scan(
			&c.ChainName, &c.UserAddress, &c.EventType,
			&c.Amount, &c.BalanceAfter, &c.BlockNumber,
			&c.EventTime, &c.TxHash,
		); err != nil {
			return nil, err
		}
		changes = append(changes, c)
	}
	return changes, rows.Err()
}

// GetUsersByChain 获取链上所有用户
func GetUsersByChain(chainName string) ([]string, error) {
	rows, err := Query(
		"SELECT DISTINCT user_address FROM user_balances WHERE chain_name = ?",
		chainName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var addr string
		if err := rows.Scan(&addr); err != nil {
			return nil, err
		}
		users = append(users, addr)
	}
	return users, rows.Err()
}
