package gocardless

import (
	"encoding/json"
	"time"
)

type (
	// Event objects represent events passed by gocardless's webhook notifications
	Event struct {
		// ID is a unique identifier, beginning with "PM".
		ID string `json:"id,omitempty"`
		// CreatedAt is a fixed timestamp, recording when the payment was created.
		CreatedAt *time.Time `json:"created_at,omitempty"`
		// ResourceType of the event is associated with
		ResourceType string `json:"resource_type,omitempty"`
		// Action performed on the resource type
		Action string `json:"action,omitempty"`
		// Links to cusomer and payment
		Links eventLinks `json:"links"`
		//
		Details struct {
			// source of event i.e. API
			Origin string `json:"origin,omitempty"`
			// description code
			Cause string `json:"cause,omitempty"`
			// long form description of the detail
			Description string `json:"description,omitempty"`
			// payment scheme
			Scheme string `json:"scheme,omitempty"`
			// scheme specfic event code
			ReasonCode string `json:"reason_code,omitempty"`
		} `json:"details"`
		// Metadata is a key-value store of custom data. Up to 3 keys are permitted, with key names up to 50
		// characters and values up to 500 characters.
		Metadata map[string]string `json:"metadata,omitempty"`
	}
	eventLinks struct {
		RefundID             string `json:"refund,omitempty"`
		MandateID            string `json:"mandate,omitempty"`
		PaymentID            string `json:"payment,omitempty"`
		InstalmentScheduleID string `json:"instalment_schedule,omitempty"`
		CreditorID           string `json:"creditor,omitempty"`
		SubscriptionID       string `json:"subscription,omitempty"`
		PayoutID             string `json:"payout,omitempty"`
	}

	// EventList a List of Events
	EventList struct {
		Events []*Event `json:"events"`
	}
)

func (e *Event) String() string {
	bs, _ := json.Marshal(e)
	return string(bs)
}
