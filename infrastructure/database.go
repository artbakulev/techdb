package infrastructure

import "github.com/jackc/pgx"

var config = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "techdb",
	User:     "techdb_user",
	Password: "techdb_password",
}

func InitDatabaseConnection() (*pgx.ConnPool, error) {
	var err error
	connection, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 100,
		})
	return connection, err
}
