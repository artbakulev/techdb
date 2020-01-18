package infrastructure

import "github.com/jackc/pgx"

var config = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "techdb",
	User:     "techdb_user",
	Password: "techdb_password",
}

var Connection *pgx.ConnPool

func InitDatabaseConnection() error {
	var err error
	Connection, err = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 50,
		})
	return err
}
