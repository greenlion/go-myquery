package myquery

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/* Establish a database connection, immediately connecting to the database */
func Connect(host string, user string, password string, port string, db string) (*sql.DB, error) {
	connstr := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	conn, err := sql.Open("mysql", connstr)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Query(db *sql.DB, sql string) (*sql.Rows, error) {
	rows, err := db.Query(sql)
	return rows, err
}

/* This function is generic in that it returns all column values as a
   map of column names to  string values.  NULL values are represented
	 by string :NULL: (this is unlikely to appear in real data)

	 This function is less useful if you need to do calculations on
	 numeric values.  It makes more sense to fetch a SQL resultset
	 into a struct in that case.
*/
func Fetch(rows *sql.Rows) (map[string]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		column_pointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			column_pointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(column_pointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]string)
		for i, column_name := range cols {
			raw := *(column_pointers[i].(*interface{}))
			if raw == nil {
				m[column_name] = ":NULL:"
			} else {
				/* convert the array of unsigned integers into a string */
				tmp := []byte{}
				for _, b := range raw.([]uint8) {
					tmp = append(tmp, byte(b))
				}

				m[column_name] = string(tmp)
			}
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		return m, nil
	}

	return nil, nil
}
