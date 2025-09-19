package mysql

import (
	"context"
	"fmt"
	"sync"
	"time"
	"web-app/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// DBManager 数据库连接管理器
type DBManager struct {
	writeDB  *sqlx.DB   // 写数据库
	readDBs  []*sqlx.DB // 读数据库列表
	rwMutex  sync.RWMutex
	config   *settings.MySQLConfig
	stats    *PoolStats
	statsMux sync.RWMutex
}

// PoolStats 连接池统计信息
type PoolStats struct {
	// 连接池状态
	MaxOpenConnections int `json:"max_open_connections"`
	OpenConnections    int `json:"open_connections"`
	InUse              int `json:"in_use"`
	Idle               int `json:"idle"`

	// 等待统计
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`

	// 读写分离统计
	WriteQueryCount int64   `json:"write_query_count"`
	ReadQueryCount  int64   `json:"read_query_count"`
	ReadWriteRatio  float64 `json:"read_write_ratio"`

	// 性能统计
	AvgQueryTime   float64 `json:"avg_query_time_ms"`
	SlowQueryCount int64   `json:"slow_query_count"`
	ErrorCount     int64   `json:"error_count"`

	LastUpdateTime time.Time `json:"last_update_time"`
}

var (
	dbManager *DBManager
	once      sync.Once
)

// InitAdvanced 初始化增强版MySQL连接池
func InitAdvanced(cfg *settings.MySQLConfig) (err error) {
	once.Do(func() {
		dbManager = &DBManager{
			config: cfg,
			stats: &PoolStats{
				LastUpdateTime: time.Now(),
			},
		}
		err = dbManager.init()
	})
	return err
}

// init 初始化数据库连接
func (dm *DBManager) init() error {
	// 初始化写数据库
	if err := dm.initWriteDB(); err != nil {
		return fmt.Errorf("初始化写数据库失败: %v", err)
	}

	// 初始化读数据库（如果启用读写分离）
	if dm.config.EnableReadWriteSplit && len(dm.config.ReadHosts) > 0 {
		if err := dm.initReadDBs(); err != nil {
			zap.L().Error("初始化读数据库失败，将使用写数据库处理读请求", zap.Error(err))
		}
	}

	// 启动统计信息更新协程
	go dm.updateStatsRoutine()

	zap.L().Info("数据库连接池初始化成功",
		zap.Int("max_open_conns", dm.config.MaxOpenConns),
		zap.Int("max_idle_conns", dm.config.MaxIdleConns),
		zap.Bool("read_write_split", dm.config.EnableReadWriteSplit),
		zap.Int("read_db_count", len(dm.readDBs)))

	return nil
}

// initWriteDB 初始化写数据库连接
func (dm *DBManager) initWriteDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		dm.config.User, dm.config.Password, dm.config.Host, dm.config.Port, dm.config.DB)

	writeDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接写数据库失败: %v", err)
	}

	// 配置连接池参数
	dm.configureDB(writeDB, "write")
	dm.writeDB = writeDB

	return nil
}

// initReadDBs 初始化读数据库连接
func (dm *DBManager) initReadDBs() error {
	for i, host := range dm.config.ReadHosts {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&timeout=10s&readTimeout=30s",
			dm.config.User, dm.config.Password, host, dm.config.Port, dm.config.DB)

		readDB, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			zap.L().Error("连接读数据库失败", zap.String("host", host), zap.Error(err))
			continue
		}

		// 读数据库连接池配置（通常比写数据库更大）
		dm.configureDB(readDB, fmt.Sprintf("read-%d", i))
		dm.readDBs = append(dm.readDBs, readDB)
	}

	if len(dm.readDBs) == 0 {
		return fmt.Errorf("所有读数据库连接均失败")
	}

	return nil
}

// configureDB 配置数据库连接池参数
func (dm *DBManager) configureDB(database *sqlx.DB, dbType string) {
	// 基础连接池配置
	database.SetMaxOpenConns(dm.config.MaxOpenConns)
	database.SetMaxIdleConns(dm.config.MaxIdleConns)

	// 连接生命周期配置
	if dm.config.ConnMaxLifetime > 0 {
		database.SetConnMaxLifetime(time.Duration(dm.config.ConnMaxLifetime) * time.Minute)
	} else {
		database.SetConnMaxLifetime(30 * time.Minute) // 默认30分钟
	}

	if dm.config.ConnMaxIdleTime > 0 {
		database.SetConnMaxIdleTime(time.Duration(dm.config.ConnMaxIdleTime) * time.Minute)
	} else {
		database.SetConnMaxIdleTime(15 * time.Minute) // 默认15分钟
	}

	zap.L().Info("数据库连接池配置完成",
		zap.String("type", dbType),
		zap.Int("max_open_conns", dm.config.MaxOpenConns),
		zap.Int("max_idle_conns", dm.config.MaxIdleConns))
}

// GetWriteDB 获取写数据库连接
func (dm *DBManager) GetWriteDB() *sqlx.DB {
	dm.rwMutex.RLock()
	defer dm.rwMutex.RUnlock()

	dm.statsMux.Lock()
	dm.stats.WriteQueryCount++
	dm.statsMux.Unlock()

	return dm.writeDB
}

// GetReadDB 获取读数据库连接（负载均衡）
func (dm *DBManager) GetReadDB() *sqlx.DB {
	dm.rwMutex.RLock()
	defer dm.rwMutex.RUnlock()

	dm.statsMux.Lock()
	dm.stats.ReadQueryCount++
	dm.statsMux.Unlock()

	// 如果没有配置读数据库，返回写数据库
	if len(dm.readDBs) == 0 {
		return dm.writeDB
	}

	// 简单的轮询负载均衡
	index := time.Now().UnixNano() % int64(len(dm.readDBs))
	return dm.readDBs[index]
}

// RecordQueryTime 记录查询时间
func (dm *DBManager) RecordQueryTime(duration time.Duration, isError bool) {
	dm.statsMux.Lock()
	defer dm.statsMux.Unlock()

	durationMs := float64(duration.Nanoseconds()) / 1e6

	// 更新平均查询时间
	totalQueries := dm.stats.WriteQueryCount + dm.stats.ReadQueryCount
	if totalQueries > 0 {
		dm.stats.AvgQueryTime = (dm.stats.AvgQueryTime*float64(totalQueries-1) + durationMs) / float64(totalQueries)
	}

	// 记录慢查询（超过1秒）
	if duration > time.Second {
		dm.stats.SlowQueryCount++
	}

	// 记录错误
	if isError {
		dm.stats.ErrorCount++
	}
}

// GetStats 获取连接池统计信息
func (dm *DBManager) GetStats() *PoolStats {
	dm.statsMux.Lock()
	defer dm.statsMux.Unlock()

	// 更新连接池状态
	if dm.writeDB != nil {
		dbStats := dm.writeDB.Stats()
		dm.stats.MaxOpenConnections = dbStats.MaxOpenConnections
		dm.stats.OpenConnections = dbStats.OpenConnections
		dm.stats.InUse = dbStats.InUse
		dm.stats.Idle = dbStats.Idle
		dm.stats.WaitCount = dbStats.WaitCount
		dm.stats.WaitDuration = dbStats.WaitDuration
		dm.stats.MaxIdleClosed = dbStats.MaxIdleClosed
		dm.stats.MaxIdleTimeClosed = dbStats.MaxIdleTimeClosed
		dm.stats.MaxLifetimeClosed = dbStats.MaxLifetimeClosed
	}

	// 计算读写比例
	total := dm.stats.WriteQueryCount + dm.stats.ReadQueryCount
	if total > 0 {
		dm.stats.ReadWriteRatio = float64(dm.stats.ReadQueryCount) / float64(total)
	}

	dm.stats.LastUpdateTime = time.Now()

	// 返回副本以避免并发问题
	return &PoolStats{
		MaxOpenConnections: dm.stats.MaxOpenConnections,
		OpenConnections:    dm.stats.OpenConnections,
		InUse:              dm.stats.InUse,
		Idle:               dm.stats.Idle,
		WaitCount:          dm.stats.WaitCount,
		WaitDuration:       dm.stats.WaitDuration,
		MaxIdleClosed:      dm.stats.MaxIdleClosed,
		MaxIdleTimeClosed:  dm.stats.MaxIdleTimeClosed,
		MaxLifetimeClosed:  dm.stats.MaxLifetimeClosed,
		WriteQueryCount:    dm.stats.WriteQueryCount,
		ReadQueryCount:     dm.stats.ReadQueryCount,
		ReadWriteRatio:     dm.stats.ReadWriteRatio,
		AvgQueryTime:       dm.stats.AvgQueryTime,
		SlowQueryCount:     dm.stats.SlowQueryCount,
		ErrorCount:         dm.stats.ErrorCount,
		LastUpdateTime:     dm.stats.LastUpdateTime,
	}
}

// updateStatsRoutine 定期更新统计信息
func (dm *DBManager) updateStatsRoutine() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := dm.GetStats()

		// 记录关键指标到日志
		if stats.OpenConnections > int(float64(stats.MaxOpenConnections)*0.8) {
			zap.L().Warn("数据库连接池使用率较高",
				zap.Int("open_connections", stats.OpenConnections),
				zap.Int("max_open_connections", stats.MaxOpenConnections),
				zap.Float64("usage_rate", float64(stats.OpenConnections)/float64(stats.MaxOpenConnections)))
		}

		if stats.WaitCount > 0 {
			zap.L().Warn("数据库连接池出现等待",
				zap.Int64("wait_count", stats.WaitCount),
				zap.Duration("wait_duration", stats.WaitDuration))
		}
	}
}

// HealthCheck 健康检查
func (dm *DBManager) HealthCheck(ctx context.Context) error {
	// 检查写数据库
	if err := dm.writeDB.PingContext(ctx); err != nil {
		return fmt.Errorf("写数据库健康检查失败: %v", err)
	}

	// 检查读数据库
	for i, readDB := range dm.readDBs {
		if err := readDB.PingContext(ctx); err != nil {
			zap.L().Error("读数据库健康检查失败", zap.Int("index", i), zap.Error(err))
			// 读数据库失败不影响整体健康状态，只记录日志
		}
	}

	return nil
}

// Close 关闭所有数据库连接
func (dm *DBManager) Close() error {
	dm.rwMutex.Lock()
	defer dm.rwMutex.Unlock()

	var errs []error

	// 关闭写数据库
	if dm.writeDB != nil {
		if err := dm.writeDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭写数据库失败: %v", err))
		}
	}

	// 关闭读数据库
	for i, readDB := range dm.readDBs {
		if err := readDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭读数据库%d失败: %v", i, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("关闭数据库连接时发生错误: %v", errs)
	}

	return nil
}

// 全局函数，保持向后兼容
func GetWriteDB() *sqlx.DB {
	if dbManager == nil {
		return db // 降级到原始连接
	}
	return dbManager.GetWriteDB()
}

func GetReadDB() *sqlx.DB {
	if dbManager == nil {
		return db // 降级到原始连接
	}
	return dbManager.GetReadDB()
}

func GetDBStats() *PoolStats {
	if dbManager == nil {
		return nil
	}
	return dbManager.GetStats()
}

func RecordQueryTime(duration time.Duration, isError bool) {
	if dbManager != nil {
		dbManager.RecordQueryTime(duration, isError)
	}
}

func DBHealthCheck(ctx context.Context) error {
	if dbManager == nil {
		return db.PingContext(ctx)
	}
	return dbManager.HealthCheck(ctx)
}
