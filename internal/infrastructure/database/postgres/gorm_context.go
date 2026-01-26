package postgres

import (
	"doan/pkg/config"
	"doan/pkg/logger"
	"doan/pkg/types"

	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger" // Sử dụng logger chuẩn của GORM hoặc adapter
)

var (
	lock     sync.Mutex
	instance *gorm.DB
)

const (
	// Các giá trị mặc định cho cấu hình pool và thời gian
	defaultMaxIdleConns    = 10
	defaultMaxOpenConns    = 100
	defaultConnMaxLifetime = 5 * time.Minute
	defaultSSLMode         = "disable"             // Cân nhắc đổi thành "require" làm mặc định nếu có thể
	healthCheckInterval    = 1 * time.Minute       // Tần suất kiểm tra kết nối DB
	pingTimeout            = 5 * time.Second       // Timeout cho mỗi lần ping DB
	slowQueryThreshold     = 50 * time.Millisecond // Ngưỡng để log câu query chậm
)

// getValidatedDbConfig lấy và xác thực cấu hình DB từ Manager.
// Trả về lỗi nếu cấu hình thiếu hoặc không hợp lệ.
func getValidatedDbConfig(configManager config.Manager) (*types.PSQLConfig, error) {
	dbConfig := &types.PSQLConfig{}
	if err := configManager.UnmarshalKey("database", dbConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
	}

	// --- Validation và Giá trị mặc định ---
	if dbConfig.Host == "" {
		return nil, fmt.Errorf("database config error: host is required")
	}
	if dbConfig.Port == "" {
		return nil, fmt.Errorf("database config error: port is required")
	}
	if dbConfig.User == "" {
		return nil, fmt.Errorf("database config error: user is required")
	}
	// Password có thể rỗng nếu dùng phương thức xác thực khác, nhưng thường là cần
	// if dbConfig.Password == "" {
	//     return nil, fmt.Errorf("database config error: password is required")
	// }
	if dbConfig.Schema == "" {
		return nil, fmt.Errorf("database config error: schema/dbname is required")
	}

	// Pool Sizes (Giả sử kiểu trong struct là int)
	if dbConfig.MaxIdleConns <= 0 {
		dbConfig.MaxIdleConns = defaultMaxIdleConns
	}
	if dbConfig.MaxOpenConns <= 0 {
		dbConfig.MaxOpenConns = defaultMaxOpenConns
	}

	// Connection Lifetime (Giả sử kiểu trong struct là string)
	if dbConfig.ConnMaxLifetime == "" {
		// Gán giá trị string để log ra dễ đọc, sẽ parse sau
		dbConfig.ConnMaxLifetime = defaultConnMaxLifetime.String()
	}

	return dbConfig, nil
}

