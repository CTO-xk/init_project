package db

import (
	"database/sql"
	"erc20-service/config"
	"erc20-service/pkg/logger"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB  *sql.DB
	log = logger.New("database")
)

// Init 初始化数据库连接
func Init(cfg config.DatabaseConfig) error {
	// MySQL DSN: user:password@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	// 验证连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("验证连接失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(1 * time.Hour)

	DB = db
	log.Info("数据库连接成功")
	return nil
}

// InitChainStatus 初始化链状态
func InitChainStatus(chains []config.ChainConfig) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, chain := range chains {
		var exists bool
		err := TxQueryRow(tx,
			"SELECT EXISTS(SELECT 1 FROM chain_status WHERE chain_name = ?)",
			chain.Name,
		).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			_, err := TxExec(tx,
				"INSERT INTO chain_status (chain_name, last_processed_block) VALUES (?, ?)",
				chain.Name, chain.StartBlock,
			)
			if err != nil {
				return err
			}
			log.Info("初始化链状态", "chain", chain.Name, "start_block", chain.StartBlock)
		}
	}

	return tx.Commit()
}

// ---- SQL Logging Helpers ----

// Exec 执行非查询语句并记录SQL与参数
func Exec(query string, args ...any) (sql.Result, error) {
	log.Info("SQL Exec", "query", query, "args", args)
	return DB.Exec(query, args...)
}

// Query 执行查询并记录SQL与参数
func Query(query string, args ...any) (*sql.Rows, error) {
	log.Info("SQL Query", "query", query, "args", args)
	return DB.Query(query, args...)
}

// QueryRow 执行单行查询并记录SQL与参数
func QueryRow(query string, args ...any) *sql.Row {
	log.Info("SQL QueryRow", "query", query, "args", args)
	return DB.QueryRow(query, args...)
}

// TxExec 在事务中执行非查询语句并记录SQL与参数
func TxExec(tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	log.Info("SQL TxExec", "query", query, "args", args)
	return tx.Exec(query, args...)
}

// TxQuery 在事务中执行查询并记录SQL与参数
func TxQuery(tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	log.Info("SQL TxQuery", "query", query, "args", args)
	return tx.Query(query, args...)
}

// TxQueryRow 在事务中执行单行查询并记录SQL与参数
func TxQueryRow(tx *sql.Tx, query string, args ...any) *sql.Row {
	log.Info("SQL TxQueryRow", "query", query, "args", args)
	return tx.QueryRow(query, args...)
}
