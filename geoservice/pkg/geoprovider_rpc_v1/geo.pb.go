// rpc contract declaration, used to generate code
// for implementation in specific language

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: geo.proto

package geoprovider_rpc_v1

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

type AddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *AddressRequest) Reset() {
	*x = AddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_geo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressRequest) ProtoMessage() {}

func (x *AddressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_geo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressRequest.ProtoReflect.Descriptor instead.
func (*AddressRequest) Descriptor() ([]byte, []int) {
	return file_geo_proto_rawDescGZIP(), []int{0}
}

func (x *AddressRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type GeoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lat string `protobuf:"bytes,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lng string `protobuf:"bytes,2,opt,name=lng,proto3" json:"lng,omitempty"`
}

func (x *GeoRequest) Reset() {
	*x = GeoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_geo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeoRequest) ProtoMessage() {}

func (x *GeoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_geo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeoRequest.ProtoReflect.Descriptor instead.
func (*GeoRequest) Descriptor() ([]byte, []int) {
	return file_geo_proto_rawDescGZIP(), []int{1}
}

func (x *GeoRequest) GetLat() string {
	if x != nil {
		return x.Lat
	}
	return ""
}

func (x *GeoRequest) GetLng() string {
	if x != nil {
		return x.Lng
	}
	return ""
}

type AddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	City   string `protobuf:"bytes,1,opt,name=city,proto3" json:"city,omitempty"`
	Street string `protobuf:"bytes,2,opt,name=street,proto3" json:"street,omitempty"`
	House  string `protobuf:"bytes,3,opt,name=house,proto3" json:"house,omitempty"`
	Lat    string `protobuf:"bytes,4,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon    string `protobuf:"bytes,5,opt,name=lon,proto3" json:"lon,omitempty"`
}

func (x *AddressResponse) Reset() {
	*x = AddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_geo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressResponse) ProtoMessage() {}

func (x *AddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_geo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressResponse.ProtoReflect.Descriptor instead.
func (*AddressResponse) Descriptor() ([]byte, []int) {
	return file_geo_proto_rawDescGZIP(), []int{2}
}

func (x *AddressResponse) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *AddressResponse) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

func (x *AddressResponse) GetHouse() string {
	if x != nil {
		return x.House
	}
	return ""
}

func (x *AddressResponse) GetLat() string {
	if x != nil {
		return x.Lat
	}
	return ""
}

func (x *AddressResponse) GetLon() string {
	if x != nil {
		return x.Lon
	}
	return ""
}

type AddressesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses []*AddressResponse `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *AddressesResponse) Reset() {
	*x = AddressesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_geo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressesResponse) ProtoMessage() {}

func (x *AddressesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_geo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressesResponse.ProtoReflect.Descriptor instead.
func (*AddressesResponse) Descriptor() ([]byte, []int) {
	return file_geo_proto_rawDescGZIP(), []int{3}
}

func (x *AddressesResponse) GetAddresses() []*AddressResponse {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_geo_proto protoreflect.FileDescriptor

var file_geo_proto_rawDesc = []byte{
	0x0a, 0x09, 0x67, 0x65, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x67, 0x65, 0x6f,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x31, 0x22,
	0x26, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x22, 0x30, 0x0a, 0x0a, 0x47, 0x65, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6e, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x6e, 0x67, 0x22, 0x77, 0x0a, 0x0f, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x68, 0x6f, 0x75, 0x73,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x61, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c,
	0x6f, 0x6e, 0x22, 0x56, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x67, 0x65, 0x6f,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x31, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x32, 0xbd, 0x01, 0x0a, 0x0d, 0x47,
	0x65, 0x6f, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x56, 0x31, 0x12, 0x5a, 0x0a, 0x0d,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x22, 0x2e,
	0x67, 0x65, 0x6f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f,
	0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x25, 0x2e, 0x67, 0x65, 0x6f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f,
	0x72, 0x70, 0x63, 0x5f, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x07, 0x47, 0x65, 0x6f, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x1e, 0x2e, 0x67, 0x65, 0x6f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x67, 0x65, 0x6f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x65,
	0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6f,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x31, 0x3b,
	0x67, 0x65, 0x6f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x70, 0x63, 0x5f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_geo_proto_rawDescOnce sync.Once
	file_geo_proto_rawDescData = file_geo_proto_rawDesc
)

func file_geo_proto_rawDescGZIP() []byte {
	file_geo_proto_rawDescOnce.Do(func() {
		file_geo_proto_rawDescData = protoimpl.X.CompressGZIP(file_geo_proto_rawDescData)
	})
	return file_geo_proto_rawDescData
}

var file_geo_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_geo_proto_goTypes = []any{
	(*AddressRequest)(nil),    // 0: geoprovider_rpc_v1.AddressRequest
	(*GeoRequest)(nil),        // 1: geoprovider_rpc_v1.GeoRequest
	(*AddressResponse)(nil),   // 2: geoprovider_rpc_v1.AddressResponse
	(*AddressesResponse)(nil), // 3: geoprovider_rpc_v1.AddressesResponse
}
var file_geo_proto_depIdxs = []int32{
	2, // 0: geoprovider_rpc_v1.AddressesResponse.addresses:type_name -> geoprovider_rpc_v1.AddressResponse
	0, // 1: geoprovider_rpc_v1.GeoProviderV1.AddressSearch:input_type -> geoprovider_rpc_v1.AddressRequest
	1, // 2: geoprovider_rpc_v1.GeoProviderV1.GeoCode:input_type -> geoprovider_rpc_v1.GeoRequest
	3, // 3: geoprovider_rpc_v1.GeoProviderV1.AddressSearch:output_type -> geoprovider_rpc_v1.AddressesResponse
	3, // 4: geoprovider_rpc_v1.GeoProviderV1.GeoCode:output_type -> geoprovider_rpc_v1.AddressesResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_geo_proto_init() }
func file_geo_proto_init() {
	if File_geo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_geo_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AddressRequest); i {
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
		file_geo_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GeoRequest); i {
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
		file_geo_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*AddressResponse); i {
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
		file_geo_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*AddressesResponse); i {
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
			RawDescriptor: file_geo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_geo_proto_goTypes,
		DependencyIndexes: file_geo_proto_depIdxs,
		MessageInfos:      file_geo_proto_msgTypes,
	}.Build()
	File_geo_proto = out.File
	file_geo_proto_rawDesc = nil
	file_geo_proto_goTypes = nil
	file_geo_proto_depIdxs = nil
}
