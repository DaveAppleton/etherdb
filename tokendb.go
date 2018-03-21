package etherdb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // DB selection
)

// Token holds ERC20 token data
type Token struct {
	TokenID  uint64
	Symbol   string
	Name     string
	Decimals uint8
	Address  string
}

var db *sql.DB

// InitDB has to be called before any activity but needs viper
func InitDB(connectString string) {
	var err error
	db, err = sql.Open("postgres", connectString) //
	if err != nil {
		fmt.Println("open token database : ", err)
		log.Fatal("open token database : ", err)
	}
	fmt.Println("database is open")
	//defer db.Close()
}

// Add this token to the database
func (t *Token) Add() (err error) {
	statement := `insert into tokens (address,name,symbol, decimals) values ($1,$2,$3,$4) returning tkn`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.Address, t.Name, t.Symbol, t.Decimals).Scan(&t.TokenID)
	return
}

// FindAll tokens that match either on name or address
func (t *Token) FindAll() (tokens []Token, err error) {
	statement := `select tkn,address,name,symbol,decimals from tokens where symbol=$1 or address=$2`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(t.Symbol, t.Address)
	if err != nil {
		return
	}
	var tkn Token
	for rows.Next() {
		err = rows.Scan(&tkn.TokenID, &tkn.Address, &tkn.Name, &tkn.Symbol, &tkn.Decimals)
		if err != nil {
			return
		}
		tokens = append(tokens, tkn)
	}
	return
}

// GetAllTokensAndMaxFrom max as a starting value, from the database
func GetAllTokensAndMaxFrom(max uint64) (tokens []Token, newMax uint64, err error) {
	newMax = max
	statement := `select tkn, address,name,symbol,decimals from tokens where tkn > $1 order by name`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(max)
	if err != nil {
		return
	}
	var tkn Token
	for rows.Next() {
		err = rows.Scan(&tkn.TokenID, &tkn.Address, &tkn.Name, &tkn.Symbol, &tkn.Decimals)
		if err != nil {
			return
		}
		tokens = append(tokens, tkn)
		if tkn.TokenID > newMax {
			newMax = tkn.TokenID
		}
	}
	return
}

// GetAllTokensAndMax from the database
func GetAllTokensAndMax() (tokens []Token, max uint64, err error) {
	return GetAllTokensAndMaxFrom(uint64(0))
}

// GetAllTokens that are in the database
func GetAllTokens() (tokens []Token, err error) {
	tokens, _, err = GetAllTokensAndMax()
	return
}
