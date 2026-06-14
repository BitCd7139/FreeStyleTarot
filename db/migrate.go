package db

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// RunMigrations 执行嵌入的 SQL migration（与 DATABASE_URL 相同，支持 Supabase Session pooler）
func RunMigrations(databaseURL string) {
	normalized, err := NormalizeDatabaseURL(databaseURL)
	if err != nil {
		zap.S().Fatalw("migration 数据库 URL 无效", "error", err)
	}
	migrateURL := ToMigrateURL(normalized)

	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		zap.S().Fatalw("加载 migration 文件失败", "error", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, migrateURL)
	if err != nil {
		zap.S().Fatalw("初始化 migration 失败", "error", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		zap.S().Fatalw("执行 migration 失败", "error", err)
	}

	zap.S().Infow("数据库 migration 完成")
}
