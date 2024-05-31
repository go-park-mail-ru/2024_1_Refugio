// 89.208.223.140
package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
	_ "time/tzdata"

	"github.com/gorilla/mux"
	"github.com/kataras/requestid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	_ "github.com/jackc/pgx/stdlib"

	"mail/cmd/configs"
	"mail/internal/models/microservice_ports"
	"mail/internal/monitoring"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/middleware"
	"mail/internal/pkg/session"
	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/websocket"

	migrate "github.com/rubenv/sql-migrate"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	auth_proto "mail/internal/microservice/auth/proto"
	email_proto "mail/internal/microservice/email/proto"
	folder_proto "mail/internal/microservice/folder/proto"
	question_proto "mail/internal/microservice/questionnaire/proto"
	session_proto "mail/internal/microservice/session/proto"
	user_proto "mail/internal/microservice/user/proto"
	authHand "mail/internal/pkg/auth/delivery/http"
	emailHand "mail/internal/pkg/email/delivery/http"
	folderHand "mail/internal/pkg/folder/delivery/http"
	gmailAuthHand "mail/internal/pkg/gmail/gmail_auth/delivery/http"
	gmailEmailHand "mail/internal/pkg/gmail/gmail_handler/delivery/http"
	oauthHand "mail/internal/pkg/oauth/delivery/http"
	questionHand "mail/internal/pkg/questionnairy/delivery/http"
	userHand "mail/internal/pkg/user/delivery/http"

	_ "mail/docs"
)

// @title API MailHub
// @version 1.0
// @description API server for MailHub

// @host mailhub.su
// @BasePath /
func main() {
	fmt.Println("starting mail")
	settingTime()

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	loggerMiddlewareAccess := initializeMiddlewareLogger()

	sessionManagerServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.SessionService))
	if err != nil {
		log.Fatalf("connection with microservice auth fail")
	}
	defer sessionManagerServiceConn.Close()
	sessionsManager := initializeSessionsManager(session_proto.NewSessionServiceClient(sessionManagerServiceConn))

	authServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.AuthService))
	if err != nil {
		log.Fatalf("connection with microservice auth fail")
	}
	defer authServiceConn.Close()

	userServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer userServiceConn.Close()
	authHandler := initializeAuthHandler(sessionsManager, auth_proto.NewAuthServiceClient(authServiceConn), user_proto.NewUserServiceClient(userServiceConn))
	userHandler := initializeUserHandler(sessionsManager, user_proto.NewUserServiceClient(userServiceConn))

	emailServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer emailServiceConn.Close()
	emailHandler := initializeEmailHandler(sessionsManager, email_proto.NewEmailServiceClient(emailServiceConn))

	folderServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		log.Fatalf("connection with microservice user fail")
	}
	defer folderServiceConn.Close()
	folderHandler := initializeFolderHandler(sessionsManager, folder_proto.NewFolderServiceClient(folderServiceConn))

	questionServiceConn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		log.Fatalf("connection with microservice question fail")
	}
	defer questionServiceConn.Close()
	questionHandler := initializeQuestionHandler(sessionsManager, question_proto.NewQuestionServiceClient(questionServiceConn))

	oauthHandler := initializeOAuthHandler(sessionsManager, user_proto.NewUserServiceClient(userServiceConn))

	oauthGMailHandler := initializeGMailAuthHandler(sessionsManager, auth_proto.NewAuthServiceClient(authServiceConn), user_proto.NewUserServiceClient(userServiceConn))
	emailGMailHandler := initializeEmailGMailHandler(sessionsManager)
	router := setupRouter(authHandler, oauthHandler, oauthGMailHandler, userHandler, emailHandler, folderHandler, questionHandler, emailGMailHandler, loggerMiddlewareAccess)

	startServer(router)
}

// settingTime setting local time on server
func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error in location detection")
	}

	time.Local = loc
}

