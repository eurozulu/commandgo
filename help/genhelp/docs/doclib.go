package docs

type DocEntry struct {
	Key   string
	Value string
}

type DocLibrary interface {
	Entry(k string) *DocEntry
}
