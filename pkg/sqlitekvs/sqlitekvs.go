package sqlitekvs

import (
	"database/sql"
	"fmt"
	"time"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// KVSCache :
type KVSCache struct {
	FilePath string
	TTL      time.Duration

	db *sql.DB
}

func (kvs *KVSCache) ensureDbOpen() error {

	if kvs.db != nil {
		return nil
	}

	driverName := "sqlite3"
	dataSourceName := kvs.FilePath
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("colud not open database '%s' with driver '%s': %w", dataSourceName, driverName, err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS kvs (key TEXT PRIMARY KEY, value TEXT NOT NULL, expires DATETIME NOT NULL);")
	if err != nil {
		return fmt.Errorf("colud not initialize database schema: %w", err)
	}

	kvs.db = db

	return nil
}

// Get :
func (kvs *KVSCache) Get(key string) (string, bool, error) {

	kvs.ensureDbOpen()

	row := kvs.db.QueryRow(`SELECT value FROM kvs WHERE key = ? AND expires > DATETIME('now');`, key)

	var content string
	err := row.Scan(&content)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return content, true, nil
}

// Set :
func (kvs *KVSCache) Set(key, value string) error {

	kvs.ensureDbOpen()

	ttl := kvs.TTL.Minutes()
	expiration := fmt.Sprintf(`DATETIME('now', '+%d minutes')`, int(ttl))

	_, err := kvs.db.Exec(fmt.Sprintf(`
		INSERT INTO kvs(key, value, expires)
		VALUES(?, ?, %s)
		ON CONFLICT(key) 
			DO UPDATE SET
				value = ?,
				expires = %s`, expiration, expiration), key, value, value)
	if err != nil {
		return fmt.Errorf("could not save item: %w", err)
	}

	return nil
}
