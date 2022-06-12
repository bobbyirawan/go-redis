package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v7"
)

func main() {
	http.HandleFunc("/payments", paymentsHandler)
	http.ListenAndServe(":8080", nil)
}

func paymentsHandler(w http.ResponseWriter, req *http.Request) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	buf := new(bytes.Buffer)

	// Include a Validation logic here to sanitize the req.Body when working in a production environment

	buf.ReadFrom(req.Body)

	paymentDetails := buf.String()

	err := redisClient.RPush("payments", paymentDetails).Err()

	if err != nil {
		fmt.Fprintf(w, err.Error()+"\r\n")
	} else {
		fmt.Fprintf(w, "Payment details accepted successfully\r\n")
	}

	// curl -i -X POST localhost:8080/payments -H "Content-Type: application/json" -d '{"first_name": "JOHN", "last_name": "DOE", "payment_mode": "CASH", "payment_ref_no": "-", "amount" : 5000.25}'
	// curl -i -X POST localhost:8080/payments -H "Content-Type: application/json" -d '{"first_name": "MARY", "last_name": "SMITH", "payment_mode": "CHEQUE", "payment_ref_no": "985", "amount" : 985.65}'
	// curl -i -X POST localhost:8080/payments -H "Content-Type: application/json" -d '{"first_name": "ANN", "last_name": "JACOBS", "payment_mode": "PAYPAL", "payment_ref_no": "ABC-XYZ", "amount" : 15.25}'
}
