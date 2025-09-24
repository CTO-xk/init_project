-- 更新现有表结构以支持更大的积分值
ALTER TABLE user_points MODIFY COLUMN total_points DECIMAL(30,6) NOT NULL DEFAULT 0;
ALTER TABLE points_calculation_history MODIFY COLUMN points_added DECIMAL(30,6) NOT NULL;
ALTER TABLE points_calculation_history MODIFY COLUMN total_points DECIMAL(30,6) NOT NULL;
