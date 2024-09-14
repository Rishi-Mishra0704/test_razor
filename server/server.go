package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/TheRSTech/test_razor/models"
	"github.com/TheRSTech/test_razor/views"
	"github.com/labstack/echo/v4"
	"github.com/razorpay/razorpay-go"
)

type Server struct {
	razor  *razorpay.Client
	router *echo.Echo
}

func NewServer(client razorpay.Client) *Server {
	server := &Server{
		razor: &client,
	}
	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	e := echo.New()

	// API routes
	e.POST("/api/connect-account", server.ConnectAccount)
	e.POST("/api/make-payment", server.MakePayment)
	e.GET("/api/transactions", server.GetTransactions)
	e.POST("/api/webhook", server.HandleWebhook) // Webhook endpoint

	// Serve the index page
	e.GET("/", func(c echo.Context) error {
		return render(c, views.Index())
	})

	server.router = e
}

func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

// ConnectAccount handles account creation and fund account linking
func (s *Server) ConnectAccount(c echo.Context) error {
	var req models.AccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
	}

	contactData := models.ContactData{
		Name:  req.Name,
		Email: req.Email,
		Type:  req.AccountType,
		Contact: models.ContactInfo{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		},
	}

	contact, err := s.razor.Customer.Create(contactData.ToMap(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create contact: " + err.Error()})
	}

	fundAccountData := models.FundAccountData{
		ContactID:   contact["id"].(string),
		AccountType: "bank_account",
		BankAccount: models.BankAccountInfo{
			Name:          req.AccountHolderName,
			IFSC:          req.IFSCCode,
			AccountNumber: req.AccountNumber,
		},
	}

	fundAccount, err := s.razor.FundAccount.Create(fundAccountData.ToMap(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create fund account: " + err.Error()})
	}

	return c.JSON(http.StatusOK, models.AccountResponse{
		Contact:     contact,
		FundAccount: fundAccount,
	})
}

// MakePayment handles bank-to-bank fund transfer
func (s *Server) MakePayment(c echo.Context) error {
	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
	}

	payoutData := models.PayoutData{
		AccountNumber: req.FromAccountID,
		FundAccountID: req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Mode:          "IMPS", // or NEFT, RTGS depending on your requirements
		Purpose:       "payout",
	}

	payout, err := s.razor.Order.Create(payoutData.ToMap(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create payout: " + err.Error()})
	}

	return c.JSON(http.StatusOK, payout)
}

// GetTransactions fetches transaction details
func (s *Server) GetTransactions(c echo.Context) error {
	accountID := c.QueryParam("account_id")
	if accountID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "account_id is required"})
	}

	from := c.QueryParam("from")
	to := c.QueryParam("to")

	options := models.TransactionOptions{
		AccountID: accountID,
		From:      from,
		To:        to,
	}

	transactions, err := s.razor.Customer.All(options.ToMap(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch transactions: " + err.Error()})
	}

	return c.JSON(http.StatusOK, transactions)
}

// HandleWebhook handles incoming Razorpay webhook events
func (s *Server) HandleWebhook(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to read request body"})
	}

	razorpaySignature := c.Request().Header.Get("X-Razorpay-Signature")
	if razorpaySignature == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Missing Razorpay signature"})
	}

	secret := "YOUR_RAZORPAY_WEBHOOK_SECRET" // Store this in environment variables
	if !s.verifyWebhookSignature(body, razorpaySignature, secret) {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid Razorpay signature"})
	}

	// Process webhook event here
	// For example, handle payout success, failure, etc.
	// Event details will be in 'body' (can unmarshal and handle as per event type)

	return c.JSON(http.StatusOK, "Webhook received and verified successfully")
}

// verifyWebhookSignature verifies the Razorpay webhook signature
func (s *Server) verifyWebhookSignature(payload []byte, receivedSignature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(receivedSignature), []byte(expectedSignature))
}
