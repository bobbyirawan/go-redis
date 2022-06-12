package main

import (
	"bytes"
	"fmt"
	"net/http"

	rds "test-queue/redis"
)

func main() {
	http.HandleFunc("/payments", paymentsHandler)
	http.ListenAndServe(":8080", nil)
}

func paymentsHandler(w http.ResponseWriter, req *http.Request) {

	buf := new(bytes.Buffer)

	// Include a Validation logic here to sanitize the req.Body when working in a production environment
	buf.ReadFrom(req.Body)
	paymentDetails := buf.String()
	err := rds.RedisClient.RPush("payments", paymentDetails).Err()
	if err != nil {
		fmt.Fprintf(w, err.Error()+"\r\n")
	} else {
		fmt.Fprintf(w, "Payment details accepted successfully\r\n")
	}
}
