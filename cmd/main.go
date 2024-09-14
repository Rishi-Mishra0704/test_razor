package main

import "github.com/TheRSTech/test_razor/server"

func main() {
	client := server.InitRazorpayClient()
	server := server.NewServer(*client)
	server.Start(":8080")
}
