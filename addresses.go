package etherdb

// Account - hex address
type Account struct {
	AcctID  uint64
	User    string
	Address string
}

// Add a new address into the db
func (a *Account) Add() (err error) {
	statement := `insert into addresses (address,userid) values ($1,$2) returning id`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(a.Address, a.User).Scan(&a.AcctID)
	return
}

// FindAll tokens that match either on name or address
func (a *Account) FindAll() (accounts []Account, err error) {
	statement := `select id, address,userid from addresses where address=$1 or user=$2`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(a.Address, a.User)
	if err != nil {
		return
	}
	var acc Account
	for rows.Next() {
		err = rows.Scan(&acc.AcctID, &acc.Address, &acc.User)
		if err != nil {
			return
		}
		accounts = append(accounts, acc)
	}
	return
}

// GetAllAccountsAndMaxFrom max as a starting value, from the database
func GetAllAccountsAndMaxFrom(max uint64) (accounts []Account, newMax uint64, err error) {
	statement := `select id , address,userid from addresses order by user where id >= $1`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(max)
	if err != nil {
		return
	}
	var acc Account
	for rows.Next() {
		err = rows.Scan(&acc.AcctID, &acc.Address, &acc.User)
		if err != nil {
			return
		}
		accounts = append(accounts, acc)
		if acc.AcctID > max {
			max = acc.AcctID
		}
	}
	return
}

// GetAllAccountsAndMax from the database
func GetAllAccountsAndMax() (tokens []Account, max uint64, err error) {
	return GetAllAccountsAndMaxFrom(0)
}

// GetAllAccounts that are in the database
func GetAllAccounts() (tokens []Account, err error) {
	tokens, _, err = GetAllAccountsAndMax()
	return
}
