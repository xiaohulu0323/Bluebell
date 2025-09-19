package mysql

import (
	"fmt"
	"web-app/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

// Init 初始化MySQL连接（兼容原有代码）
func Init(cfg *settings.MySQLConfig) (err error) {
	// 优先使用增强版连接池
	if err = InitAdvanced(cfg); err == nil {
		zap.L().Info("使用增强版数据库连接池")
		return nil
	}

	// 降级到原始连接池
	zap.L().Warn("增强版连接池初始化失败，降级到原始连接池", zap.Error(err))

	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	zap.L().Info("原始数据库连接池初始化成功")
	return
}

// Close 关闭MySQL连接
func Close() {
	// 优先关闭增强版连接池
	if dbManager != nil {
		if err := dbManager.Close(); err != nil {
			zap.L().Error("关闭增强版连接池失败", zap.Error(err))
		}
		return
	}

	// 降级关闭原始连接
	if db != nil {
		_ = db.Close()
	}
}
