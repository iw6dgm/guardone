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


