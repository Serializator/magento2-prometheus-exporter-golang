package collector

type invoicesResponse struct {
	Items []struct {
		State State `json:"state"`
	} `json:"items"`
}

type State int

const (
	Open = iota + 1
	Paid
	Cancelled
)

func (state State) String() string {
	return toString[state]
}

var toString = map[State]string{
	Open:      "open",
	Paid:      "paid",
	Cancelled: "cancelled",
}
