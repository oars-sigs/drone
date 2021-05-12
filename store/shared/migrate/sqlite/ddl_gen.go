package sqlite

import (
	"database/sql"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-tpipe-template",
		stmt: createTableTpipeTemplate,
	},
	{
		name: "create-table-tpipe-pipeline",
		stmt: createTableTpipePipeline,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

//
// 001_create_table_tpipe_template.sql
//

var createTableTpipeTemplate = `
CREATE TABLE IF NOT EXISTS tpipe_templates (
	template_uuid TEXT,
	template_name TEXT,
	template_format TEXT,
	template_type TEXT,
    template_content TEXT,
	template_updated INT,
	template_created INT,
	UNIQUE ( template_uuid ),
	UNIQUE ( template_name )
);
`

//
// 002_create_table_tpipe_pipeline.sql
//

var createTableTpipePipeline = `
CREATE TABLE IF NOT EXISTS tpipe_pipelines (
	pipeline_uuid TEXT,
	pipeline_name TEXT,
	pipeline_repo TEXT,
	pipeline_slug TEXT,
	pipeline_ref TEXT,
	pipeline_sync INT(2) DEFAULT 0,
	pipeline_content TEXT,
	pipeline_created INTEGER,
	pipeline_updated INTEGER,
	UNIQUE ( pipeline_uuid ) 
);
`
