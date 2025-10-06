package db

const SCHEMA = `
	CREATE TABLE IF NOT EXISTS flights(
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  flight_no  TEXT NOT NULL,
  dep_date   TEXT NOT NULL,
  UNIQUE(flight_no, dep_date)
);

CREATE TABLE IF NOT EXISTS seats(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  flight_id    INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
  label        TEXT NOT NULL, -- e.g., 12A
  cabin        TEXT NOT NULL CHECK (cabin IN ('ECONOMY','BUSINESS','FIRST')),
  is_assigned  INTEGER NOT NULL DEFAULT 0,
  UNIQUE(flight_id, label)
);

CREATE TABLE IF NOT EXISTS vouchers(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  code         TEXT NOT NULL UNIQUE,
  flight_id    INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
  cabin        TEXT NOT NULL CHECK (cabin IN ('ECONOMY','BUSINESS','FIRST')),
  redeemed     INTEGER NOT NULL DEFAULT 0,
  expires_at   TEXT,
  redeemed_at  TEXT
);

CREATE TABLE IF NOT EXISTS seat_assignments(
  voucher_id   INTEGER NOT NULL REFERENCES vouchers(id) ON DELETE CASCADE,
  seat_id      INTEGER NOT NULL REFERENCES seats(id) ON DELETE CASCADE,
  assigned_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  PRIMARY KEY (voucher_id),
  UNIQUE (seat_id)
);

CREATE INDEX IF NOT EXISTS idx_seats_flight_cabin ON seats(flight_id, cabin);`
