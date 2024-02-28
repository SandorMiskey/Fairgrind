-- CREATE DATABASE IF NOT EXISTS {{.DB_SCHEMA}};
-- USE {{.DB_SCHEMA}};

CREATE USER '{{.DB_USER}}'@'%' IDENTIFIED BY '{{.DB_PASSWORD}}';
CREATE USER '{{.DB_USER}}'@'localhost' IDENTIFIED BY '{{.DB_PASSWORD}}';

GRANT ALL ON {{.DB_SCHEMA}}.* TO '{{.DB_USER}}'@'%';
GRANT ALL ON {{.DB_SCHEMA}}.* TO '{{.DB_USER}}'@'localhost';

GRANT SELECT, REPLICATION CLIENT, REPLICATION SLAVE ON *.* TO '{{.DB_USER}}'@'%';
GRANT SELECT, REPLICATION CLIENT, REPLICATION SLAVE ON *.* TO '{{.DB_USER}}'@'localhost';

--
--
--

CREATE USER '{{.DB_USER2}}'@'%' IDENTIFIED BY '{{.DB_PASSWORD2}}';
CREATE USER '{{.DB_USER2}}'@'localhost' IDENTIFIED BY '{{.DB_PASSWORD2}}';

GRANT ALL ON {{.DB_SCHEMA}}.* TO '{{.DB_USER2}}'@'%';
GRANT ALL ON {{.DB_SCHEMA}}.* TO '{{.DB_USER2}}'@'localhost';

GRANT SELECT, REPLICATION CLIENT, REPLICATION SLAVE ON *.* TO '{{.DB_USER2}}'@'%';
GRANT SELECT, REPLICATION CLIENT, REPLICATION SLAVE ON *.* TO '{{.DB_USER2}}'@'localhost';