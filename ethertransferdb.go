package etherdb

import (
	_ "github.com/lib/pq" // DB selection
)

// EtherTransfer = the data for a token transfer :-)
type EtherTransfer struct {
	TransferID  uint64
	BlockNumber uint64
	BlockHash   string
	Index       uint
	TxHash      string
	Source      string
	Dest        string
	Amount      string
}

// Add this token to the database
func (tt *EtherTransfer) Add() (err error) {
	statement := `insert into ethertransfers (blocknumber,blockhash,index,txhash,source,dest,amount) values ($1,$2,$3,$4,$5,$6,$7) returning transferid`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(tt.BlockNumber, tt.BlockHash, tt.Index, tt.TxHash, tt.Source, tt.Dest, tt.Amount).Scan(&tt.TransferID)
	return
}

// Find transfers that match tokenID
func (tt *EtherTransfer) Find() (transfers []TokenTransfer, err error) {
	statement := `select transferid,blocknumber,blockhash,index,txhash,source,dest,amount from ethertransfers`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	var t TokenTransfer
	for rows.Next() {
		err = rows.Scan(&t.TransferID, &t.TokenID, &t.BlockNumber, &t.BlockHash, &t.Index, &t.TxHash, &t.Source, &t.Dest, &t.Amount)
		if err != nil {
			return
		}
		transfers = append(transfers, t)
	}
	return
}
