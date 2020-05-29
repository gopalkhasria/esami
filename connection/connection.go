package connection

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
)

//DB connection for use
var DB *sql.DB

//Connect connection to database
func Connect() {
	DB, err := pgx.Connect(context.Background(), "postgresql://postgres:root@localhost:5432/Bitcoin")
	if err != nil || DB == nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
}
