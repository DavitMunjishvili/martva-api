package models

// AvailableDate represents the availability of a booking date.
type AvailableDate struct {
	BookingDate       string `json:"bookingDate"`
	BookingDateStatus int    `json:"bookingDateStatus"`
}

// CenterResult represents the result for a single center, including available dates.
type CenterResult struct {
	CenterID   int             `json:"centerId"`
	CenterName string          `json:"centerName"`
	Dates      []AvailableDate `json:"dates"`
	Error      string          `json:"error,omitempty"`
}
