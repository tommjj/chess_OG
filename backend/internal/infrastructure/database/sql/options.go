package sql

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// createEnumType creates an enum type in the database.
func createEnumType(db *gorm.DB, name string, values []string) error {
	var quotedValues []string
	for _, v := range values {
		quotedValues = append(quotedValues, fmt.Sprintf("'%s'", v))
	}
	valueList := strings.Join(quotedValues, ", ")

	sql := fmt.Sprintf(`
		DO $$
			BEGIN
    			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
        		CREATE TYPE %s AS ENUM (%s);
    		END IF;
		END
		$$;`,
		name, name, valueList)

	if err := db.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

// WithEnumType is a helper function to create enum types if they don't exist.
func WithEnumType(name string, values []string) DBOptionFunc {
	return func(db *gorm.DB) error {
		return createEnumType(db, name, values)
	}
}

// createExtensionsNX is a helper function to create extensions if they don't exist.
func createExtensionsNX(db *gorm.DB, names ...string) error {
	for _, name := range names {
		if tx := db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\";", name)); tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

// WithExtensions is a helper function to create extensions if they don't exist.
func WithExtensions(names ...string) DBOptionFunc {
	return func(db *gorm.DB) error {
		return createExtensionsNX(db, names...)
	}
}

// WithSchema is a helper function to create tables if they don't exist.
func WithSchema(dst ...any) DBOptionFunc {
	return func(db *gorm.DB) error {
		return db.AutoMigrate(dst...)
	}
}

// WithDebug enables GORM debug mode to log SQL statements.
func WithDebug() DBOptionFunc {
	return func(db *gorm.DB) error {
		*db = *db.Debug() // enable debug mode
		return nil
	}
}

// WithLogger sets the logger for GORM.
func WithLogger(l logger.Interface) DBOptionFunc {
	return func(db *gorm.DB) error {
		db.Config.Logger = l
		return nil
	}
}

// WithLoggerLevel sets the logger level for GORM.
func WithLoggerLevel(level logger.LogLevel) DBOptionFunc {
	return func(db *gorm.DB) error {
		db.Config.Logger = db.Config.Logger.LogMode(level)
		return nil
	}
}

// WithIndex create table index
func WithIndex(tableName string, indexName string, expressions ...string) DBOptionFunc {
	return func(db *gorm.DB) error {
		query := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS "%s" ON "%s" (%s)`,
			indexName,
			tableName,
			strings.Join(expressions, ", "),
		)
		return db.Exec(query).Error
	}
}

// WithScript run database script
func WithScript(script string) DBOptionFunc {
	return func(db *gorm.DB) error {
		return db.Exec(script).Error
	}
}
