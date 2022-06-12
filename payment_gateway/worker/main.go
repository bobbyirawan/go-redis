package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	go func() {
		for {
			result, err := redisClient.BLPop(0*time.Second, "payments").Result()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				params := map[string]interface{}{}
				err := json.NewDecoder(strings.NewReader(string(result[1]))).Decode(&params)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					paymentId, err := savePayment(params)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						fmt.Println("Payment # " + strconv.FormatInt(paymentId, 10) + " processed successfully.\r\n")
					}
				}

			}
		}
	}()

	select {}
}

func savePayment(params map[string]interface{}) (int64, error) {

	db, err := sql.Open("mysql", "root:rahasia@tcp(127.0.0.1:4444)/appointment_offline_simulation")

	if err != nil {
		return 0, err
	}

	defer db.Close()

	queryString := `insert into payments (
                        payment_date,
                        first_name,
                        last_name,
                        payment_mode,
                        payment_ref_no,
                        amount
                    ) values (
                        ?,
                        ?,
                        ?,
                        ?,
                        ?,
                        ?
                    )`

	stmt, err := db.Prepare(queryString)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	paymentDate := time.Now().Format("2006-01-02 15:04:05")
	firstName := params["first_name"]
	lastName := params["last_name"]
	paymentMode := params["payment_mode"]
	paymentRefNo := params["payment_ref_no"]
	amount := params["amount"]

	res, err := stmt.Exec(paymentDate, firstName, lastName, paymentMode, paymentRefNo, amount)

	if err != nil {
		return 0, err
	}

	paymentId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return paymentId, nil
}
