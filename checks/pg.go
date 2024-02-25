package checks

import (
	"database/sql"
	_ "github.com/lib/pq"
	"monitoring/spec"
	"strconv"
)

func checkPg(url string, checkSpec *spec.CheckSpec) error {

	db, err := sql.Open("postgres", url)
	defer db.Close()

	if nil != err {
		return spec.ConnectionFailed
	}

	if nil != db.Ping() {
		return spec.ConnectionFailed
	}

	err = handlePgCheckMethods(checkSpec, db)

	return err
}

func handlePgCheckMethods(checkSpec *spec.CheckSpec, db *sql.DB) error {
	switch checkSpec.Method {
	case spec.CONNECT:
		return nil

	case spec.CONTAINS:
		return spec.NotImplemented

	case spec.STATUS:
		return spec.NotImplemented

	case spec.QUERY:
		var result int
		row := db.QueryRow("SELECT 1+1")
		err := row.Scan(&result)
		if nil != err {
			return spec.QueryFailed
		}

		if nil == checkSpec.Expect {
			return spec.ExpectUndefined
		}

		if strconv.Itoa(result) != *checkSpec.Expect {
			return spec.ExpectFailed
		}
		return nil

	}

	return spec.NoCheckPerformed
}

func (c *CheckHandler) CheckPG(spec *spec.CheckSpec) (*map[string]string, error) {
	return nil, checkPg(spec.DSN, spec)
}
