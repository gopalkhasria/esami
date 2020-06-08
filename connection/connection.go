package connection

import (
	"database/sql"
	"fmt"

	//using for connection
	_ "github.com/lib/pq"
)

/*const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "Bitcoin"
)*/

//Db exported for use in another file
var Db *sql.DB

//Connect connect to the databse
func Connect() {
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)*/
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=root dbname=Bitcoin sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlStatement := `CREATE TABLE IF NOT EXISTS users( id serial PRIMARY KEY, name VARCHAR(50) NOT NULL, email VARCHAR(50) NOT NULL UNIQUE, pass VARCHAR(70) NOT NULL)`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS keys( id serial PRIMARY KEY, private_key TEXT NOT NULL UNIQUE,public_key TEXT NOT NULL UNIQUE, user_id integer REFERENCES users(id))`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS transactions( id serial PRIMARY KEY, hash TEXT NOT NULL,
		sender TEXT NOT NULL,sign TEXT NOT NULL, block TEXT);`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS outputs( id serial PRIMARY KEY, parent int REFERENCES transactions(id), pkscript TEXT NOT NULL,
		amount TEXT NOT NULL, used BOOLEAN);`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS inputs( id serial PRIMARY KEY, transaction int REFERENCES transactions(id),
		keyHash TEXT NOT NULL, sign TEXT NOT NULL, output int REFERENCES outputs(id));`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS block( id serial PRIMARY KEY, hash TEXT NOT NULL, nounce int NOT NULL, previousHash TEXT NOT NULL, created TIMESTAMP DEFAULT now())`
	db.QueryRow(sqlStatement)
	//sqlStatement = `TRUNCATE block CASCADE`
	//db.QueryRow(sqlStatement)
	sqlStatement = `INSERT INTO block (hash, nounce, previousHash) VALUES ('0', 0, '0')`
	db.QueryRow(sqlStatement)
	Db = db
	fmt.Println("Successfully connected!")
}
