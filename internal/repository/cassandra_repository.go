package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gocql/gocql"
)

type CassandraRepository interface {
	Insert(ctx context.Context, dataHeaders []string, values []interface{}) error
}

type cassandraRepository struct {
	session *gocql.Session
}

func NewCassandraRepository(session *gocql.Session) CassandraRepository {
	return &cassandraRepository{
		session: session,
	}
}

func (r *cassandraRepository) Insert(ctx context.Context, dataHeaders []string, values []interface{}) error {
	query := fmt.Sprintf("INSERT INTO domain (%s) VALUES (%s)",
		strings.Join(dataHeaders, ","),
		generateQuestionMarks(len(dataHeaders)),
	)

	if err := r.session.Query(query, values...).WithContext(ctx).Exec(); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func generateQuestionMarks(n int) string {
	marks := make([]string, n)
	for i := range marks {
		marks[i] = "?"
	}
	return strings.Join(marks, ",")
}
