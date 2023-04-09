package zinc

// ZincService is a service that interacts with the zinc server.
type ZincService struct {
	Url      string
	User     string
	Password string
}

// NewZincService returns a new zinc service.
func NewZincService(url, user, password string) *ZincService {
	return &ZincService{
		Url:      url,
		User:     user,
		Password: password,
	}
}

// StartZincService starts the zinc service singleton.
func StartZincService(url, user, password string) {
	Service = NewZincService(url, user, password)
}

// ZincService Singleton
var Service *ZincService
