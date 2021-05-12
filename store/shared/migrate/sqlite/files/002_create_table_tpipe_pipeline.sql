-- name: create-table-tpipe-pipeline

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