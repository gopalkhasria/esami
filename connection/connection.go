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
	sqlStatement = `CREATE TABLE IF NOT EXISTS keys( id serial PRIMARY KEY, private_key TEXT NOT NULL UNIQUE, public_key TEXT NOT NULL UNIQUE, user_id integer REFERENCES users(id))`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS status( id serial PRIMARY KEY, value varchar(15));`
	db.QueryRow(sqlStatement)
	//sqlStatement = `INSERT INTO status(value) VALUES ('not confirmed'), ('confirmed'),('rejected');`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS transactions( id serial PRIMARY KEY, hash TEXT NOT NULL UNIQUE,
		sender TEXT NOT NULL UNIQUE,sign TEXT NOT NULL UNIQUE, amount INT, isUsed BOOLEAN, status integer REFERENCES status(id));`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS outputs( id serial PRIMARY KEY, parent int REFERENCES transactions(id), out_transaction int REFERENCES transactions(id),
		condition TEXT NOT NULL);`
	db.QueryRow(sqlStatement)
	sqlStatement = `CREATE TABLE IF NOT EXISTS inputs( id serial PRIMARY KEY, input int REFERENCES transactions(id), out_transaction int REFERENCES transactions(id),
		pkScript TEXT NOT NULL, keyHash TEXT NOT NULL, sign TEXT NOT NULL);`
	db.QueryRow(sqlStatement)
	Db = db
	fmt.Println("Successfully connected!")
}
