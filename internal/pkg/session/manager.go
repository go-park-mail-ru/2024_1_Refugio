package session

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"net/http"
	"time"

	"mail/internal/microservice/models/proto_converters"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"

	session_proto "mail/internal/microservice/session/proto"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
)

// GlobalSessionManager is a global instance of SessionsManager.
var (
	GlobalSessionManager = &SessionsManager{}
)

// SessionsManager manages user sessions.
type SessionsManager struct {
}

// InitializationGlobalSeaaionManager initializes the global session manager.
func InitializationGlobalSeaaionManager(sessionManager *SessionsManager) {
	GlobalSessionManager = sessionManager
}

// NewSessionsManager creates a new instance of SessionsManager.
func NewSessionsManager() *SessionsManager {
	return &SessionsManager{}
}

// SetSession set the session in the request.
func (sm *SessionsManager) SetSession(sessionId string, w http.ResponseWriter, r *http.Request, ctx context.Context) error {
	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return fmt.Errorf("connection fail")
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	sess, errStatus := sessionServiceClient.GetSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.GetSessionRequest{SessionId: sessionId},
	)
	if errStatus != nil {
		return fmt.Errorf("session already exist")
	}

	w.Header().Set("X-Csrf-Token", sess.Session.CsrfToken)

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.Session.SessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, sessionCookie)

	return nil
}

// GetSession retrieves the session from the request.
func (sm *SessionsManager) GetSession(r *http.Request, ctx context.Context) *api.Session {
	sessionCookie, _ := r.Cookie("session_id")

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return nil
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	sessionProto, errStatus := sessionServiceClient.GetSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.GetSessionRequest{SessionId: sessionCookie.Value},
	)
	if errStatus != nil {
		return nil
	}

	sessionCore := proto_converters.SessionConvertProtoInCore(*sessionProto.Session)

	return converters.SessionConvertCoreInApi(*sessionCore)
}

// Check checks the validity of the session and CSRF token in the request.
func (sm *SessionsManager) Check(r *http.Request, ctx context.Context) (*api.Session, error) {
	csrfToken := r.Header.Get("X-Csrf-Token")
	if r.URL.Path != "/api/v1/verify-auth" && csrfToken == "" {
		return nil, fmt.Errorf("CSRF token not found in request headers")
	}

	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, fmt.Errorf("no session found")
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return nil, fmt.Errorf("connection fail")
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	sessionProto, errStatus := sessionServiceClient.GetSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.GetSessionRequest{SessionId: sessionCookie.Value},
	)
	if errStatus != nil {
		return nil, fmt.Errorf("no session found")
	}

	sessionCore := proto_converters.SessionConvertProtoInCore(*sessionProto.Session)

	if r.URL.Path != "/api/v1/verify-auth" && sessionCore.CsrfToken != csrfToken {
		return nil, fmt.Errorf("CSRF token mismatch")
	}

	return converters.SessionConvertCoreInApi(*sessionCore), nil
}

// CheckLogin checks if the login associated with the session matches the provided login.
func (sm *SessionsManager) CheckLogin(login string, r *http.Request, ctx context.Context) error {
	sessionCookie, _ := r.Cookie("session_id")

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return fmt.Errorf("connection fail")
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	loginProto, errStatus := sessionServiceClient.GetLoginBySession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.GetLoginBySessionRequest{SessionId: sessionCookie.Value},
	)
	if errStatus != nil {
		return errStatus
	}

	if loginProto.Login != login {
		return fmt.Errorf("no right sender email")
	}

	return nil
}

// GetLoginBySession retrieves the login associated with the session from the request.
func (sm SessionsManager) GetLoginBySession(r *http.Request, ctx context.Context) (string, error) {
	sessionCookie, _ := r.Cookie("session_id")

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return "", fmt.Errorf("connection fail")
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	loginProto, errStatus := sessionServiceClient.GetLoginBySession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.GetLoginBySessionRequest{SessionId: sessionCookie.Value},
	)
	if errStatus != nil {
		return "", errStatus
	}

	return loginProto.Login, nil
}

// Create creates a new session for the user and sets the session ID cookie in the response.
func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32, ctx context.Context) (*api.Session, error) {
	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return nil, fmt.Errorf("connection fail")
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	sessionId, errStatus := sessionServiceClient.CreateSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.CreateSessionRequest{Session: &session_proto.Session{UserId: userID,
			Device:   "",
			LifeTime: 60 * 60 * 24},
		},
	)
	if errStatus != nil {
		return nil, fmt.Errorf("session already exist")
	}

	sess, errStatus := sessionServiceClient.GetSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.GetSessionRequest{SessionId: sessionId.SessionId},
	)
	if errStatus != nil {
		return nil, fmt.Errorf("session already exist")
	}

	w.Header().Set("X-Csrf-Token", sess.Session.CsrfToken)

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.Session.SessionId,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, sessionCookie)

	sessionCore := proto_converters.SessionConvertProtoInCore(*sess.Session)

	return converters.SessionConvertCoreInApi(*sessionCore), nil
}

// DestroyCurrent destroys the current session by deleting the session ID cookie from the response.
func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request, ctx context.Context) error {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		return fmt.Errorf("connection fail")
	}
	defer conn.Close()

	sessionServiceClient := session_proto.NewSessionServiceClient(conn)
	status, errStatus := sessionServiceClient.DeleteSession(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
		&session_proto.DeleteSessionRequest{SessionId: sessionCookie.Value},
	)
	if errStatus != nil {
		return fmt.Errorf("session delete fail")
	}

	if !status.Status {
		return fmt.Errorf("no session found")
	}

	sessionCookieToDelete := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &sessionCookieToDelete)

	return nil
}