// initializeDatabase database initialization
func initializeDatabase() *sql.DB {
	db, err := sql.Open("pgx", configs.DSN)
	if err != nil {
		log.Fatalln("Can't parse config", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("Database is not available", err)
	}

	db.SetMaxOpenConns(10)

	return db
}

// migrateDatabase applying database migration
func migrateDatabase(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
}

// initializeSessionsManager initializing session manager
func initializeSessionsManager(sessionServiceClient session_proto.SessionServiceClient) *session.SessionsManager {
	sessionsManager := session.NewSessionsManager(sessionServiceClient)
	session.InitializationGlobalSessionManager(sessionsManager)

	startSessionCleaner(24*time.Hour, sessionServiceClient)

	return sessionsManager
}

// initializeAuthHandler initializing authorization handler
func initializeAuthHandler(sessionsManager *session.SessionsManager, authServiceClient auth_proto.AuthServiceClient, userServiceClient user_proto.UserServiceClient) *authHand.AuthHandler {
	return &authHand.AuthHandler{
		Sessions:          sessionsManager,
		AuthServiceClient: authServiceClient,
		UserServiceClient: userServiceClient,
	}
}

// initializeOAuthHandler initializing authorization handler
func initializeOAuthHandler(sessionsManager *session.SessionsManager, userServiceClient user_proto.UserServiceClient) *oauthHand.OAuthHandler {
	return &oauthHand.OAuthHandler{
		Sessions:          sessionsManager,
		UserServiceClient: userServiceClient,
	}
}

// initializeEmailHandler initializing email handler
func initializeEmailHandler(sessionsManager *session.SessionsManager, emailServiceClient email_proto.EmailServiceClient) *emailHand.EmailHandler {
	minioClient, err := minio.New(configs.ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(configs.ACCESSKEYID, configs.SECRETACCESSKEY, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	bucketName := "files"
	location := "eu-central-1"

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println("failed to check bucket existence")
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			fmt.Println("failed to create bucket")
		}
		fmt.Printf("Bucket has been successfully created: %s\n", bucketName)
	} else {
		fmt.Printf("Bucket %s already exists\n", bucketName)
	}

	err = minioClient.SetBucketPolicy(ctx, bucketName, generatePolicy(bucketName))
	if err != nil {
		fmt.Println("failed to set bucket policy")
	} else {
		fmt.Println("bucket policy set successfully")
	}

	return &emailHand.EmailHandler{
		Sessions:           sessionsManager,
		EmailServiceClient: emailServiceClient,
		MinioClient:        minioClient,
	}
}

// initializeUserHandler initializing user handler
func initializeUserHandler(sessionsManager *session.SessionsManager, userServiceClient user_proto.UserServiceClient) *userHand.UserHandler {
	minioClient, err := minio.New(configs.ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(configs.ACCESSKEYID, configs.SECRETACCESSKEY, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	bucketName := "photos"
	location := "eu-central-1"

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println("failed to check bucket existence")
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			fmt.Println("failed to create bucket")
		}
		fmt.Printf("Bucket has been successfully created: %s\n", bucketName)
	} else {
		fmt.Printf("Bucket %s already exists\n", bucketName)
	}

	err = minioClient.SetBucketPolicy(ctx, bucketName, generatePolicy(bucketName))
	if err != nil {
		fmt.Println("failed to set bucket policy")
	} else {
		fmt.Println("bucket policy set successfully")
	}

	return &userHand.UserHandler{
		Sessions:          sessionsManager,
		UserServiceClient: userServiceClient,
		MinioClient:       minioClient,
	}
}

// initializeGMailAuthHandler initializes the GMail authentication handler
func initializeGMailAuthHandler(sessionsManager *session.SessionsManager, authServiceClient auth_proto.AuthServiceClient, userServiceClient user_proto.UserServiceClient) *gmailAuthHand.GMailAuthHandler {
	return &gmailAuthHand.GMailAuthHandler{
		Sessions:          sessionsManager,
		AuthServiceClient: authServiceClient,
		UserServiceClient: userServiceClient,
	}
}

// initializeEmailGMailHandler initializes the GMail email handler using
func initializeEmailGMailHandler(sessionsManager *session.SessionsManager) *gmailEmailHand.GMailEmailHandler {
	return &gmailEmailHand.GMailEmailHandler{
		Sessions: sessionsManager,
	}
}

// generatePolicy policy generation for minio
func generatePolicy(bucketName string) string {
	return fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Effect": "Allow","Principal": {"AWS": ["*"]},"Action": ["s3:GetBucketLocation"],"Resource": ["arn:aws:s3:::%s"]},{"Effect": "Allow","Principal": {"AWS": ["*"]},"Action": ["s3:GetObject"],"Resource": ["arn:aws:s3:::%s/*"]}]}`, bucketName, bucketName)
}

// initializeFolderHandler initializing folder handler
func initializeFolderHandler(sessionsManager *session.SessionsManager, folderServiceClient folder_proto.FolderServiceClient) *folderHand.FolderHandler {
	return &folderHand.FolderHandler{
		Sessions:            sessionsManager,
		FolderServiceClient: folderServiceClient,
	}
}

// initializeQuestionHandler initializing question handler
func initializeQuestionHandler(sessionsManager *session.SessionsManager, questionServiceClient question_proto.QuestionServiceClient) *questionHand.QuestionHandler {
	return &questionHand.QuestionHandler{
		Sessions:              sessionsManager,
		QuestionServiceClient: questionServiceClient,
	}
}

// initializeMiddlewareLogger initializing logger
func initializeMiddlewareLogger() *middleware.Logger {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}

	logrusAccess := logger.InitializationAccessLog(f)
	loggerAccess := new(middleware.Logger)
	loggerAccess.Logger = logrusAccess

	return loggerAccess
}

// setupRouter configuring routers
func setupRouter(authHandler *authHand.AuthHandler, oauthHandler *oauthHand.OAuthHandler, oauthGMailHandler *gmailAuthHand.GMailAuthHandler, userHandler *userHand.UserHandler, emailHandler *emailHand.EmailHandler, folderHandler *folderHand.FolderHandler, questionHandler *questionHand.QuestionHandler, emailGMailHandler *gmailEmailHand.GMailEmailHandler, logger *middleware.Logger) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/testAuth/auth-vk/getAuthUrlSignUpVK", oauthHandler.GetSignUpURLVK).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/testAuth/auth-vk/getAuthUrlLoginVK", oauthHandler.GetLoginURLVK).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/testAuth/auth-vk/auth/{code}", oauthHandler.AuthVK).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/testAuth/auth-vk/loginVK/{code}", oauthHandler.LoginVK).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/testAuth/auth-vk/signupVK", oauthHandler.SignupVK).Methods("POST", "OPTIONS")

	auth := setupAuthRouter(authHandler, oauthHandler, oauthGMailHandler, emailHandler, logger)
	router.PathPrefix("/api/v1/auth").Handler(auth)

	logRouter := setupLogRouter(emailHandler, userHandler, folderHandler, questionHandler, emailGMailHandler, logger)
	router.PathPrefix("/api/v1").Handler(logRouter)

	staticDir := "/media/"
	staticFileServer := http.StripPrefix(staticDir, http.FileServer(http.Dir("./avatars")))
	router.PathPrefix(staticDir).Handler(staticFileServer)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Metrics = monitoring.RegisterMonitoring(router)

	return logger.AccessLogMiddleware(router)
}

// setupAuthRouter configuring authorization router
func setupAuthRouter(authHandler *authHand.AuthHandler, oauthHandler *oauthHand.OAuthHandler, oauthGMailHandler *gmailAuthHand.GMailAuthHandler, emailHandler *emailHand.EmailHandler, logger *middleware.Logger) http.Handler {
	auth := mux.NewRouter().PathPrefix("/api/v1/auth").Subrouter()
	auth.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware)

	r := websocket.NewRoom()
	auth.Handle("/web/websocket_connection/{login}", r)
	go r.Run()

	auth.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	auth.HandleFunc("/signup", authHandler.Signup).Methods("POST", "OPTIONS")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")
	auth.HandleFunc("/sendOther", emailHandler.SendFromAnotherDomain).Methods("POST", "OPTIONS")
	auth.HandleFunc("/addFileOther", emailHandler.AddFileFromAnotherDomain).Methods("POST", "OPTIONS")
	auth.HandleFunc("/addFileToEmailOther/{id}/file/{file-id}", emailHandler.AddFileToEmailFromAnotherDomain).Methods("POST", "OPTIONS")

	auth.HandleFunc("/getAuthURL", oauthGMailHandler.GetAuthURL).Methods("GET", "OPTIONS")
	auth.HandleFunc("/gAuth", oauthGMailHandler.GoogleAuth).Methods("GET", "OPTIONS")
	auth.HandleFunc("/signupGMailUser", oauthGMailHandler.SugnupGMail).Methods("POST", "OPTIONS")

	return auth
}

// templateHandler represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("/home/sergey/mailhub/2024_1_Refugio/cmd/mail/templates", t.filename)))
	})
	if err := t.templ.Execute(w, r); err != nil {
		fmt.Println("Error executing template:", err)
	}
}

// setupLogRouter configuring router with logger
func setupLogRouter(emailHandler *emailHand.EmailHandler, userHandler *userHand.UserHandler, folderHandler *folderHand.FolderHandler, questionHandler *questionHand.QuestionHandler, emailGMailHandler *gmailEmailHand.GMailEmailHandler, logger *middleware.Logger) http.Handler {
	logRouter := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	logRouter.Use(logger.AccessLogMiddleware, middleware.PanicMiddleware, middleware.AuthMiddleware)

	logRouter.HandleFunc("/verify-auth", userHandler.VerifyAuth).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/user/get", userHandler.GetUserBySession).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/user/update", userHandler.UpdateUserData).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/user/delete/{id}", userHandler.DeleteUserData).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/user/avatar/upload", userHandler.UploadUserAvatar).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/user/avatar/delete", userHandler.DeleteUserAvatar).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/user/count", userHandler.GetCountUsers).Methods("GET", "OPTIONS")

	logRouter.HandleFunc("/emails/incoming", emailHandler.Incoming).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/sent", emailHandler.Sent).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/draft", emailHandler.Draft).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/emails/spam", emailHandler.Spam).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/{id}", emailHandler.GetByID).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/update/{id}", emailHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/email/delete/{id}", emailHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/email/send", emailHandler.Send).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/email/adddraft", emailHandler.AddDraft).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/email/sendToOtherDomain/{id}", emailHandler.SendEmailToOtherDomains).Methods("POST", "OPTIONS")

	logRouter.HandleFunc("/email/{id}/addattachment", emailHandler.AddAttachment).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/email/addfile", emailHandler.AddFile).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/email/{id}/file/{file-id}", emailHandler.AddFileToEmail).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/email/get/file/{id}", emailHandler.GetFileByID).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/{id}/get/files/", emailHandler.GetFilesByEmailID).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/email/delete/file/{id}", emailHandler.DeleteFileByID).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/email/update/file/{id}", emailHandler.UpdateFileByID).Methods("PUT", "OPTIONS")

	logRouter.HandleFunc("/questions", questionHandler.GetAllQuestions).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/questions", questionHandler.AddQuestion).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/answers", questionHandler.AddAnswer).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/statistics", questionHandler.GetStatistics).Methods("GET", "OPTIONS")

	logRouter.HandleFunc("/folder/add", folderHandler.Add).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/folder/all", folderHandler.GetAll).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/folder/delete/{id}", folderHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/folder/update/{id}", folderHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/folder/add_email", folderHandler.AddEmailInFolder).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/folder/delete_email", folderHandler.DeleteEmailInFolder).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/folder/all_emails/{id}", folderHandler.GetAllEmailsInFolder).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/folder/allname/{id}", folderHandler.GetAllName).Methods("GET", "OPTIONS")

	logRouter.HandleFunc("/gmail/emails/incoming", emailGMailHandler.GetIncoming).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/emails/sent", emailGMailHandler.GetSent).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/emails/spam", emailGMailHandler.GetSpam).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/email/{id}", emailGMailHandler.GetById).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/email/update/{id}", emailGMailHandler.Update).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/gmail/email/delete/{id}", emailGMailHandler.Delete).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/gmail/email/send", emailGMailHandler.Send).Methods("POST", "OPTIONS")

	logRouter.HandleFunc("/gmail/drafts", emailGMailHandler.GetDrafts).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/draft/adddraft", emailGMailHandler.AddDraft).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/gmail/draft/sendDraft", emailGMailHandler.SendDraft).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/gmail/draft/{id}", emailGMailHandler.GetByIdDraft).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/draft/update/{id}", emailGMailHandler.UpdateDraft).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/gmail/draft/delete/{id}", emailGMailHandler.DeleteDraft).Methods("DELETE", "OPTIONS")

	logRouter.HandleFunc("/gmail/labels", emailGMailHandler.GetLabels).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/labels/email/{id}", emailGMailHandler.GetAllNameLabels).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/label/{name}/emails", emailGMailHandler.GetAllEmailsInLabel).Methods("GET", "OPTIONS")
	logRouter.HandleFunc("/gmail/label/create", emailGMailHandler.CreateLabel).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/gmail/label/delete/{id}", emailGMailHandler.DeleteLabel).Methods("DELETE", "OPTIONS")
	logRouter.HandleFunc("/gmail/label/update/{id}", emailGMailHandler.UpdateLabel).Methods("PUT", "OPTIONS")
	logRouter.HandleFunc("/gmail/label/add_email", emailGMailHandler.AddEmailInLabel).Methods("POST", "OPTIONS")
	logRouter.HandleFunc("/gmail/label/delete_email", emailGMailHandler.DeleteEmailInLabel).Methods("DELETE", "OPTIONS")

	return logRouter
}

// startServer starting server
func startServer(router http.Handler) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:8081", "http://89.208.223.140:8081", "http://mailhub.su:8081", "http://mailhub.su:8080", "http://localhost:8080", "http://localhost:8081", "http://89.208.223.140:8080",
			"https://127.0.0.1:8081", "https://89.208.223.140:8081", "https://mailhub.su:8081", "https://mailhub.su:8080", "https://localhost:8080", "https://localhost:8081", "https://89.208.223.140:8080",
			"https://127.0.0.1", "https://89.208.223.140", "https://mailhub.su", "https://mailhub.su", "https://localhost", "https://localhost", "https://89.208.223.140"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Csrf-Token", "Content-Type", "AuthToken"},
		ExposedHeaders:   []string{"X-Csrf-Token", "AuthToken"},
	})

	corsHandler := c.Handler(router)

	port := 8080
	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), requestid.Handler(corsHandler))
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}

// startSessionCleaner starting session cleanup
func startSessionCleaner(interval time.Duration, sessionServiceClient session_proto.SessionServiceClient) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Println("Failed to create logfile" + "log.txt")
				}
				defer f.Close()

				c := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
				ctx := context.WithValue(c, "requestID", "DeleteExpiredSessionsNULL")

				req, err := sessionServiceClient.CleanupExpiredSessions(
					metadata.NewOutgoingContext(ctx,
						metadata.New(map[string]string{"requestID": ctx.Value("requestID").(string)})),
					&session_proto.CleanupExpiredSessionsRequest{},
				)
				if err != nil {
					fmt.Printf("Error cleaning expired sessions: %v\n", err)
					return
				}
				fmt.Println(req)
			}
		}
	}()
}
