-- name: create-table-tpipe-template

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