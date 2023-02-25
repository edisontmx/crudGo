package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (*sql.DB, error) {
	connString := "root:@tcp(localhost:3308)/sistema?parseTime=true"

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	ctx := context.Background()
	db, err := conexionDB()
	if err != nil {
		panic(err)
	}

	err = insert(ctx, db, "abiud Medina", "correoabiud@correo.com")

	http.HandleFunc("/", index)
	http.HandleFunc("/crear", Crear)

	log.Println("Servidor corriendo")

	http.ListenAndServe(":8080", nil)

	db.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hola developer")
	plantillas.ExecuteTemplate(w, "index", nil)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func insert(ctx context.Context, db *sql.DB, nombre string, correo string) error {

	queryADD := `INSERT INTO empleados( nombre, correo) values(?, ?)`

	result, err := db.ExecContext(ctx, queryADD, nombre, correo)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Println("inserted id: ", id)

	return nil
}
