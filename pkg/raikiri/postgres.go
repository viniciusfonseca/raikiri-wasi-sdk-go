package raikiri

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rajatjindal/wasi-go-sdk/pkg/wasihttp"
)

type PostgresConnection struct {
	connection_id string
}

type PgConnectionBuilder struct {
	connection_string_secret_name string
}

var Handle = wasihttp.Handle

func NewPgConnectionBuilder() *PgConnectionBuilder {

	return &PgConnectionBuilder{}
}

func (builder *PgConnectionBuilder) ConnectionStringSecretName(connection_string_secret_name string) *PgConnectionBuilder {
	builder.connection_string_secret_name = connection_string_secret_name
	return builder
}

func (builder *PgConnectionBuilder) Build() (*PostgresConnection, error) {

	client := wasihttp.NewClient()
	req, err := http.NewRequest(http.MethodGet, "https://raikiri.db/postgres", nil)

	if err != nil {
		return nil, err
	}

	if builder.connection_string_secret_name == "" {
		builder.connection_string_secret_name = "POSTGRES_CONNECTION_STRING"
	}

	req.Header.Set("Connection-String-Secret-Name", builder.connection_string_secret_name)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	connection_id, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &PostgresConnection{connection_id: string(connection_id)}, nil
}

func (connection *PostgresConnection) ExecuteSql(query string, params []interface{}) (int, error) {

	var value struct {
		Query  string
		Params []interface{}
	}

	value.Query = query
	if params == nil {
		params = []interface{}{}
	}
	value.Params = params

	body, err := json.Marshal(value)

	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://raikiri.db/execute", bytes.NewReader(body))

	if err != nil {
		return 0, err
	}

	req.Header.Set("Connection-Id", connection.connection_id)

	client := wasihttp.NewClient()
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	count, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	countInt, err := strconv.Atoi(string(count))

	if err != nil {
		return 0, err
	}

	return countInt, nil
}

func (connection *PostgresConnection) QuerySql(query string, params []interface{}) ([]byte, error) {

	var value struct {
		Query  string
		Params []interface{}
	}

	value.Query = query
	if params == nil {
		params = []interface{}{}
	}
	value.Params = params

	body, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://raikiri.db/query", bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Connection-Id", connection.connection_id)

	client := wasihttp.NewClient()
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return respBody, nil
}