// GetDBContext khởi tạo và trả về singleton instance của GORM DB.
// Sử dụng double-checked locking để đảm bảo thread-safety.
// appCtx là context gốc của ứng dụng, dùng để quản lý lifecycle của goroutine kiểm tra sức khỏe.
// Trả về lỗi nếu có vấn đề xảy ra trong quá trình khởi tạo.
func GetDBContext(appCtx context.Context, log logger.Logger, configManager config.Manager) (*gorm.DB, error) {
	// --- Double-Checked Locking: Kiểm tra nhanh bên ngoài lock ---
	if instance != nil {
		return instance, nil
	}

	lock.Lock()
	defer lock.Unlock()

	// --- Double-Checked Locking: Kiểm tra lại bên trong lock ---
	if instance == nil {
		log.Info(appCtx, "Initializing database connection...")

		// 1. Lấy và xác thực cấu hình
		dbConfig, err := getValidatedDbConfig(configManager)
		if err != nil {
			// Lỗi đã được định dạng trong getValidatedDbConfig
			log.Error(appCtx, "Invalid database configuration", "error", err)
			return nil, err
		}
		log.Info(appCtx, "Database configuration loaded", "host", dbConfig.Host, "port", dbConfig.Port, "user", dbConfig.User, "schema", dbConfig.Schema)

		// 2. Tạo chuỗi DSN
		dbSource := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password, // Password lấy trực tiếp từ config
			dbConfig.Schema,
		)

		logLevel := gormlogger.Warn // Mặc định log lỗi và cảnh báo
		isDebug := configManager.GetBool("is_debug")
		if isDebug {
			logLevel = gormlogger.Info // Log cả các câu SQL
		}
		gormLog := NewGormLogger(log, slowQueryThreshold, logLevel, true)

		gormConfig := &gorm.Config{
			Logger: gormLog,
			// Các cấu hình GORM khác nếu cần (NamingStrategy, etc.)
		}

		// 4. Mở kết nối GORM
		db, err := gorm.Open(postgres.Open(dbSource), gormConfig)
		if err != nil {
			log.Error(appCtx, "Failed to open database connection", "error", err)
			return nil, fmt.Errorf("gorm.Open failed: %w", err)
		}

		if isDebug {
			// Không cần db.Debug() nữa nếu logger GORM đã ở level Info
			log.Info(appCtx, "GORM debug mode enabled (SQL logging via GORM logger)")
		}

		// 5. Lấy *sql.DB gốc để cấu hình pool
		sqlDB, err := db.DB()
		if err != nil {
			log.Error(appCtx, "Failed to get underlying *sql.DB from GORM", "error", err)
			// Không có cách nào dễ dàng để đóng db GORM vừa mở ở đây.
			// Lỗi này rất hiếm khi xảy ra.
			return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
		}

		// 6. Cấu hình Connection Pool
		connMaxLifetime, err := time.ParseDuration(dbConfig.ConnMaxLifetime)
		if err != nil {
			log.Warn(appCtx, "Invalid ConnMaxLifetime format in config, using default",
				"configValue", dbConfig.ConnMaxLifetime,
				"defaultValue", defaultConnMaxLifetime,
				"error", err)
			connMaxLifetime = defaultConnMaxLifetime
		}
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(connMaxLifetime)

		log.Info(appCtx, "Database connection pool configured",
			"MaxIdleConns", dbConfig.MaxIdleConns,
			"MaxOpenConns", dbConfig.MaxOpenConns,
			"ConnMaxLifetime", connMaxLifetime,
		)

		// 7. Khởi chạy Goroutine kiểm tra sức khỏe (chỉ chạy một lần)
		// Goroutine này sẽ tự động dừng khi appCtx bị cancel.
		go runDBHealthCheck(appCtx, sqlDB, log)

		log.Info(appCtx, "Database connection initialized successfully.")
		instance = db // Gán instance thành công
	}

	return instance, nil
}

// runDBHealthCheck chạy goroutine nền để kiểm tra kết nối DB định kỳ.
// Goroutine sẽ dừng khi context appCtx bị hủy (cancelled).
func runDBHealthCheck(appCtx context.Context, sqlDB *sql.DB, log logger.Logger) {
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	log.Info(appCtx, "Starting periodic database health check goroutine", "interval", healthCheckInterval)

	for {
		select {
		case <-appCtx.Done(): // Lắng nghe tín hiệu shutdown từ context gốc
			log.Info(appCtx, "Stopping database health check goroutine due to context cancellation.")
			return // Thoát khỏi goroutine một cách an toàn
		case <-ticker.C:
			// Tạo context con với timeout riêng cho việc ping
			pingCtx, cancel := context.WithTimeout(appCtx, pingTimeout)

			err := sqlDB.PingContext(pingCtx)
			cancel() // Luôn gọi cancel sau khi dùng xong context con

			if err != nil {
				// Chỉ log lỗi, không panic.
				// Có thể tăng mức độ log (Error) nếu Ping lỗi liên tục.
				log.Warn(appCtx, "Periodic database ping failed", "error", err)
				// Không log stats nếu ping lỗi
				continue
			}

			// Log pool stats nếu ping thành công (có thể dùng Debug level)
			stats := sqlDB.Stats()
			log.Debug(appCtx, "Database connection pool stats",
				"MaxOpen", stats.MaxOpenConnections, // Tổng số kết nối tối đa được phép
				"Open", stats.OpenConnections, // Số kết nối đang mở (idle + in-use)
				"InUse", stats.InUse, // Số kết nối đang được sử dụng
				"Idle", stats.Idle, // Số kết nối đang chờ (idle)
				"WaitCount", stats.WaitCount, // Tổng số lần phải chờ kết nối
				"WaitDuration", stats.WaitDuration, // Tổng thời gian chờ kết nối
				"MaxIdleClosed", stats.MaxIdleClosed, // Số kết nối bị đóng do vượt MaxIdleConns
				"MaxLifetimeClosed", stats.MaxLifetimeClosed, // Số kết nối bị đóng do vượt ConnMaxLifetime
			)
		}
	}
}
