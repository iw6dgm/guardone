CREATE TABLE reading (
	created_at INTEGER(4) NOT NULL DEFAULT (strftime('%s','now')),
	ssid text NOT NULL,
	rssi INTEGER NOT NULL,
	primary key (created_at, ssid)
);

CREATE TABLE network (
	ssid TEXT NOT NULL PRIMARY KEY,
	full_name TEXT NOT NULL
);

CREATE VIEW report_v1 AS
SELECT
DATETIME(created_at,'unixepoch') AS created_at_datetime,r.ssid,IFNULL(n.full_name,'UNDEFINED') AS network_name,r.rssi
FROM reading r
LEFT JOIN network n ON r.ssid = n.ssid
ORDER BY created_at ASC;

