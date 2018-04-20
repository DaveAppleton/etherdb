package etherdb

import (
	"database/sql"

	_ "github.com/lib/pq" // DB selection
)

// TokenTransfer = the data for a token transfer :-)
type TokenTransfer struct {
	TransferID  uint64
	TokenID     uint64
	BlockNumber uint64
	BlockHash   string
	Timestamp   uint64
	Index       uint
	TxHash      string
	Source      string
	Dest        string
	Amount      string
}

type NamedTokenTransfer struct {
	TransferID    uint64
	TokenID       uint64
	TokenName     string
	TokenSymbol   string
	TokenDecimals uint8
	TokenAddress  string
	BlockNumber   uint64
	BlockHash     string
	Timestamp     uint64
	Index         uint
	TxHash        string
	Source        string
	Dest          string
	Amount        string
}

// Add this token to the database
func (tt *TokenTransfer) Add() (err error) {
	statement := `insert into tokentransfers (tokenid,blocknumber,blockhash,index,txhash,source,dest,amount,timestamp) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning transferid`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(tt.TokenID, tt.BlockNumber, tt.BlockHash, tt.Index, tt.TxHash, tt.Source, tt.Dest, tt.Amount, tt.Timestamp).Scan(&tt.TransferID)
	return
}

// AddIfNotFound only adds if the txhash is not found
func (tt *TokenTransfer) AddIfNotFound() (err error) {
	statement := `insert into tokentransfers (txhash,tokenid,blocknumber,blockhash,index,source,dest,amount,timestamp) 
					values ($1,$2,$3,$4,$5,$6,$7,$8,$9) on conflict do nothing`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(tt.TxHash, tt.TokenID, tt.BlockNumber, tt.BlockHash, tt.Index, tt.Source, tt.Dest, tt.Amount, tt.Timestamp)
	return
}

func getTransfers(rows *sql.Rows) (transfers []TokenTransfer, err error) {
	var t TokenTransfer
	for rows.Next() {
		err = rows.Scan(&t.TransferID, &t.TokenID, &t.BlockNumber, &t.BlockHash, &t.Index, &t.TxHash, &t.Source, &t.Dest, &t.Amount, &t.Timestamp)
		if err != nil {
			return
		}
		transfers = append(transfers, t)
	}
	return
}

func getNamedTransfers(rows *sql.Rows) (transfers []NamedTokenTransfer, err error) {
	var t NamedTokenTransfer
	for rows.Next() {
		err = rows.Scan(&t.TransferID, &t.TokenID, &t.BlockNumber, &t.BlockHash, &t.Index, &t.TxHash, &t.Source, &t.Dest, &t.Amount, &t.Timestamp, &t.TokenName, &t.TokenSymbol, &t.TokenDecimals, &t.TokenAddress)
		if err != nil {
			return
		}
		transfers = append(transfers, t)
	}
	return
}

// Find transfers that match tokenID
func (tt *TokenTransfer) Find() (transfers []TokenTransfer, err error) {
	statement := `select transferid,tokenid,blocknumber,blockhash,index,txhash,source,dest,amount,timestamp from tokentransfers where tokenid=$1`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(tt.TokenID)
	if err != nil {
		return
	}
	return getTransfers(rows)

}

// FindByAddress returns transfers of specific token to or from an address
func (tt *TokenTransfer) FindByAddress(addr string) (transfers []TokenTransfer, err error) {
	statement := `select transferid,tokenid,blocknumber,blockhash,index,txhash,source,dest,amount,timestamp from tokentransfers where tokenid=$1 and (source=$2 or dest=$2)`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(tt.TokenID, addr)
	if err != nil {
		return
	}
	return getTransfers(rows)
}

// FindAllByAddress returns transfers of any token to or from an address newest first
func (tt *NamedTokenTransfer) FindAllByAddress(addr string) (transfers []NamedTokenTransfer, err error) {
	statement := `select transferid,A.tokenid,blocknumber,blockhash,index,txhash,source,dest,amount,timestamp,name,symbol,decimals,address
		from tokentransfers A, tokens B
		 where A.tokenid = B.tkn
		 and (source=$1 or dest=$1)
		 order by blocknumber desc, index desc`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(addr)
	if err != nil {
		return
	}
	return getNamedTransfers(rows)
}

// MaxBlock for a specified token
func (tt *TokenTransfer) MaxBlock() (max int64, err error) {
	statement := `select max(blocknumber) from tokentransfers where tokenid=$1`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	err = stmt.QueryRow(tt.TokenID).Scan(&max)
	if err != nil {
		max = 0
		err = nil
	}
	return
}
