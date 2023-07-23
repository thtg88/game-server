// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: internal/msgs/game.proto

package msgs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Level uint64 `protobuf:"varint,2,opt,name=level,proto3" json:"level,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_msgs_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_internal_msgs_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_internal_msgs_game_proto_rawDescGZIP(), []int{0}
}

func (x *Player) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Player) GetLevel() uint64 {
	if x != nil {
		return x.Level
	}
	return 0
}

// The request message containing ...
type PlayRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Player *Player `protobuf:"bytes,1,opt,name=player,proto3" json:"player,omitempty"`
}

func (x *PlayRequest) Reset() {
	*x = PlayRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_msgs_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayRequest) ProtoMessage() {}

func (x *PlayRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_msgs_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayRequest.ProtoReflect.Descriptor instead.
func (*PlayRequest) Descriptor() ([]byte, []int) {
	return file_internal_msgs_game_proto_rawDescGZIP(), []int{1}
}

func (x *PlayRequest) GetPlayer() *Player {
	if x != nil {
		return x.Player
	}
	return nil
}

// The response message containing ...
type PlayReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PlayReply) Reset() {
	*x = PlayReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_msgs_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayReply) ProtoMessage() {}

func (x *PlayReply) ProtoReflect() protoreflect.Message {
	mi := &file_internal_msgs_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayReply.ProtoReflect.Descriptor instead.
func (*PlayReply) Descriptor() ([]byte, []int) {
	return file_internal_msgs_game_proto_rawDescGZIP(), []int{2}
}

func (x *PlayReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_internal_msgs_game_proto protoreflect.FileDescriptor

var file_internal_msgs_game_proto_rawDesc = []byte{
	0x0a, 0x18, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6d, 0x73, 0x67, 0x73, 0x2f,
	0x67, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6d, 0x73, 0x67, 0x22,
	0x2e, 0x0a, 0x06, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x22,
	0x32, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23,
	0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x6d, 0x73, 0x67, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x06, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x22, 0x25, 0x0a, 0x09, 0x50, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x34, 0x0a, 0x04, 0x47, 0x61,
	0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x04, 0x50, 0x6c, 0x61, 0x79, 0x12, 0x10, 0x2e, 0x6d, 0x73, 0x67,
	0x2e, 0x50, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x6d,
	0x73, 0x67, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x30, 0x01,
	0x42, 0x53, 0x0a, 0x17, 0x69, 0x6f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x6d, 0x73, 0x67, 0x73, 0x42, 0x09, 0x47, 0x61, 0x6d,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x74, 0x67, 0x38, 0x38, 0x2f, 0x67, 0x61, 0x6d, 0x65,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x6d, 0x73, 0x67, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_msgs_game_proto_rawDescOnce sync.Once
	file_internal_msgs_game_proto_rawDescData = file_internal_msgs_game_proto_rawDesc
)

func file_internal_msgs_game_proto_rawDescGZIP() []byte {
	file_internal_msgs_game_proto_rawDescOnce.Do(func() {
		file_internal_msgs_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_msgs_game_proto_rawDescData)
	})
	return file_internal_msgs_game_proto_rawDescData
}

var file_internal_msgs_game_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_msgs_game_proto_goTypes = []interface{}{
	(*Player)(nil),      // 0: msg.Player
	(*PlayRequest)(nil), // 1: msg.PlayRequest
	(*PlayReply)(nil),   // 2: msg.PlayReply
}
var file_internal_msgs_game_proto_depIdxs = []int32{
	0, // 0: msg.PlayRequest.player:type_name -> msg.Player
	1, // 1: msg.Game.Play:input_type -> msg.PlayRequest
	2, // 2: msg.Game.Play:output_type -> msg.PlayReply
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_msgs_game_proto_init() }
func file_internal_msgs_game_proto_init() {
	if File_internal_msgs_game_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_msgs_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Player); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_msgs_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_msgs_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_msgs_game_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_msgs_game_proto_goTypes,
		DependencyIndexes: file_internal_msgs_game_proto_depIdxs,
		MessageInfos:      file_internal_msgs_game_proto_msgTypes,
	}.Build()
	File_internal_msgs_game_proto = out.File
	file_internal_msgs_game_proto_rawDesc = nil
	file_internal_msgs_game_proto_goTypes = nil
	file_internal_msgs_game_proto_depIdxs = nil
}
