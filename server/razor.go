package server

import (
	"github.com/razorpay/razorpay-go"
)

func InitRazorpayClient() *razorpay.Client {
	razor := razorpay.NewClient("key", "secret")
	return razor
}
