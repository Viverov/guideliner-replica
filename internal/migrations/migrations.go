package migrations

import "github.com/go-gormigrate/gormigrate/v2"

func GetMigrationsList() []*gormigrate.Migration {
	var migrationsList []*gormigrate.Migration

	migrationsList = append(migrationsList, createUserTable())

	return migrationsList
}
