package etherdb

import "fmt"

// TransferChannel is a channel returning all the latest transfers
type TransferChannel struct {
	ID    uint64
	Ch    chan TokenTransfer
	Error error
}

// MakeTokenTransferChannel returns a channel to token transfers
func MakeTokenTransferChannel(ID uint64) (ch TransferChannel) {
	ch = TransferChannel{ID: ID, Ch: make(chan TokenTransfer, 50)}
	return
}

// ReadTokenTransfers send db to channel
func (ch *TransferChannel) ReadTokenTransfers() (err error) {
	defer close(ch.Ch)
	tt := TokenTransfer{TokenID: ch.ID}
	ttRes, err := tt.Find()
	if err != nil {
		ch.Error = err
		return
	}
	for _, transfer := range ttRes {
		fmt.Println(transfer.BlockNumber)
		ch.Ch <- transfer
	}
	return
}
