package mssql

import (
	"database/sql"
	"errors"
	"log"
	"modelhelper/cli/modelhelper"

	mssql "github.com/denisenkom/go-mssqldb"
)

type msSqlDatabaseService struct {
	connectionService modelhelper.ConnectionService
	sourceFactory     modelhelper.SourceFactoryService
}

// BulkCopy implements modelhelper.DatabaseService.
func (m *msSqlDatabaseService) BulkCopy(source string, dest string, sourceQuery string, table string) (int, error) {
	if source == "" || dest == "" {
		return -1, errors.New("source and/or destination connection not defined")
	}

	cons, err := m.connectionService.Connections()
	if err != nil {
		return -1, err
	}

	if conSource, ok := cons[source]; !ok {
		return -1, errors.New("source connection not found")
	} else if conSource.Type != "mssql" {
		return -1, errors.New("source connection not of type mssql")
	}

	if conDst, ok := cons[dest]; !ok {
		return -1, errors.New("destination connection not found")
	} else if conDst.Type != "mssql" {
		return -1, errors.New("destination connection not of type mssql")
	}

	// srcSqlServer := NewMsSqlSource(m.connectionService, source)

	src, _ := m.sourceFactory.NewSource("mssql", source)

	// t, _ := src.EntitiesFromNames([]string{table})
	t, err := src.Entity(table)
	if err != nil {
		return -1, err
	}
	if t == nil {
		return -1, errors.New("table not found")
	}

	dst, _ := m.sourceFactory.NewSource("mssql", dest)
	_, err = dst.Entity(table)
	if err != nil {
		return -1, errors.New("destination table not found")
	}

	destinationDbConnection, err := sql.Open("mssql", cons[dest].ConnectionString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer destinationDbConnection.Close()

	sourceDbTransaction, err := destinationDbConnection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	cols := []string{}

	for _, col := range t.Columns {
		cols = append(cols, col.Name)
	}

	stmt, err := sourceDbTransaction.Prepare(mssql.CopyIn(table, mssql.BulkOptions{}, cols...))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get data from source
	sourceDbConnection, err := sql.Open("mssql", cons[source].ConnectionString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer sourceDbConnection.Close()

	rows, err := sourceDbConnection.Query(sourceQuery)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var values []interface{}
		for range cols {
			var value interface{}
			values = append(values, &value)
		}

		err = rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			log.Fatal(err)
		}
	}

	result, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = sourceDbTransaction.Commit()
	if err != nil {
		log.Fatal(err)
	}
	rowCount, err := result.RowsAffected()
	return int(rowCount), err

}

func NewMsSqlDatabaseService(connectionService modelhelper.ConnectionService, sourceFactory modelhelper.SourceFactoryService) modelhelper.DatabaseService {
	return &msSqlDatabaseService{connectionService: connectionService, sourceFactory: sourceFactory}
}
