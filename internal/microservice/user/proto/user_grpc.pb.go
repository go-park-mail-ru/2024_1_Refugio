// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: user.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersReply, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserReply, error)
	GetUserByLogin(ctx context.Context, in *GetUserByLoginRequest, opts ...grpc.CallOption) (*GetUserByLoginReply, error)
	IsLoginUnique(ctx context.Context, in *IsLoginUniqueRequest, opts ...grpc.CallOption) (*IsLoginUniqueReply, error)
	DeleteUserById(ctx context.Context, in *DeleteUserByIdRequest, opts ...grpc.CallOption) (*DeleteUserByIdReply, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserReply, error)
	UploadUserAvatar(ctx context.Context, in *UploadUserAvatarRequest, opts ...grpc.CallOption) (*UploadUserAvatarReply, error)
	DeleteUserAvatar(ctx context.Context, in *DeleteUserAvatarRequest, opts ...grpc.CallOption) (*DeleteUserAvatarReply, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserReply, error)
	GetUserByVKId(ctx context.Context, in *GetUserVKIdRequest, opts ...grpc.CallOption) (*GetUserReply, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersReply, error) {
	out := new(GetUsersReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/GetUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserReply, error) {
	out := new(GetUserReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserByLogin(ctx context.Context, in *GetUserByLoginRequest, opts ...grpc.CallOption) (*GetUserByLoginReply, error) {
	out := new(GetUserByLoginReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/GetUserByLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) IsLoginUnique(ctx context.Context, in *IsLoginUniqueRequest, opts ...grpc.CallOption) (*IsLoginUniqueReply, error) {
	out := new(IsLoginUniqueReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/IsLoginUnique", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUserById(ctx context.Context, in *DeleteUserByIdRequest, opts ...grpc.CallOption) (*DeleteUserByIdReply, error) {
	out := new(DeleteUserByIdReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/DeleteUserById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserReply, error) {
	out := new(UpdateUserReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UploadUserAvatar(ctx context.Context, in *UploadUserAvatarRequest, opts ...grpc.CallOption) (*UploadUserAvatarReply, error) {
	out := new(UploadUserAvatarReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/UploadUserAvatar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUserAvatar(ctx context.Context, in *DeleteUserAvatarRequest, opts ...grpc.CallOption) (*DeleteUserAvatarReply, error) {
	out := new(DeleteUserAvatarReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/DeleteUserAvatar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserReply, error) {
	out := new(CreateUserReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserByVKId(ctx context.Context, in *GetUserVKIdRequest, opts ...grpc.CallOption) (*GetUserReply, error) {
	out := new(GetUserReply)
	err := c.cc.Invoke(ctx, "/proto.UserService/GetUserByVKId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersReply, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserReply, error)
	GetUserByLogin(context.Context, *GetUserByLoginRequest) (*GetUserByLoginReply, error)
	IsLoginUnique(context.Context, *IsLoginUniqueRequest) (*IsLoginUniqueReply, error)
	DeleteUserById(context.Context, *DeleteUserByIdRequest) (*DeleteUserByIdReply, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserReply, error)
	UploadUserAvatar(context.Context, *UploadUserAvatarRequest) (*UploadUserAvatarReply, error)
	DeleteUserAvatar(context.Context, *DeleteUserAvatarRequest) (*DeleteUserAvatarReply, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error)
	GetUserByVKId(context.Context, *GetUserVKIdRequest) (*GetUserReply, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) GetUsers(context.Context, *GetUsersRequest) (*GetUsersReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServiceServer) GetUserByLogin(context.Context, *GetUserByLoginRequest) (*GetUserByLoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByLogin not implemented")
}
func (UnimplementedUserServiceServer) IsLoginUnique(context.Context, *IsLoginUniqueRequest) (*IsLoginUniqueReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsLoginUnique not implemented")
}
func (UnimplementedUserServiceServer) DeleteUserById(context.Context, *DeleteUserByIdRequest) (*DeleteUserByIdReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserById not implemented")
}
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) UploadUserAvatar(context.Context, *UploadUserAvatarRequest) (*UploadUserAvatarReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadUserAvatar not implemented")
}
func (UnimplementedUserServiceServer) DeleteUserAvatar(context.Context, *DeleteUserAvatarRequest) (*DeleteUserAvatarReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserAvatar not implemented")
}
func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) GetUserByVKId(context.Context, *GetUserVKIdRequest) (*GetUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByVKId not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/GetUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUsers(ctx, req.(*GetUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserByLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/GetUserByLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByLogin(ctx, req.(*GetUserByLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_IsLoginUnique_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsLoginUniqueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).IsLoginUnique(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/IsLoginUnique",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).IsLoginUnique(ctx, req.(*IsLoginUniqueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUserById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUserById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/DeleteUserById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUserById(ctx, req.(*DeleteUserByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UploadUserAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadUserAvatarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UploadUserAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/UploadUserAvatar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UploadUserAvatar(ctx, req.(*UploadUserAvatarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUserAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserAvatarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUserAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/DeleteUserAvatar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUserAvatar(ctx, req.(*DeleteUserAvatarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserByVKId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserVKIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByVKId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UserService/GetUserByVKId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByVKId(ctx, req.(*GetUserVKIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUsers",
			Handler:    _UserService_GetUsers_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserService_GetUser_Handler,
		},
		{
			MethodName: "GetUserByLogin",
			Handler:    _UserService_GetUserByLogin_Handler,
		},
		{
			MethodName: "IsLoginUnique",
			Handler:    _UserService_IsLoginUnique_Handler,
		},
		{
			MethodName: "DeleteUserById",
			Handler:    _UserService_DeleteUserById_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "UploadUserAvatar",
			Handler:    _UserService_UploadUserAvatar_Handler,
		},
		{
			MethodName: "DeleteUserAvatar",
			Handler:    _UserService_DeleteUserAvatar_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "GetUserByVKId",
			Handler:    _UserService_GetUserByVKId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
