package es

// Scanner is the interface each index must imply.
type Scanner interface {
	IndexName() string
	Depth() int
	Mappings() map[string]map[string]map[string]string
	ScanUpdate(int, int) []interface{}
	ScanDelete(int, int) []interface{}
}
