package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/utils/constants"

	folder_mock "mail/internal/microservice/folder/mock"
	folder_proto "mail/internal/microservice/folder/proto"
	api "mail/internal/models/delivery_models"
	session_mock "mail/internal/pkg/session/mock"
)

func TestFolderHandler_Add_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	newFolder := api.Folder{
		ID:   1,
		Name: "Test Folder",
	}
	jsonData, _ := json.Marshal(newFolder)
	req := httptest.NewRequest("POST", "/api/v1/folder/add", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(2), nil)

	mockFolderServiceClient.EXPECT().CreateFolder(gomock.Any(), gomock.Any()).
		Return(&folder_proto.FolderWithID{Folder: &folder_proto.Folder{
			Id:   1,
			Name: "Test Folder",
		}, Id: uint32(1)}, nil)

	handler.Add(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := `{"status":200,"body":{"folder":{"id":1,"name":"Test Folder"}}}` + "\n"
	assert.Equal(t, expectedResponse, rr.Body.String())
}

func TestFolderHandler_Add_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("POST", "/api/v1/folder/add", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Add(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponse := `{"status":400,"body":{"error":"Bad JSON in request"}}` + "\n"
	assert.Equal(t, expectedResponse, rr.Body.String())
}

func TestFolderHandler_Add_BadSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	newFolder := api.Folder{
		ID:   1,
		Name: "Test Folder",
	}
	jsonData, _ := json.Marshal(newFolder)
	req := httptest.NewRequest("POST", "/api/v1/folder/add", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(0), errors.New("bad session"))

	handler.Add(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedResponse := `{"status":400,"body":{"error":"Bad session"}}` + "\n"
	assert.Equal(t, expectedResponse, rr.Body.String())
}

func TestFolderHandler_Add_FailedToAddFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	newFolder := api.Folder{
		ID:   1,
		Name: "Test Folder",
	}
	jsonData, _ := json.Marshal(newFolder)
	req := httptest.NewRequest("POST", "/api/v1/folder/add", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(2), nil)

	mockFolderServiceClient.EXPECT().CreateFolder(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("failed to add folder"))

	handler.Add(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedResponse := `{"status":500,"body":{"error":"Failed to add folder message"}}` + "\n"
	assert.Equal(t, expectedResponse, rr.Body.String())
}

func TestFolderHandler_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(2), nil)

	mockFolderServiceClient.EXPECT().GetAllFolders(gomock.Any(), gomock.Any()).
		Return(&folder_proto.Folders{
			Folders: []*folder_proto.Folder{
				{Id: 1, Name: "Folder 1"},
				{Id: 2, Name: "Folder 2"},
			},
		}, nil)

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFolderHandler_GetAll_FailedToGetProfileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(0), errors.New("failed to get profile ID"))

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_GetAll_FailedToGetFolders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().GetAllFolders(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("failed to get folders"))

	handler.GetAll(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFolderHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete/1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().DeleteFolder(gomock.Any(), gomock.Any()).
		Return(&folder_proto.FolderStatus{Status: true}, nil)

	handler.Delete(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFolderHandler_Delete_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete/not_an_integer", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Delete(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_Delete_FailedToDeleteFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete/1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().DeleteFolder(gomock.Any(), gomock.Any()).
		Return(&folder_proto.FolderStatus{Status: false}, nil)

	handler.Delete(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFolderHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	newFolder := api.Folder{
		ID:   1,
		Name: "Updated Test Folder",
	}
	jsonData, _ := json.Marshal(newFolder)
	req := httptest.NewRequest("PUT", "/api/v1/folder/update/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().UpdateFolder(gomock.Any(), gomock.Any()).
		Return(&folder_proto.FolderStatus{Status: true}, nil)

	handler.Update(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFolderHandler_Update_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"ID":"not_an_integer","Name":"Updated Folder"}`

	req := httptest.NewRequest("PUT", "/api/v1/folder/update/not_an_integer", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Update(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_Update_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("PUT", "/api/v1/folder/update/1", strings.NewReader("not_a_json_body"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Update(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_Update_FailedToUpdateFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"ID":1,"Name":"Updated Folder"}`

	req := httptest.NewRequest("PUT", "/api/v1/folder/update/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().UpdateFolder(gomock.Any(), gomock.Any()).
		Return(&folder_proto.FolderStatus{Status: false}, nil)

	handler.Update(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFolderHandler_AddEmailInFolder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("POST", "/api/v1/folder/add_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)
	mockFolderServiceClient.EXPECT().CheckEmailProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().AddEmailInFolder(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	handler.AddEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFolderHandler_AddEmailInFolder_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID:1,"InvalidField":1}`
	req := httptest.NewRequest("POST", "/api/v1/folder/add_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.AddEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_AddEmailInFolder_BadSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("POST", "/api/v1/folder/add_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(0), errors.New("session error"))

	handler.AddEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_AddEmailInFolder_CheckFolderProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("POST", "/api/v1/folder/add_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("folder profile check error"))

	handler.AddEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_AddEmailInFolder_CheckEmailProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("POST", "/api/v1/folder/add_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().CheckEmailProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("email profile check error"))

	handler.AddEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_AddEmailInFolder_AddEmailInFolderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("POST", "/api/v1/folder/add_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().CheckEmailProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().AddEmailInFolder(gomock.Any(), gomock.Any()).Return(nil, errors.New("add email in folder error"))

	handler.AddEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFolderHandler_DeleteEmailInFolder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().CheckEmailProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().DeleteEmailInFolder(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	handler.DeleteEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFolderHandler_DeleteEmailInFolder_BadJSONRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1`
	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.DeleteEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_DeleteEmailInFolder_BadSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), errors.New("session error"))

	handler.DeleteEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_DeleteEmailInFolder_CheckFolderProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("check folder profile error"))

	handler.DeleteEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_DeleteEmailInFolder_CheckEmailProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().CheckEmailProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("check email profile error"))

	handler.DeleteEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_DeleteEmailInFolder_DeleteEmailInFolderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	reqBody := `{"FolderID":1,"EmailID":1}`
	req := httptest.NewRequest("DELETE", "/api/v1/folder/delete_email", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().CheckEmailProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockFolderServiceClient.EXPECT().DeleteEmailInFolder(gomock.Any(), gomock.Any()).Return(nil, errors.New("delete email in folder error"))

	handler.DeleteEmailInFolder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestFolderHandler_GetAllEmailsInFolder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all_emails/1", nil)
	req.Header.Set("X-Csrf-Token", "token")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockSessionsManager.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).Return("loginUser", nil)

	mockFolderServiceClient.EXPECT().GetAllEmailsInFolder(gomock.Any(), gomock.Any()).Return(&folder_proto.ObjectsEmail{
		Emails: []*folder_proto.ObjectEmail{}}, nil)

	handler.GetAllEmailsInFolder(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFolderHandler_GetAllEmailsInFolder_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all_emails/abc", nil)
	req.Header.Set("X-Csrf-Token", "token")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "i"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	handler.GetAllEmailsInFolder(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_GetAllEmailsInFolder_BadSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all_emails/1", nil)
	req.Header.Set("X-Csrf-Token", "token")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(0), errors.New("session error"))

	handler.GetAllEmailsInFolder(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_GetAllEmailsInFolder_CheckFolderProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all_emails/1", nil)
	req.Header.Set("X-Csrf-Token", "token")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(0), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("check folder profile error"))

	handler.GetAllEmailsInFolder(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestFolderHandler_GetAllEmailsInFolder_GetAllEmailsInFolderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionsManager := session_mock.NewMockSessionsManager(ctrl)
	mockFolderServiceClient := folder_mock.NewMockFolderServiceClient(ctrl)

	handler := &FolderHandler{
		Sessions:            mockSessionsManager,
		FolderServiceClient: mockFolderServiceClient,
	}

	req := httptest.NewRequest("GET", "/api/v1/folder/all_emails/1", nil)
	req.Header.Set("X-Csrf-Token", "token")

	ctx := req.Context()
	ctx = context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), "testID")
	r := req.WithContext(ctx)
	vars := map[string]string{"id": "1"}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	mockSessionsManager.EXPECT().GetProfileIDBySessionID(gomock.Any(), gomock.Any()).Return(uint32(0), nil)

	mockFolderServiceClient.EXPECT().CheckFolderProfile(gomock.Any(), gomock.Any()).Return(&folder_proto.FolderEmailStatus{Status: true}, nil)

	mockSessionsManager.EXPECT().GetLoginBySession(gomock.Any(), gomock.Any()).Return("loginUser", nil)

	mockFolderServiceClient.EXPECT().GetAllEmailsInFolder(gomock.Any(), gomock.Any()).Return(nil, errors.New("get all emails in folder error"))

	handler.GetAllEmailsInFolder(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
