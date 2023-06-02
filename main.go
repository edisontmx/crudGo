package main

import (
	/* "context" */
	"crudgo/envs"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

/* func conexionDB() (*sql.DB, error) {
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

} */

func conexionDB() (*sql.DB, error) {

	Driver := envs.Get("DB_DRIVER", "mysql")
	if Driver == "" {
		Driver = "mysql"
	}
	Usuario := envs.Get("DB_USER", "root")
	if Usuario == "" {
		Usuario = "root"
	}
	Clave := envs.Get("DB_KEY", "")
	Puerto := envs.Get("DB_PORT", "@tcp(localhost:3308)/")
	if Puerto == "" {
		Puerto = "@tcp(localhost:3308)/"
	}
	Nombre := envs.Get("DB_NAME", "sistema")
	if Nombre == "" {
		Nombre = "sistema"
	}

	//connString := "root:@tcp(localhost:3308)/sistema?parseTime=true"

	db, err := sql.Open(Driver, Usuario+Clave+Puerto /* "@tcp(localhost:3308)/" */ +Nombre)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}

/* func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Clave := "root"
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Clave+"@tcp(localhost:3308)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
} */

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	/* ctx := context.Background() */
	db, err := conexionDB()
	if err != nil {
		panic(err)
	}

	/* err = insert(ctx, db, "abiud Medina", "correoabiud@correo.com") */

	http.HandleFunc("/", index)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)

	db.Close()
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)
	db, err := conexionDB()

	queryDelete, err := db.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		return
	}
	queryDelete.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)

}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hola developer")

	db, err := conexionDB()

	queryAll, err := db.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for queryAll.Next() {
		var id int
		var nombre, correo string
		err = queryAll.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}

	//fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(w, "index", arregloEmpleado)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	db, err := conexionDB()

	queryOne, err := db.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)
	empleado := Empleado{}
	for queryOne.Next() {
		var id int
		var nombre, correo string
		err = queryOne.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		fmt.Println(empleado)
		plantillas.ExecuteTemplate(w, "editar", empleado)
	}
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		db, err := conexionDB()

		queryAdd, err := db.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?,?)")

		if err != nil {
			return
		}
		queryAdd.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)
	}
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		db, err := conexionDB()

		queryUpdate, err := db.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")

		if err != nil {
			return
		}
		queryUpdate.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)
	}
}

/* func insert(ctx context.Context, db *sql.DB, nombre string, correo string) error {

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
*/
