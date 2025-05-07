package main

import (
	"net/http"

	"github.com/viniciusfonseca/raikiri-wasi-sdk-go/pkg/raikiri"
)

func init() {
	raikiri.Handle(func(w http.ResponseWriter, r *http.Request) {

		connectionSetup := raikiri.NewSqlConnectionSetup()
		connectionSetup.ConnectionType("postgres")

		conn, err := connectionSetup.Init()

		if err != nil {
			panic(err)
		}

		p := raikiri.SqlParams
		conn.ExecuteSql("INSERT INTO accounts (id, balance) VALUES ($1, $2);", p("1", 0))

		accounts, err := conn.QuerySql("SELECT id, balance FROM accounts;", p())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(accounts)
	})
}

func main() {}
