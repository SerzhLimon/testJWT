package postgres

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/SerzhLimon/testJWT/config"
)

func InitPostgresClient(cfg config.PostgresConfig) (*sql.DB, error) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	options := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode)

	database, err := sql.Open("postgres", options)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":    cfg.Host,
			"port":    cfg.Port,
			"user":    cfg.User,
			"dbname":  cfg.DBName,
			"sslmode": cfg.SSLMode,
			"error":   err.Error(),
		}).Error("Failed to open PostgreSQL connection")
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host":    cfg.Host,
			"port":    cfg.Port,
			"user":    cfg.User,
			"dbname":  cfg.DBName,
			"sslmode": cfg.SSLMode,
			"error":   err.Error(),
		}).Error("Failed to ping PostgreSQL database")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"host":    cfg.Host,
		"port":    cfg.Port,
		"user":    cfg.User,
		"dbname":  cfg.DBName,
		"sslmode": cfg.SSLMode,
	}).Info("Successful connection to PostgreSQL")

	return database, nil
}
