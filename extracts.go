package ogame

// Compile time checks to ensure type satisfies Extractor interface
var _ Extractor = ExtractorV6{}
var _ Extractor = (*ExtractorV6)(nil)
var _ Extractor = ExtractorV7{}
var _ Extractor = (*ExtractorV7)(nil)

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML []byte) int64 {
	return extractUniverseSpeedV6(pageHTML)
}
