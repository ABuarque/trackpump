package service

// IDService defines how id serices should work
type IDService interface {
	Get() (string, error)
}
