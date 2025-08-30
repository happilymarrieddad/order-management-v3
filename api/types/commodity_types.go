package types

import "fmt"

// CommodityType represents the type of a commodity.
type CommodityType int

const (
	// CommodityTypeUnknown is for when the type is not known.
	CommodityTypeUnknown CommodityType = iota
	// CommodityTypeProduce represents produce commodities.
	CommodityTypeProduce
)

// String returns the string representation of a CommodityType.
func (ct CommodityType) String() string {
	switch ct {
	case CommodityTypeProduce:
		return "produce"
	default:
		return "unknown"
	}
}

// ParseCommodityType converts a string to a CommodityType.
func ParseCommodityType(s string) (CommodityType, error) {
	switch s {
	case "produce":
		return CommodityTypeProduce, nil
	default:
		return CommodityTypeUnknown, fmt.Errorf("unknown commodity type: %s", s)
	}
}
