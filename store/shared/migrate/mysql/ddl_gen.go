package mysql

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
	{
		name: "alter-table-pipelines-add-column-sync",
		stmt: alterTablePipelinesAddColumnSync,
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
	template_uuid VARCHAR(40),
	template_name VARCHAR(255),
	template_format VARCHAR(255),
	template_type VARCHAR(255),
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
	pipeline_uuid VARCHAR(40),
	pipeline_name VARCHAR(255),
	pipeline_repo VARCHAR(1024),
	pipeline_slug VARCHAR(1024),
	pipeline_ref VARCHAR(255),
	pipeline_content MEDIUMTEXT,
	pipeline_created INTEGER,
	pipeline_updated INTEGER,
	UNIQUE ( pipeline_uuid ) 
);
`

var alterTablePipelinesAddColumnSync = `
ALTER TABLE tpipe_pipelines ADD COLUMN pipeline_sync INT(2) DEFAULT 0;
`
