-- name: create-table-tpipe-template

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