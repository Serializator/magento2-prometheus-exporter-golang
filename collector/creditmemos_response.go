package collector

import "fmt"

type creditmemosResponse struct {
	Items []struct {
		State CreditMemoState `json:"state"`
	} `json:"items"`
}

type CreditMemoState int

const (
	CreditMemoOpen = iota + 1
	CreditMemoRefunded
	CreditMemoCancelled
)

func (state CreditMemoState) String() (string, error) {
	switch state {
	case CreditMemoOpen:
		return "open", nil
	case CreditMemoRefunded:
		return "refunded", nil
	case CreditMemoCancelled:
		return "cancelled", nil
	}

	return "", fmt.Errorf("unknown creditmemo state (%d)", state)
}
