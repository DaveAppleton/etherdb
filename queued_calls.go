package etherdb

// Queue holds ERC20 token data
type Queue struct {
	ID       uint64
	URL      string
	Data     string
	Received bool
	Count    uint16
}

// Add this token to the database
func (q *Queue) Add() (err error) {
	statement := `insert into queues (url,data,count,received) values ($1,$2,0,0) returning txn`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(q.URL, q.Data).Scan(&q.ID)
	return
}

// Update record using AccID as key
func (q *Queue) Update() (changes int64, err error) {
	statement := `update queues set (url=$1),(data=$2),(received=$3),(count=$4) where id=$5`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	result, err := stmt.Exec(q.URL, q.Data, q.Received, q.Count, q.ID)
	if err != nil {
		return
	}
	changes, err = result.RowsAffected()
	return
}

// GetAllUnhandledQueueFrom max as a starting value, from the database
func GetAllUnhandledQueueFrom(max uint64) (queues []Queue, newMax uint64, err error) {
	statement := `select id, url,data,count,received from queues order by id where id > $1 and not received`
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	rows, err := stmt.Query(max)
	if err != nil {
		return
	}
	var q Queue
	for rows.Next() {
		err = rows.Scan(q.ID, q.URL, q.Data, q.Count, q.Received)
		if err != nil {
			return
		}
		queues = append(queues, q)
		if q.ID > max {
			max = q.ID
		}
	}
	return
}

// GetAllUnhandledQueue from the database
func GetAllUnhandledQueue() (queues []Queue, newMax uint64, err error) {
	return GetAllUnhandledQueueFrom(uint64(0))
}
