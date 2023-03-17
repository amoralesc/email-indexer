package zinc

// ZincService is a service that interacts with the zinc server.
type ZincService struct {
	Url      string
	User     string
	Password string
}

// NewZincService creates a new zinc service (Dependency Injection).
func NewZincService(url, user, password string) *ZincService {
	return &ZincService{
		Url:      url,
		User:     user,
		Password: password,
	}
}

// ZincService Singleton
var Service *ZincService
