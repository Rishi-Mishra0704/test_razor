package models

type AccountRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	AccountType       string `json:"account_type"`
	AccountNumber     string `json:"account_number"`
	IFSCCode          string `json:"ifsc_code"`
	BankName          string `json:"bank_name"`
	AccountHolderName string `json:"account_holder_name"`
}

type ContactData struct {
	Name    string      `json:"name"`
	Email   string      `json:"email"`
	Type    string      `json:"type"`
	Contact ContactInfo `json:"contact"`
}

type ContactInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type FundAccountData struct {
	ContactID   string          `json:"contact_id"`
	AccountType string          `json:"account_type"`
	BankAccount BankAccountInfo `json:"bank_account"`
}

type BankAccountInfo struct {
	Name          string `json:"name"`
	IFSC          string `json:"ifsc"`
	AccountNumber string `json:"account_number"`
}

type PaymentRequest struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
}

type PayoutData struct {
	AccountNumber string `json:"account_number"`
	FundAccountID string `json:"fund_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Mode          string `json:"mode"`
	Purpose       string `json:"purpose"`
}

type TransactionOptions struct {
	AccountID string `json:"account_id"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	Contact     map[string]interface{} `json:"contact"`
	FundAccount map[string]interface{} `json:"fund_account"`
}

// Helper methods to convert structs to map[string]interface{}
func (c ContactData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    c.Name,
		"email":   c.Email,
		"type":    c.Type,
		"contact": c.Contact,
	}
}

func (f FundAccountData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"contact_id":   f.ContactID,
		"account_type": f.AccountType,
		"bank_account": f.BankAccount,
	}
}

func (p PayoutData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"account_number":  p.AccountNumber,
		"fund_account_id": p.FundAccountID,
		"amount":          p.Amount,
		"currency":        p.Currency,
		"mode":            p.Mode,
		"purpose":         p.Purpose,
	}
}

func (t TransactionOptions) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"account_id": t.AccountID,
	}
	if t.From != "" {
		result["from"] = t.From
	}
	if t.To != "" {
		result["to"] = t.To
	}
	return result
}
