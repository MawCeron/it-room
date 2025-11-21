package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
	}

	for _, p := range pragmas {
		if _, err := conn.Exec(p); err != nil {
			conn.Close()
			return nil, fmt.Errorf("could not set pragma %s: %w", p, err)
		}
	}

	db := &DB{Conn: conn}

	if err := db.applyInitialMigration(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

func (d *DB) Close() error { return d.Conn.Close() }

func (d *DB) applyInitialMigration() error {
	var count int
	row := d.Conn.QueryRow("SELECT count(name) FROM sqlite_master WHERE type='table'")
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	path := "migrations/schema.sql"

	var r io.Reader
	var f *os.File
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		ff, err := os.Open(path)
		if err != nil {
			return err
		}
		f = ff
		defer f.Close()
		r = f
	}

	if r == nil {
		return fmt.Errorf("migration file not found in %v", path)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	stmts := splitSQLStatements(string(b))
	tx, err := d.Conn.Begin()
	if err != nil {
		return err
	}

	for _, s := range stmts {
		if len(s) == 0 {
			continue
		}
		if _, err := tx.Exec(s); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func splitSQLStatements(sqlText string) []string {
	var out []string
	cur := ""
	for _, r := range sqlText {
		cur += string(r)
		if r == ';' {
			out = append(out, cur)
			cur = ""
		}
	}

	if len(cur) > 0 {
		out = append(out, cur)
	}

	return out
}
