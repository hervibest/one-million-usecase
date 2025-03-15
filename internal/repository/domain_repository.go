package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DomainRepository interface {
	Insert(ctx context.Context, dataHeaders []string, values []interface{}) error
}

type domainRepository struct {
	db *pgxpool.Pool
}

func NewDomainRepository(db *pgxpool.Pool) DomainRepository {
	return &domainRepository{
		db: db,
	}
}

func (r *domainRepository) Insert(ctx context.Context, dataHeaders []string, values []interface{}) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release() // Pastikan koneksi dikembalikan ke pool

	query := fmt.Sprintf("INSERT INTO domain (%s) VALUES (%s)",
		strings.Join(dataHeaders, ","),
		strings.Join(generateQuestionsMark(len(dataHeaders)), ","),
	)

	_, err = conn.Exec(ctx, query, values...)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}
