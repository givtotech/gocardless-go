package gocardless

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	subscriptionEndpoint = "subscriptions"
)

type (
	// Subscription objects represent payments according to a schedule
	Subscription struct {
		// ID is a unique identifier, beginning with "SB".
		ID string `json:"id,omitempty"`
		// CreatedAt is a fixed timestamp, recording when the subscription was created.
		CreatedAt *time.Time `json:"created_at,omitempty"`
		// Amount in pence (GBP), cents (AUD/EUR), öre (SEK), or øre (DKK).
		// e.g 1000 is 10 GBP in pence
		Amount int `json:"amount"`
		// Currency currency code, defaults to national currency of country_code
		Currency string `json:"currency"`
		// Status status of subscription.
		Status string `json:"status,omitempty"`
		// Name of subscription.
		Name string `json:"name,omitempty"`
		// StartDate A future date on which the subscription should start.
		StartDate *Date `json:"start_date,omitempty"`
		// Number of interval_units between customer charge dates.
		Interval int `json:"interval,omitempty"`
		// The total number of payments that should be taken by this subscription.
		Count int `json:"count,omitempty"`
		// The unit of time between customer charge dates. One of weekly, monthly or yearly.
		IntervalUnit string `json:"interval_unit"`
		// As per RFC 2445. The day of the month to charge customers on. 1-28 or -1 to indicate the last day of the month.
		DayOfMonth int `json:"day_of_month,omitempty"`
		// Name of the month on which to charge a customer. Must be lowercase. Only applies when the interval_unit is yearly
		Month int `json:"month,omitempty"`
		//An optional payment reference.
		PaymentReference int `json:"payment_reference,omitempty"`
		// The amount to be deducted from each payment as an app fee
		AppFee int `json:"app_fee,omitempty"`
		//
		UpcomingPayments []subscriptionPayment `json:"upcoming_payments,omitempty"`
		// Metadata is a key-value store of custom data. Up to 3 keys are permitted, with key names up to 50
		// characters and values up to 500 characters.
		Metadata map[string]string `json:"metadata,omitempty"`
		// Links to cusomer and payment
		Links subscriptionLinks `json:"links"`
		// On failure, automatically retry payments using intelligent retries
		Retry bool `json:"retry_if_possible,omitempty"`
	}
	subscriptionLinks struct {
		MandateID string `json:"mandate"`
	}
	subscriptionPayment struct {
		// ChargeDate A future date on which the payment should be collected.
		// If not specified, the payment will be collected as soon as possible
		ChargeDate *Date `json:"charge_date,omitempty"`
		// Amount in pence (GBP), cents (AUD/EUR), öre (SEK), or øre (DKK).
		// e.g 1000 is 10 GBP in pence
		Amount int `json:"amount"`
	}

	// paymentWrapper is a utility struct used to wrap and unwrap the JSON request being passed to the remote API
	subscriptionWrapper struct {
		Subscription *Subscription `json:"subscriptions"`
	}

	// SubscriptionListResponse a List response of Subscription instances
	SubscriptionListResponse struct {
		Subscriptions []*Subscription `json:"subscriptions"`
		Meta          Meta            `json:"meta,omitempty"`
	}
)

func (s *Subscription) String() string {
	bs, _ := json.Marshal(s)
	return string(bs)
}

// NewSubscription instantiate new subscription object
func NewSubscription(amount int, currency string, intervalUnit string, mandateID string) *Subscription {
	return &Subscription{
		Amount:       amount,
		Currency:     currency,
		IntervalUnit: intervalUnit,
		Links:        subscriptionLinks{MandateID: mandateID},
	}
}

// AddMetadata adds new metadata item to payment object
func (s *Subscription) AddMetadata(key, value string) {
	s.Metadata[key] = value
}

// CreateSubscription creates a new subscription object.
//
// Relative endpoint: POST /subscriptions
func (c *Client) CreateSubscription(subscription *Subscription) error {
	subscriptionReq := &subscriptionWrapper{subscription}

	err := c.post(subscriptionEndpoint, subscriptionReq, subscriptionReq)
	if err != nil {
		return err
	}

	return err
}

// GetSubscriptions returns a cursor-paginated list of your subscriptions.
//
// Relative endpoint: GET /subscriptions
func (c *Client) GetSubscriptions() (*SubscriptionListResponse, error) {
	list := &SubscriptionListResponse{}

	err := c.get(subscriptionEndpoint, list)
	if err != nil {
		return nil, err
	}
	return list, err
}

// GetSubscription retrieves the details of an existing subscription.
//
// Relative endpoint: GET /subscriptions/SB123
func (c *Client) GetSubscription(id string) (*Subscription, error) {
	wrapper := &subscriptionWrapper{}

	err := c.get(fmt.Sprintf(`%s/%s`, subscriptionEndpoint, id), wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.Subscription, err
}

// UpdateSubscription Updates a sSubscription object. Supports all of the fields supported when creating a subscription.
//
// Relative endpoint: PUT /subscriptions/SB123
func (c *Client) UpdateSubscription(subscription *Subscription) error {
	// allows only metadata
	subscriptionMeta := map[string]interface{}{
		"subscriptions": map[string]interface{}{
			"metadata": subscription.Metadata,
		},
	}

	subscriptionReq := &subscriptionWrapper{subscription}

	err := c.put(fmt.Sprintf(`%s/%s`, subscriptionEndpoint, subscription.ID), subscriptionMeta, subscriptionReq)
	if err != nil {
		return err
	}
	return err
}

// CancelSubscription immediately cancels a subscription.
//
// Relative endpoint: POST /subscriptions/SU123/actions/cancel
func (c *Client) CancelSubscription(subscription *Subscription) error {
	// allows only metadata
	subscriptionMeta := map[string]interface{}{
		"subscriptions": map[string]interface{}{
			"metadata": subscription.Metadata,
		},
	}

	wrapper := &subscriptionWrapper{subscription}

	err := c.post(fmt.Sprintf(`%s/%s/actions/cancel`, subscriptionEndpoint, subscription.ID), subscriptionMeta, wrapper)
	if err != nil {
		return err
	}
	return err
}

// PauseSubscription pauses a active subscription
//
// Relative endpoint: POST /subscriptions/SU123/actions/pause
func (c *Client) PauseSubscription(subscription *Subscription) error {
	// allows only metadata
	subscriptionMeta := map[string]interface{}{
		"subscriptions": map[string]interface{}{
			"metadata": subscription.Metadata,
		},
	}

	wrapper := &subscriptionWrapper{subscription}
	err := c.post(fmt.Sprintf(`%s/%s/actions/pause`, subscriptionEndpoint, subscription.ID), subscriptionMeta, wrapper)
	if err != nil {
		return err
	}
	return err
}

// ResumeSubscription resumes a paused subscription
//
// Relative endpoint: POST /subscriptions/SU123/actions/resume
func (c *Client) ResumeSubscription(subscription *Subscription) error {
	// allows only metadata
	subscriptionMeta := map[string]interface{}{
		"subscriptions": map[string]interface{}{
			"metadata": subscription.Metadata,
		},
	}

	wrapper := &subscriptionWrapper{subscription}
	err := c.post(fmt.Sprintf(`%s/%s/actions/resume`, subscriptionEndpoint, subscription.ID), subscriptionMeta, wrapper)
	if err != nil {
		return err
	}
	return err
}
