package utils

import (
	"errors"
	"net/http"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// HandlePGError identifica erros específicos do PostgreSQL e retorna uma mensagem amigável e o status HTTP adequado
func HandlePGError(err error) (string, int) {
	var pgErr *pq.Error
	if !errors.As(err, &pgErr) {
		return "Erro interno no banco de dados", http.StatusInternalServerError
	}

	// Log estruturado para facilitar a depuração
	log.WithFields(log.Fields{
		"constraint": pgErr.Constraint,
		"code":       pgErr.Code,
		"message":    pgErr.Message,
		"detail":     pgErr.Detail,
	}).Error("Postgres Driver Error")

	// SQLSTATE codes: https://www.postgresql.org/docs/current/errcodes-appendix.html

	switch pgErr.Code {
	case "23505": // unique_violation
		if msg, found := LookupPGMessage(pgErr.Detail, UniqueViolations); found {
			return msg, http.StatusConflict
		}
		if msg, found := LookupPGMessage(pgErr.Constraint, UniqueViolations); found {
			return msg, http.StatusConflict
		}
		return "Já existe um registro com estes dados.", http.StatusConflict

	case "23503": // foreign_key_violation
		// Prefer match by constraint symbol (ent geralmente popula Constraint), depois por Detail
		if msg, found := LookupPGMessage(pgErr.Constraint, ForeignKeyViolations); found {
			return msg, http.StatusNotFound
		}
		if msg, found := LookupPGMessage(pgErr.Detail, ForeignKeyViolations); found {
			return msg, http.StatusNotFound
		}
		return "Um registro relacionado não foi encontrado.", http.StatusNotFound

	default:
		return "Erro na operação de banco de dados.", http.StatusInternalServerError
	}
}
