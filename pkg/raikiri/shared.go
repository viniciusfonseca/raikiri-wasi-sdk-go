package raikiri

type DbConnection interface {
	Execute(*DbConnection, []byte) int
	Query(*DbConnection, []byte) []byte
}
