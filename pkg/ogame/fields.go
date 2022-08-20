package ogame

// Fields planet fields stats
type Fields struct {
	Built int64
	Total int64
}

// HasFieldAvailable returns either or not we can still build on this planet
func (f Fields) HasFieldAvailable() bool {
	return f.Built < f.Total
}
