package db

import (
	"database/sql"
	"time"
)

// PointsCalculation 积分计算结果
type PointsCalculation struct {
	ChainName    string
	UserAddress  string
	PeriodStart  time.Time
	PeriodEnd    time.Time
	PointsAdded  float64
	TotalPoints  float64
	CalculatedAt time.Time
}

// GetUserLastCalculatedTime 获取用户上次积分计算时间
func GetUserLastCalculatedTime(chainName, userAddr string) (time.Time, error) {
	var lastTime time.Time
	err := QueryRow(`
        SELECT last_calculated_at FROM user_points 
        WHERE chain_name = ? AND user_address = ?
    `, chainName, userAddr).Scan(&lastTime)

	if err == sql.ErrNoRows {
		// 首次计算，返回创建时间
		return time.Now().Add(-24 * time.Hour), nil
	}
	return lastTime, err
}

// UpdateUserPoints 更新用户积分
func UpdateUserPoints(calc PointsCalculation) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 更新总积分
	_, err = TxExec(tx, `
        INSERT INTO user_points (chain_name, user_address, total_points, last_calculated_at)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE 
            total_points = VALUES(total_points),
            last_calculated_at = VALUES(last_calculated_at),
            updated_at = CURRENT_TIMESTAMP
    `,
		calc.ChainName, calc.UserAddress, calc.TotalPoints, calc.PeriodEnd,
	)
	if err != nil {
		return err
	}

	// 记录计算历史
	_, err = TxExec(tx, `
        INSERT INTO points_calculation_history (
            chain_name, user_address, period_start, period_end,
            points_added, total_points
        ) VALUES (?, ?, ?, ?, ?, ?)
    `,
		calc.ChainName, calc.UserAddress, calc.PeriodStart, calc.PeriodEnd,
		calc.PointsAdded, calc.TotalPoints,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetUserTotalPoints 获取用户总积分
func GetUserTotalPoints(chainName, userAddr string) (float64, error) {
	var total float64
	err := QueryRow(`
        SELECT total_points FROM user_points 
        WHERE chain_name = ? AND user_address = ?
    `, chainName, userAddr).Scan(&total)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	return total, err
}

// HasPointsCalculated 检查指定时间段是否已经计算过积分
func HasPointsCalculated(chainName, userAddr string, start, end time.Time) (bool, error) {
	var count int
	err := QueryRow(`
        SELECT COUNT(*) FROM points_calculation_history 
        WHERE chain_name = ? AND user_address = ? 
        AND period_start <= ? AND period_end >= ?
    `, chainName, userAddr, end, start).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetMissingCalculationPeriods 获取缺失的积分计算时间段
func GetMissingCalculationPeriods(chainName, userAddr string, start, end time.Time, intervalMinutes int) ([]TimePeriod, error) {
	var periods []TimePeriod

	// 获取已计算的时间段
	rows, err := Query(`
        SELECT period_start, period_end FROM points_calculation_history 
        WHERE chain_name = ? AND user_address = ? 
        AND period_start >= ? AND period_end <= ?
        ORDER BY period_start
    `, chainName, userAddr, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calculatedPeriods []TimePeriod
	for rows.Next() {
		var p TimePeriod
		if err := rows.Scan(&p.Start, &p.End); err != nil {
			return nil, err
		}
		calculatedPeriods = append(calculatedPeriods, p)
	}

	// 找出缺失的时间段
	current := start
	interval := time.Duration(intervalMinutes) * time.Minute

	for current.Before(end) {
		periodEnd := current.Add(interval)
		if periodEnd.After(end) {
			periodEnd = end
		}

		// 检查这个时间段是否已计算
		hasCalculated := false
		for _, calcPeriod := range calculatedPeriods {
			if !current.Before(calcPeriod.End) && !periodEnd.After(calcPeriod.Start) {
				hasCalculated = true
				break
			}
		}

		if !hasCalculated {
			periods = append(periods, TimePeriod{
				Start: current,
				End:   periodEnd,
			})
		}

		current = periodEnd
	}

	return periods, nil
}

// TimePeriod 时间段结构
type TimePeriod struct {
	Start time.Time
	End   time.Time
}
