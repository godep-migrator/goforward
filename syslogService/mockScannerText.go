package syslogService

//Mock ScannerText

type MockScannerText struct {
	TValue string
}

//simple mock function to replace the scanner text for unit tests
func (s *MockScannerText) Text() (r string) {
	return s.TValue
}
