CREATE TABLE IF NOT EXISTS `hosts` (
 `host_id` INTEGER PRIMARY KEY AUTOINCREMENT,
 `hostname` TEXT KEY,
 `addr` TEXT NOT NULL,
 `addrtype` TEXT NOT NULL,
 `updated_at` TEXT NOT NULL,
 UNIQUE(hostname,addr),
 FOREIGN KEY (host_id) REFERENCES ports (host_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `ports` (
  `host_id` INTEGER NOT NULL,
 `protocol` TEXT NOT NULL,
 `port_id` TEXT NOT NULL,
 `state` TEXT NOT NULL,
 `reason` TEXT NOT NULL,
 `name` TEXT NOT NULL,
 `start_time` TEXT NOT NULL,
 UNIQUE(host_id,port_id,start_time)
);