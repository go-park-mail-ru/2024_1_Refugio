package microservice_ports

type MicroServicePorts string

const (
	UserService     MicroServicePorts = "8001"
	EmailService    MicroServicePorts = "8002"
	SessionService  MicroServicePorts = "8003"
	AuthService     MicroServicePorts = "8004"
	FolderService   MicroServicePorts = "8005"
	QuestionService MicroServicePorts = "8006"
)

// GetPorts returns the port number associated with a MicroServicePorts enum value.
func GetPorts(gender MicroServicePorts) string {
	switch gender {
	case UserService:
		return "8001"
	case EmailService:
		return "8002"
	case SessionService:
		return "8003"
	case AuthService:
		return "8004"
	case FolderService:
		return "8005"
	case QuestionService:
		return "8006"
	default:
		return ""
	}
}
