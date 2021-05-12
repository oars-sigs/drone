-- name: create-table-tpipe-pipeline

CREATE TABLE IF NOT EXISTS tpipe_pipelines (
	pipeline_uuid VARCHAR(40),
	pipeline_name VARCHAR(255),
	pipeline_repo VARCHAR(1024),
	pipeline_slug VARCHAR(1024),
	pipeline_ref VARCHAR(255),
	pipeline_sync INT(2)  DEFAULT 0,
	pipeline_content TEXT,
	pipeline_created INTEGER,
	pipeline_updated INTEGER,
	UNIQUE ( pipeline_uuid ) 
);