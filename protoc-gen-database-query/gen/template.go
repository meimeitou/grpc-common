package gen

import (
	"fmt"

	"github.com/meimeitou/grpc-common/common/database"
	"github.com/meimeitou/grpc-common/internal/descriptor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func getMethodDatabaseOption(reg *descriptor.Registry, meth *descriptor.Method) (*database.Query, error) {
	opts, err := extractDatabaseOptionFromMethodDescriptor(meth.MethodDescriptorProto)
	if err != nil {
		return nil, err
	}
	if opts != nil {
		return opts, nil
	}
	opts, ok := reg.GetDatabaseMethodOption(meth.FQMN())
	if !ok {
		return nil, nil
	}
	return opts, nil
}

func extractDatabaseOptionFromMethodDescriptor(meth *descriptorpb.MethodDescriptorProto) (*database.Query, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, database.E_Query) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, database.E_Query)
	opts, ok := ext.(*database.Query)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an Operation", ext)
	}
	return opts, nil
}
