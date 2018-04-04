package etherdb

/// QueueSending - currently a work in progress

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// QueueSend (tx, user, destination, valueStr, gasLimit, gasPrice, data)
//    Entry added to TX Queue, which is retried every 10 seconds
//    must wait for 1 of
//    * Tx Time out - how do we report it?
//    * Tx successful - (does it)
func QueueSend(tx common.Hash, user string, destination string, valueStr string, gasLimit string, gasPrice string, data string) {
	queueSQL := "insert into queueSend (txHash,user,dest,value,gasLimit,gasPrice,data) values ($1,$2,$3,$4,$5,$6,$7)"
	db.Exec(queueSQL, tx.Hex(), user, destination, valueStr, gasLimit, gasPrice, data)
}

// BumpRetry to indicate the number of times this has been retried
func BumpRetry(id *big.Int) {
	sql := "update queueSend set retryCount=retryCoun+1 where id=$1"
	_, err := db.Exec(sql, id)
	if err != nil {
		log.Println("bumpRetry ", err)
	}
}

func LogSecond(id *big.Int, user string, dest string, val string, tx common.Hash) {
	sql := "insert into seconds (id,user,dest,value,hash) values ($1,$2,$3,$4)"
	db.Exec(sql, id, user, dest, val, tx.Hex())
}

// SetFail message to show why not sent yet
func SetFail(id *big.Int, errIn error) {
	sql := "update queueSend set status = 'fail' and error='$1' where id=$2"
	_, err := db.Exec(sql, errIn.Error(), id)
	if err != nil {
		log.Println(err)
	}
}

// queue structure
// txHash
// user - whose account it is
// destination - who to send it to
// gasLimit
// gasPrice
// data
// status
// retryCount

// queueLoop
//

// func queueLoop() {
// 	var id *big.Int
// 	var hashStr string
// 	var userStr string
// 	var destStr string
// 	var valStr string
// 	var gasLStr string
// 	var gasPStr string
// 	var dataStr string
// 	var statStr string
// 	var retryCount *big.Int
// 	for {
// 		sql := `select id,txHash, user, destination, value, gasLmit, gasPrice, data, status, retryCount
// 				from queueSend where status = 'waiting'`
// 		for {
// 			rows, err := db.Query(sql, id, hashStr, userStr, destStr, valStr, gasLStr, gasPStr, dataStr, statStr, retryCount)
// 			if err != nil {
// 				log.Fatal("Q Loop ", err)
// 			}
// 			for rows.Next() {
// 				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 				defer cancel()
// 				_, ispend, err := client.TransactionByHash(ctx, common.HexToHash(hashStr))
// 				if err != nil {
// 					// update this record as error.
// 					setFail(id, err)
// 					continue
// 				}
// 				if ispend {
// 					bumpRetry(id)
// 					continue
// 					// still waiting
// 				}
// 				// if tx has passed, fire and forget the sendTx
// 				txr, err := client.TransactionReceipt(ctx, common.HexToHash(hashStr))
// 				if strings.Compare(txr.Status, "0") == 0 {
// 					// this should never happen - primary send fail
// 					setFail(id, errors.New("Primary send failed"))
// 				}
// 				userKey := ethKeys.NewKey("userKeys/" + userStr)
// 				userKey.LoadKey()                    // should work
// 				dest := common.HexToAddress(destStr) // should be a valid address
// 				value, _ := new(big.Int).SetString(valStr, 10)
// 				gasLimit, _ := strconv.Atoi(gasLStr)
// 				gasPrice, _ := new(big.Int).SetString(gasPStr, 10)
// 				data := common.Hex2Bytes(dataStr)
// 				tx2, err := sendEthereum(userKey, dest, value, uint64(gasLimit), gasPrice, data)
// 				logSecond(id, userStr, destStr, valStr, tx2.Hash())
// 			}
// 		}
// 	}
// }
