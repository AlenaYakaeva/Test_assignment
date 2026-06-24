package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	conn *pgxpool.Pool
}

func New(dbDSN string) (*Storage, error) {
	conn, err := pgxpool.New(context.TODO(), dbDSN)
	if err != nil {
		return &Storage{}, err
	}
	return &Storage{
		conn: conn,
	}, nil
}

func (s *Storage) Close() {
	s.conn.Close()
}

func RunMigrations(dbDSN string) error {
	m, err := migrate.New("file://migration", dbDSN)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no change")
		} else {
			return err
		}
	}
	fmt.Println("migrations complete")
	return nil
}
