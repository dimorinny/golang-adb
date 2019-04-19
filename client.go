package golangadb

type Client interface {
	Devices() ([]Device, error)
}
