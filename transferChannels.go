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

// NamedTransferChannel is a channel returning all the latest transfers
type NamedTransferChannel struct {
	ID    uint64
	Ch    chan NamedTokenTransfer
	Error error
}

// MakeTokenTransferChannel returns a channel to token transfers
func MakeNamedTokenTransferChannel(ID uint64) (ch NamedTransferChannel) {
	ch = NamedTransferChannel{ID: ID, Ch: make(chan NamedTokenTransfer, 50)}
	return
}

// ReadTokenTransfers send db to channel
func (ch *NamedTransferChannel) ReadNamedTokenTransfers(addr string) (err error) {
	defer close(ch.Ch)
	tt := NamedTokenTransfer{TokenID: ch.ID}
	ttRes, err := tt.FindAllByAddress(addr)
	if err != nil {
		ch.Error = err
		return
	}
	fmt.Println(len(ttRes), " results found for ",addr)
	for _, transfer := range ttRes {
		ch.Ch <- transfer
	}
	return
}
