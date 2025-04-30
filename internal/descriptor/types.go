package descriptor

import (
	"fmt"
	"strings"

	"github.com/meimeitou/grpc-common/common/database"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// IsWellKnownType returns true if the provided fully qualified type name is considered 'well-known'.
func IsWellKnownType(typeName string) bool {
	_, ok := wellKnownTypeConv[typeName]
	return ok
}

// GoPackage represents a golang package.
type GoPackage struct {
	// Path is the package path to the package.
	Path string
	// Name is the package name of the package
	Name string
	// Alias is an alias of the package unique within the current invocation of gRPC-Gateway generator.
	Alias string
}

// Standard returns whether the import is a golang standard package.
func (p GoPackage) Standard() bool {
	return !strings.Contains(p.Path, ".")
}

// String returns a string representation of this package in the form of import line in golang.
func (p GoPackage) String() string {
	if p.Alias == "" {
		return fmt.Sprintf("%q", p.Path)
	}
	return fmt.Sprintf("%s %q", p.Alias, p.Path)
}

// ResponseFile wraps pluginpb.CodeGeneratorResponse_File.
type ResponseFile struct {
	*pluginpb.CodeGeneratorResponse_File
	// GoPkg is the Go package of the generated file.
	GoPkg GoPackage
}

// File wraps descriptorpb.FileDescriptorProto for richer features.
type File struct {
	*descriptorpb.FileDescriptorProto
	// GoPkg is the go package of the go file generated from this file.
	GoPkg GoPackage
	// GeneratedFilenamePrefix is used to construct filenames for generated
	// files associated with this source file.
	//
	// For example, the source file "dir/foo.proto" might have a filename prefix
	// of "dir/foo". Appending ".pb.go" produces an output file of "dir/foo.pb.go".
	GeneratedFilenamePrefix string
	// Messages is the list of messages defined in this file.
	Messages []*Message
	// Enums is the list of enums defined in this file.
	Enums []*Enum
	// Services is the list of services defined in this file.
	Services []*Service
}

// Pkg returns package name or alias if it's present
func (f *File) Pkg() string {
	pkg := f.GoPkg.Name
	if alias := f.GoPkg.Alias; alias != "" {
		pkg = alias
	}
	return pkg
}

// proto2 determines if the syntax of the file is proto2.
func (f *File) proto2() bool {
	return f.Syntax == nil || f.GetSyntax() == "proto2"
}

// Message describes a protocol buffer message types.
type Message struct {
	*descriptorpb.DescriptorProto
	// File is the file where the message is defined.
	File *File
	// Outers is a list of outer messages if this message is a nested type.
	Outers []string
	// Fields is a list of message fields.
	Fields []*Field
	// Index is proto path index of this message in File.
	Index int
	// ForcePrefixedName when set to true, prefixes a type with a package prefix.
	ForcePrefixedName bool
}

// FQMN returns a fully qualified message name of this message.
func (m *Message) FQMN() string {
	components := []string{""}
	if m.File.Package != nil {
		components = append(components, m.File.GetPackage())
	}
	components = append(components, m.Outers...)
	components = append(components, m.GetName())
	return strings.Join(components, ".")
}

// GoType returns a go type name for the message type.
// It prefixes the type name with the package alias if
// its belonging package is not "currentPackage".
func (m *Message) GoType(currentPackage string) string {
	var components []string
	components = append(components, m.Outers...)
	components = append(components, m.GetName())

	name := strings.Join(components, "_")
	if !m.ForcePrefixedName && m.File.GoPkg.Path == currentPackage {
		return name
	}
	return fmt.Sprintf("%s.%s", m.File.Pkg(), name)
}

// Enum describes a protocol buffer enum types.
type Enum struct {
	*descriptorpb.EnumDescriptorProto
	// File is the file where the enum is defined
	File *File
	// Outers is a list of outer messages if this enum is a nested type.
	Outers []string
	// Index is a enum index value.
	Index int
	// ForcePrefixedName when set to true, prefixes a type with a package prefix.
	ForcePrefixedName bool
}

// FQEN returns a fully qualified enum name of this enum.
func (e *Enum) FQEN() string {
	components := []string{""}
	if e.File.Package != nil {
		components = append(components, e.File.GetPackage())
	}
	components = append(components, e.Outers...)
	components = append(components, e.GetName())
	return strings.Join(components, ".")
}

// GoType returns a go type name for the enum type.
// It prefixes the type name with the package alias if
// its belonging package is not "currentPackage".
func (e *Enum) GoType(currentPackage string) string {
	var components []string
	components = append(components, e.Outers...)
	components = append(components, e.GetName())

	name := strings.Join(components, "_")
	if !e.ForcePrefixedName && e.File.GoPkg.Path == currentPackage {
		return name
	}
	return fmt.Sprintf("%s.%s", e.File.Pkg(), name)
}

// Service wraps descriptorpb.ServiceDescriptorProto for richer features.
type Service struct {
	*descriptorpb.ServiceDescriptorProto
	// File is the file where this service is defined.
	File *File
	// Methods is the list of methods defined in this service.
	Methods []*Method
	// ForcePrefixedName when set to true, prefixes a type with a package prefix.
	ForcePrefixedName bool
}

// FQSN returns the fully qualified service name of this service.
func (s *Service) FQSN() string {
	components := []string{""}
	if s.File.Package != nil {
		components = append(components, s.File.GetPackage())
	}
	components = append(components, s.GetName())
	return strings.Join(components, ".")
}

// InstanceName returns object name of the service with package prefix if needed
func (s *Service) InstanceName() string {
	if !s.ForcePrefixedName {
		return s.GetName()
	}
	return fmt.Sprintf("%s.%s", s.File.Pkg(), s.GetName())
}

// ClientConstructorName returns name of the Client constructor with package prefix if needed
func (s *Service) ClientConstructorName() string {
	constructor := "New" + s.GetName() + "Client"
	if !s.ForcePrefixedName {
		return constructor
	}
	return fmt.Sprintf("%s.%s", s.File.Pkg(), constructor)
}

// Method wraps descriptorpb.MethodDescriptorProto for richer features.
type Method struct {
	*descriptorpb.MethodDescriptorProto
	// Service is the service which this method belongs to.
	Service *Service
	// RequestType is the message type of requests to this method.
	RequestType *Message
	// ResponseType is the message type of responses from this method.
	ResponseType *Message
	Binding      *Binding
}

// FQMN returns a fully qualified rpc method name of this method.
func (m *Method) FQMN() string {
	var components []string
	components = append(components, m.Service.FQSN())
	components = append(components, m.GetName())
	return strings.Join(components, ".")
}

// Binding describes how an HTTP endpoint is bound to a gRPC method.
type Binding struct {
	// Method is the method which the endpoint is bound to.
	Method *Method
	// ChQuery is the query option for clickhouse
	ChQuery *database.QueryDefine
	// mysql/postgres
	Query *database.QueryDefine
}

func (b *Binding) ValidQuery() error {
	if err := b.validQueryOption(b.ChQuery); err != nil {
		return err
	}
	if err := b.validQueryOption(b.Query); err != nil {
		return err
	}
	return nil
}

func (b *Binding) validQueryOption(q *database.QueryDefine) error {

	return nil
}

// Field wraps descriptorpb.FieldDescriptorProto for richer features.
type Field struct {
	*descriptorpb.FieldDescriptorProto
	// Message is the message type which this field belongs to.
	Message *Message
	// FieldMessage is the message type of the field.
	// FieldMessage *Message
	// ForcePrefixedName when set to true, prefixes a type with a package prefix.
	ForcePrefixedName bool
	Repeated          bool
}

// FQFN returns a fully qualified field name of this field.
func (f *Field) FQFN() string {
	return strings.Join([]string{f.Message.FQMN(), f.GetName()}, ".")
}

// check field type is a well-known type
func (f *Field) IsWellKnownType() bool {
	if _, ok := wellKnownTypeConv[f.FQFN()]; ok {
		return true
	}
	return false
}

// check field type is a common type
func (f *Field) IsCommonType() bool {
	if _, ok := proto3ConvertFuncs[f.GetType()]; ok {
		return true
	}
	return false
}

var (
	proto3ConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.Bool",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.String",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:    "runtime.Bytes",
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.Enum",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64",
	}

	proto3OptionalConvertFuncs = func() map[descriptorpb.FieldDescriptorProto_Type]string {
		result := make(map[descriptorpb.FieldDescriptorProto_Type]string)
		for typ, converter := range proto3ConvertFuncs {
			// TODO: this will use convert functions from proto2.
			//       The converters returning pointers should be moved
			//       to a more generic file.
			result[typ] = converter + "P"
		}
		return result
	}()

	// TODO: replace it with a IIFE
	proto3RepeatedConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolSlice",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.StringSlice",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:    "runtime.BytesSlice",
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumSlice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64Slice",
	}

	proto2ConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64P",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32P",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64P",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64P",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32P",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64P",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32P",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolP",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.StringP",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(yugui) Handle bytes
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32P",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumP",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32P",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64P",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32P",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64P",
	}

	proto2RepeatedConvertFuncs = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolSlice",
		descriptorpb.FieldDescriptorProto_TYPE_STRING:  "runtime.StringSlice",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(maros7) Handle bytes
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumSlice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32Slice",
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64Slice",
	}

	wellKnownTypeConv = map[string]string{
		".google.protobuf.Timestamp":   "runtime.Timestamp",
		".google.protobuf.Duration":    "runtime.Duration",
		".google.protobuf.StringValue": "runtime.StringValue",
		".google.protobuf.FloatValue":  "runtime.FloatValue",
		".google.protobuf.DoubleValue": "runtime.DoubleValue",
		".google.protobuf.BoolValue":   "runtime.BoolValue",
		".google.protobuf.BytesValue":  "runtime.BytesValue",
		".google.protobuf.Int32Value":  "runtime.Int32Value",
		".google.protobuf.UInt32Value": "runtime.UInt32Value",
		".google.protobuf.Int64Value":  "runtime.Int64Value",
		".google.protobuf.UInt64Value": "runtime.UInt64Value",
	}
)

func ProtoTypeToGoType(t descriptorpb.FieldDescriptorProto_Type) string {
	switch t {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return "float64"
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "float32"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return "uint32"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return "int32" // 通常enum底层是int32
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		return "struct" // 复杂类型，后面可以再细分
	default:
		return "unknown"
	}
}
