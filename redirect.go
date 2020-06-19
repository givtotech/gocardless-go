package gocardless

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	redirectEndpoint = "redirect_flows"
)

type (
	// Redirect Mandates represent the Direct Debit mandate with a customer.
	Redirect struct {
		// ID is a unique identifier, beginning with "RE".
		ID string `json:"id,omitempty"`
		// customer’s session token which must be provided when the redirect flow is completed.
		SessionToken string `json:"session_token"`
		// The URL to redirect to upon successful mandate setup.
		SuccessRedirectURL string `json:"success_redirect_url"`
		// The URL to redirect a customer to enter payment details.
		RedirectURL string `json:"redirect_url,omitempty"`
		// Description of the item the customer is paying for
		Description string `json:"description,omitempty"`
		// payment scheme
		Scheme string `json:"scheme,omitempty"`
		// CreatedAt is a fixed timestamp, recording when the redrect was created.
		CreatedAt *time.Time `json:"created_at,omitempty"`
		// prefill the details of a customer
		Customer Customer `json:"prefilled_customer"`
		// Links links to cusomer and bank accounts
		Links redirectLinks `json:"links"`
		// Metadata is a key-value store of custom data. Up to 3 keys are permitted, with key names up to 50
		// characters and values up to 500 characters.
		Metadata map[string]string `json:"metadata,omitempty"`
	}
	redirectLinks struct {
		CreditorID            string `json:"creditor,omitempty"`
		CustomerID            string `json:"customer,omitempty"`
		CustomerBankAccountID string `json:"customer_bank_account,omitempty"`
		MandateID             string `json:"mandate,omitempty"`
	}
	// redirectWrapper is a utility struct used to wrap and unwrap the JSON request being passed to the remote API
	redirectWrapper struct {
		Redirect *Redirect `json:"redirect_flows"`
	}
)

func (r *Redirect) String() string {
	bs, _ := json.Marshal(r)
	return string(bs)
}

// NewRedirect instantiate new Redirect object
func NewRedirect(sessionToken string, redirectURL string) *Redirect {
	return &Redirect{
		SessionToken:       sessionToken,
		SuccessRedirectURL: redirectURL,
	}
}

// AddMetadata adds new metadata item to mandate object
func (r *Redirect) AddMetadata(key, value string) {
	r.Metadata[key] = value
}

// CreateRedirect creates a new redirect object.
//
// Relative endpoint: POST /redirect_flows
func (c *Client) CreateRedirect(redirect *Redirect) error {
	redirectReq := &redirectWrapper{redirect}

	err := c.post(redirectEndpoint, redirectReq, redirectReq)
	if err != nil {
		return err
	}

	return err
}

// GetRedirect retrieves the details of an existing redirect.
//
// Relative endpoint: GET /redirect_flows/RE123
func (c *Client) GetRedirect(id string) (*Redirect, error) {
	wrapper := &redirectWrapper{}

	err := c.get(fmt.Sprintf(`%s/%s`, redirectEndpoint, id), wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.Redirect, err
}

// CompleteRedirect Completes a redirect object. creates a customer, customer bank account, and mandate objects
//
// Relative endpoint: POST /redirect_flows/RE123/actions/complete
func (c *Client) CompleteRedirect(redirect *Redirect) error {
	//
	rdData := map[string]interface{}{
		"data": map[string]interface{}{
			"session_token": redirect.SessionToken,
		},
	}

	redirectReq := &redirectWrapper{redirect}

	err := c.post(fmt.Sprintf(`%s/%s/actions/complete`, redirectEndpoint, redirect.ID), rdData, redirectReq)
	if err != nil {
		return err
	}
	return err
}
