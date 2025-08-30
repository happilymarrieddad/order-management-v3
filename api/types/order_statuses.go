package types

// OrderStatus defines the type for an order's status. The value is a snake_case string
// that is safe for programmatic use and database storage.
type OrderStatus string

// Defines all possible order statuses as constants for type safety.
const (
	OrderStatusPendingAcceptance OrderStatus = "pending_acceptance"
	OrderStatusPendingBooking    OrderStatus = "pending_booking"
	OrderStatusHold              OrderStatus = "hold"
	OrderStatusBooked            OrderStatus = "booked"
	OrderStatusShippedInTransit  OrderStatus = "shipped_in_transit"
	OrderStatusDelivered         OrderStatus = "delivered"
	OrderStatusReadyToInvoice    OrderStatus = "ready_to_invoice"
	OrderStatusInvoiced          OrderStatus = "invoiced"
	OrderStatusRejected          OrderStatus = "rejected"
	OrderStatusCancelled         OrderStatus = "cancelled"
	OrderStatusHoldForPOD        OrderStatus = "hold_for_pod"
	OrderStatusOrderTemplate     OrderStatus = "order_template"
	OrderStatusPaidInFull        OrderStatus = "paid_in_full"
)

// OrderStatusInfo holds the metadata associated with an OrderStatus.
type OrderStatusInfo struct {
	DisplayName string // The full, human-readable name (e.g., "Pending Acceptance")
	ShortName   string // A shorter name for condensed UIs (e.g., "Pending")
}

// OrderStatuses provides a lookup map to get metadata for any given OrderStatus.
// This acts as the single source of truth for all status properties.
var OrderStatuses = map[OrderStatus]OrderStatusInfo{
	OrderStatusPendingAcceptance: {"Pending Acceptance", "Pending"},
	OrderStatusPendingBooking:    {"Pending Booking", "Booking"},
	OrderStatusHold:              {"Hold", "Hold"},
	OrderStatusBooked:            {"Booked", "Booked"},
	OrderStatusShippedInTransit:  {"Shipped / In-Transit", "In-Transit"},
	OrderStatusDelivered:         {"Delivered", "Delivered"},
	OrderStatusReadyToInvoice:    {"Ready to Invoice", "Ready to Invoice"},
	OrderStatusInvoiced:          {"Invoiced", "Invoiced"},
	OrderStatusRejected:          {"Rejected", "Rejected"},
	OrderStatusCancelled:         {"Cancelled", "Cancelled"},
	OrderStatusHoldForPOD:        {"Hold for POD", "Invoice Approval"},
	OrderStatusOrderTemplate:     {"Order Template", "Order Template"},
	OrderStatusPaidInFull:        {"Paid In Full", "Paid in Full"},
}

// DisplayName returns the human-readable display name of the status.
// It provides a safe way to access the metadata.
func (s OrderStatus) DisplayName() string {
	if info, ok := OrderStatuses[s]; ok {
		return info.DisplayName
	}
	return string(s) // Fallback to the raw value if not found
}

// ShortName returns the short name of the status.
func (s OrderStatus) ShortName() string {
	if info, ok := OrderStatuses[s]; ok {
		return info.ShortName
	}
	return string(s) // Fallback to the raw value if not found
}

// IsValid checks if the status is a defined, valid status.
func (s OrderStatus) IsValid() bool {
	_, ok := OrderStatuses[s]
	return ok
}

// AllOrderStatuses returns a slice of all valid order statuses.
// Useful for populating dropdowns or for validation lists.
func AllOrderStatuses() []OrderStatus {
	statuses := make([]OrderStatus, 0, len(OrderStatuses))
	for status := range OrderStatuses {
		statuses = append(statuses, status)
	}
	return statuses
}
