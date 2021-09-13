package config

const AppName = "mysql-resources-db-go-service"

var AppVersion string

type Config struct {
	Port       int  `mapstructure:"server_port" default:"8080"`
	DebugPProf bool `mapstructure:"debug_pprof" default:"false"`

	MySQLDBAddress            string `mapstructure:"mysql_db_address" validate:"required"`
	MySQLDBPort               int    `mapstructure:"mysql_db_port" default:"3306"`
	MySQLDBUser               string `mapstructure:"mysql_db_user" validate:"required"`
	MySQLDBPassword           string `mapstructure:"mysql_db_password" validate:"required"`
	MySQLDBName               string `mapstructure:"mysql_db_name" default:"resource_database"`
	MySQLDBMigrationDirectory string `mapstructure:"mysql_db_migration_dir" validate:"required"`
}
