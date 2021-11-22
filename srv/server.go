package srv

type Server interface {
	ListenAndServe() error
}
