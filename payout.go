package gocardless

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	payoutEndpoint = "payouts"
)

type (
	// Payout objects represent payouts from a customer to a creditor, taken against a Direct Debit payout.
	Payout struct {
		// ID is a unique identifier, beginning with "PO".
		ID string `json:"id,omitempty"`
		// Amount in pence (GBP), cents (AUD/EUR), öre (SEK), or øre (DKK).
		// e.g 1000 is 10 GBP in pence
		Amount int `json:"amount"`
		// ArrivalDate Date the payout is due to arrive in the creditor’s bank account.
		// If not specified, the payout will be collected as soon as possible
		ArrivalDate *Date `json:"arrival_date,omitempty"`
		// Fees that have already been deducted from the payout amount in minor unit
		// e.g 1000 is 10 GBP in pence
		DeductedFees int `json:"deducted_fees"`
		// Currency currency code, defaults to national currency of country_code
		Currency string `json:"currency"`
		// CreatedAt is a fixed timestamp, recording when the payout was created.
		CreatedAt *time.Time `json:"created_at,omitempty"`
		// Whether a payout contains merchant revenue or partner fees.
		PayoutType string `json:"payout_type,omitempty"`
		// Reference An optional payout reference that will appear on your customer’s bank statement
		Reference string `json:"reference,omitempty"`
		// Status status of payout.
		Status string `json:"status,omitempty"`
		// foreign exchange info
		FX fxInfo `json:"fx"`
		// ISO 4217 code for the currency in which tax is paid out to the tax authorities of your tax jurisdiction.
		TaxCurrency string `json:"tax_currency"`
		// Metadata is a key-value store of custom data. Up to 3 keys are permitted, with key names up to 50
		// characters and values up to 500 characters.
		Metadata map[string]string `json:"metadata,omitempty"`
		// Links to cusomer and payout
		Links payoutLinks `json:"links"`
	}
	payoutLinks struct {
		CreditorID     string `json:"creditor,omitempty"`
		CreditorBankID string `json:"creditor_bank_account,omitempty"`
	}
	fxInfo struct {
		EstimatedRate float64 `json:"estimated_exchange_rate,omitempty"` // Rate used in the foreign exchange of the amount into the fx_currency.
		Rate          float64 `json:"exchange_rate,omitempty"`           // Rate used in the foreign exchange of the amount into the fx_currency.
		Amount        int     `json:"fx_amount,omitempty"`               // Amount that was paid out in the fx_currency after foreign exchange.
		Currency      string  `json:"fx_currency,omitempty"`             // ISO 4217 code for the currency in which amounts will be paid out (after foreign exchange
	}
	// payoutWrapper is a utility struct used to wrap and unwrap the JSON request being passed to the remote API
	payoutWrapper struct {
		Payout *Payout `json:"payouts"`
	}

	// PayoutListResponse a List response of Payout instances
	PayoutListResponse struct {
		Payouts []*Payout `json:"payouts"`
		Meta    Meta      `json:"meta,omitempty"`
	}
)

func (p *Payout) String() string {
	bs, _ := json.Marshal(p)
	return string(bs)
}

// AddMetadata adds new metadata item to payout object
func (p *Payout) AddMetadata(key, value string) {
	p.Metadata[key] = value
}

// GetPayouts returns a cursor-paginated list of your payouts.
//
// Relative endpoint: GET /payouts
func (c *Client) GetPayouts(ctx context.Context) (*PayoutListResponse, error) {
	list := &PayoutListResponse{}

	err := c.get(ctx, payoutEndpoint, list)
	if err != nil {
		return nil, err
	}
	return list, err
}

// GetPayout retrieves the details of an existing payout.
//
// Relative endpoint: GET /payouts/PO123
func (c *Client) GetPayout(ctx context.Context, id string) (*Payout, error) {
	wrapper := &payoutWrapper{}

	err := c.get(ctx, fmt.Sprintf(`%s/%s`, payoutEndpoint, id), wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.Payout, err
}

// UpdatePayout Updates a payout object. Supports all of the fields supported when creating a payout.
//
// Relative endpoint: PUT /payouts/PM123
func (c *Client) UpdatePayout(ctx context.Context, payout *Payout) error {
	// allows only metadata
	payoutMeta := map[string]interface{}{
		"payouts": map[string]interface{}{
			"metadata": payout.Metadata,
		},
	}

	payoutReq := &payoutWrapper{payout}

	err := c.put(ctx, fmt.Sprintf(`%s/%s`, payoutEndpoint, payout.ID), payoutMeta, payoutReq)
	if err != nil {
		return err
	}
	return err
}
