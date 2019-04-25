# go-myquery
A wrapper around the Golang MySQL driver that can fetch query results as maps of column names to string values
```
package myquery // import "github.com/greenlion/go-myquery/myquery"

func Connect(host string, user string, password string, port string, db string) (*sql.DB, error)
func Fetch(rows *sql.Rows) (map[string]string, error)
func Query(db *sql.DB, sql string) (*sql.Rows, error)
```

## Example usage
```
import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/greenlion/go-myquery/myquery"
)

func main() {

  /* connection to the database */
  conn, err := myquery.Connect("127.0.0.1", "root", "", "3306", "performance_schema")
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  /* run a query and get a pointer to the resultset */
  stmt, err := myquery.Query(conn, "select * from threads")
  if err != nil {
    panic(err)
  }
  /* note that fetching an entire resultset will automatically close the 
     resultset.  You can't fetch the same resultset twice or rewind
     the resultset. deferred Close() is still a good idea in case the
     entire resultset is not fetched.
  */
  defer stmt.Close()

  /* fetch each row from the resultset into a map[string]string and exit when no more rows */
  for {
    row, err := myquery.Fetch(stmt)
    if err != nil {
      panic(err)
    }
    if row == nil {
      break
    }
    /* print the row */
    fmt.Println(row)
  }

}
```
