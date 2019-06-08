CREATE TABLE domains (
	name STRING NOT NULL,
	servers_changed STRING NULL,
	ssl_grade STRING NULL,
	previous_ssl_grade STRING NULL,
	logo STRING NULL,
	title STRING NULL,
	is_down STRING NULL,
	review_time DATE NULL,
	CONSTRAINT "primary" PRIMARY KEY (name ASC),
	UNIQUE INDEX domains_servers_changed_key (servers_changed ASC),
	FAMILY "primary" (name, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down, review_time)
)

CREATE TABLE servers (
	address STRING NOT NULL,
	ssl_grade STRING NULL,
	country STRING NULL,
	owner STRING NULL,
	domain_name STRING NULL,
	CONSTRAINT "primary" PRIMARY KEY (address ASC),
	UNIQUE INDEX servers_ssl_grade_key (ssl_grade ASC),
	CONSTRAINT fk_domain_name_ref_domains FOREIGN KEY (domain_name) REFERENCES domains (name) ON DELETE CASCADE,
	INDEX servers_auto_index_fk_domain_name_ref_domains (domain_name ASC),
	FAMILY "primary" (address, ssl_grade, country, owner, domain_name)
)