package collector

import (
	"fmt"
)

type invoicesResponse struct {
	Items []struct {
		State InvoiceState `json:"state"`
	} `json:"items"`
}

type InvoiceState int

const (
	InvoiceOpen = iota + 1
	InvoicePaid
	InvoiceCancelled
)

func (state InvoiceState) String() (string, error) {
	switch state {
	case InvoiceOpen:
		return "open", nil
	case InvoicePaid:
		return "paid", nil
	case InvoiceCancelled:
		return "cancelled", nil
	}

	return "", fmt.Errorf("unknown invoice state (%d)", state)
}
