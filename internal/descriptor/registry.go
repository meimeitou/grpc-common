package descriptor

import (
	"fmt"
	"sort"
	"strings"

	"github.com/golang/glog"
	"github.com/meimeitou/grpc-common/common/database"
	"github.com/meimeitou/grpc-common/internal/codegenerator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Registry is a registry of information extracted from pluginpb.CodeGeneratorRequest.
type Registry struct {
	// msgs is a mapping from fully-qualified message name to descriptor
	msgs map[string]*Message

	// enums is a mapping from fully-qualified enum name to descriptor
	enums map[string]*Enum

	// files is a mapping from file path to descriptor
	files map[string]*File

	// prefix is a prefix to be inserted to golang package paths generated from proto package names.
	prefix string

	// pkgMap is a user-specified mapping from file path to proto package.
	pkgMap map[string]string

	// pkgAliases is a mapping from package aliases to package paths in go which are already taken.
	pkgAliases map[string]string

	// repeatedPathParamSeparator specifies how path parameter repeated fields are separated
	repeatedPathParamSeparator repeatedFieldSeparator

	// useGoTemplate determines whether you want to use GO templates
	// in your protofile comments
	useGoTemplate bool

	// ignoreComments determines whether all protofile comments should be excluded from output
	ignoreComments bool

	// disableDefaultErrors disables the generation of the default error types.
	// This is useful for users who have defined custom error handling.
	disableDefaultErrors bool

	// simpleOperationIDs removes the service prefix from the generated
	// operationIDs. This risks generating duplicate operationIDs.
	simpleOperationIDs bool

	standalone bool
	// warnOnUnboundMethods causes the registry to emit warning logs if an RPC method
	// has no HttpRule annotation.
	warnOnUnboundMethods bool

	// proto3OptionalNullable specifies whether Proto3 Optional fields should be marked as x-nullable.
	proto3OptionalNullable bool

	// fileOptions is a mapping of file name to additional OpenAPI file options
	databaseOptions map[string]*database.Query

	// omitPackageDoc, if false, causes a package comment to be included in the generated code.
	omitPackageDoc bool

	// annotationMap is used to check for duplicate HTTP annotations
	annotationMap map[annotationIdentifier]struct{}
}

type repeatedFieldSeparator struct {
	name string
	sep  rune
}

type annotationIdentifier struct {
	method       string
	pathTemplate string
	service      *Service
}

// NewRegistry returns a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		msgs:       make(map[string]*Message),
		enums:      make(map[string]*Enum),
		files:      make(map[string]*File),
		pkgMap:     make(map[string]string),
		pkgAliases: make(map[string]string),
		repeatedPathParamSeparator: repeatedFieldSeparator{
			name: "csv",
			sep:  ',',
		},
		databaseOptions: make(map[string]*database.Query),
		annotationMap:   make(map[annotationIdentifier]struct{}),
	}
}

// Load loads definitions of services, methods, messages, enumerations and fields from "req".
func (r *Registry) Load(req *pluginpb.CodeGeneratorRequest) error {
	gen, err := protogen.Options{}.New(req)
	if err != nil {
		return err
	}
	// Note: keep in mind that this might be not enough because
	// protogen.Plugin is used only to load files here.
	// The support for features must be set on the pluginpb.CodeGeneratorResponse.
	codegenerator.SetSupportedFeaturesOnPluginGen(gen)
	return r.load(gen)
}

func (r *Registry) LoadFromPlugin(gen *protogen.Plugin) error {
	return r.load(gen)
}

func (r *Registry) load(gen *protogen.Plugin) error {
	filePaths := make([]string, 0, len(gen.FilesByPath))
	for filePath := range gen.FilesByPath {
		filePaths = append(filePaths, filePath)
	}
	sort.Strings(filePaths)

	for _, filePath := range filePaths {
		r.loadFile(filePath, gen.FilesByPath[filePath])
	}

	for _, filePath := range filePaths {
		if !gen.FilesByPath[filePath].Generate {
			continue
		}
		file := r.files[filePath]
		if err := r.loadServices(file); err != nil {
			return err
		}
	}

	return nil
}

// loadFile loads messages, enumerations and fields from "file".
// It does not loads services and methods in "file".  You need to call
// loadServices after loadFiles is called for all files to load services and methods.
func (r *Registry) loadFile(filePath string, file *protogen.File) {
	pkg := GoPackage{
		Path: string(file.GoImportPath),
		Name: string(file.GoPackageName),
	}
	if r.standalone {
		pkg.Alias = "ext" + cases.Title(language.AmericanEnglish).String(pkg.Name)
	}

	if err := r.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
		for i := 0; ; i++ {
			alias := fmt.Sprintf("%s_%d", pkg.Name, i)
			if err := r.ReserveGoPackageAlias(alias, pkg.Path); err == nil {
				pkg.Alias = alias
				break
			}
		}
	}
	f := &File{
		FileDescriptorProto:     file.Proto,
		GoPkg:                   pkg,
		GeneratedFilenamePrefix: file.GeneratedFilenamePrefix,
	}

	r.files[filePath] = f
	r.registerMsg(f, nil, file.Proto.MessageType)
	r.registerEnum(f, nil, file.Proto.EnumType)
}

func (r *Registry) registerMsg(file *File, outerPath []string, msgs []*descriptorpb.DescriptorProto) {
	for i, md := range msgs {
		m := &Message{
			File:              file,
			Outers:            outerPath,
			DescriptorProto:   md,
			Index:             i,
			ForcePrefixedName: r.standalone,
		}
		for _, fd := range md.GetField() {
			fdata := &Field{
				Message:              m,
				FieldDescriptorProto: fd,
				ForcePrefixedName:    r.standalone,
				Repeated:             fd.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED,
			}
			m.Fields = append(m.Fields, fdata)
		}
		file.Messages = append(file.Messages, m)
		r.msgs[m.FQMN()] = m
		glog.V(1).Infof("register name: %s", m.FQMN())

		var outers []string
		outers = append(outers, outerPath...)
		outers = append(outers, m.GetName())
		r.registerMsg(file, outers, m.GetNestedType())
		r.registerEnum(file, outers, m.GetEnumType())
	}
}

func (r *Registry) registerEnum(file *File, outerPath []string, enums []*descriptorpb.EnumDescriptorProto) {
	for i, ed := range enums {
		e := &Enum{
			File:                file,
			Outers:              outerPath,
			EnumDescriptorProto: ed,
			Index:               i,
			ForcePrefixedName:   r.standalone,
		}
		file.Enums = append(file.Enums, e)
		r.enums[e.FQEN()] = e
		glog.V(1).Infof("register enum name: %s", e.FQEN())
	}
}

// LookupMsg looks up a message type by "name".
// It tries to resolve "name" from "location" if "name" is a relative message name.
func (r *Registry) LookupMsg(location, name string) (*Message, error) {
	glog.V(1).Infof("lookup %s from %s", name, location)
	if strings.HasPrefix(name, ".") {
		m, ok := r.msgs[name]
		if !ok {
			return nil, fmt.Errorf("no message found: %s", name)
		}
		return m, nil
	}

	if !strings.HasPrefix(location, ".") {
		location = fmt.Sprintf(".%s", location)
	}
	components := strings.Split(location, ".")
	for len(components) > 0 {
		fqmn := strings.Join(append(components, name), ".")
		if m, ok := r.msgs[fqmn]; ok {
			return m, nil
		}
		components = components[:len(components)-1]
	}
	return nil, fmt.Errorf("no message found: %s", name)
}

// LookupEnum looks up a enum type by "name".
// It tries to resolve "name" from "location" if "name" is a relative enum name.
func (r *Registry) LookupEnum(location, name string) (*Enum, error) {
	glog.V(1).Infof("lookup enum %s from %s", name, location)
	if strings.HasPrefix(name, ".") {
		e, ok := r.enums[name]
		if !ok {
			return nil, fmt.Errorf("no enum found: %s", name)
		}
		return e, nil
	}

	if !strings.HasPrefix(location, ".") {
		location = fmt.Sprintf(".%s", location)
	}
	components := strings.Split(location, ".")
	for len(components) > 0 {
		fqen := strings.Join(append(components, name), ".")
		if e, ok := r.enums[fqen]; ok {
			return e, nil
		}
		components = components[:len(components)-1]
	}
	return nil, fmt.Errorf("no enum found: %s", name)
}

// LookupFile looks up a file by name.
func (r *Registry) LookupFile(name string) (*File, error) {
	f, ok := r.files[name]
	if !ok {
		return nil, fmt.Errorf("no such file given: %s", name)
	}
	return f, nil
}

// AddPkgMap adds a mapping from a .proto file to proto package name.
func (r *Registry) AddPkgMap(file, protoPkg string) {
	r.pkgMap[file] = protoPkg
}

// SetPrefix registers the prefix to be added to go package paths generated from proto package names.
func (r *Registry) SetPrefix(prefix string) {
	r.prefix = prefix
}

// SetStandalone registers standalone flag to control package prefix
func (r *Registry) SetStandalone(standalone bool) {
	r.standalone = standalone
}

// ReserveGoPackageAlias reserves the unique alias of go package.
// If succeeded, the alias will be never used for other packages in generated go files.
// If failed, the alias is already taken by another package, so you need to use another
// alias for the package in your go files.
func (r *Registry) ReserveGoPackageAlias(alias, pkgpath string) error {
	if taken, ok := r.pkgAliases[alias]; ok {
		if taken == pkgpath {
			return nil
		}
		return fmt.Errorf("package name %s is already taken. Use another alias", alias)
	}
	r.pkgAliases[alias] = pkgpath
	return nil
}

// GetAllFQMNs returns a list of all FQMNs
func (r *Registry) GetAllFQMNs() []string {
	keys := make([]string, 0, len(r.msgs))
	for k := range r.msgs {
		keys = append(keys, k)
	}
	return keys
}

// GetAllFQENs returns a list of all FQENs
func (r *Registry) GetAllFQENs() []string {
	keys := make([]string, 0, len(r.enums))
	for k := range r.enums {
		keys = append(keys, k)
	}
	return keys
}

// GetRepeatedPathParamSeparator returns a rune spcifying how
// path parameter repeated fields are separated.
func (r *Registry) GetRepeatedPathParamSeparator() rune {
	return r.repeatedPathParamSeparator.sep
}

// GetRepeatedPathParamSeparatorName returns the name path parameter repeated
// fields repeatedFieldSeparator. I.e. 'csv', 'pipe', 'ssv' or 'tsv'
func (r *Registry) GetRepeatedPathParamSeparatorName() string {
	return r.repeatedPathParamSeparator.name
}

// SetRepeatedPathParamSeparator sets how path parameter repeated fields are
// separated. Allowed names are 'csv', 'pipe', 'ssv' and 'tsv'.
func (r *Registry) SetRepeatedPathParamSeparator(name string) error {
	var sep rune
	switch name {
	case "csv":
		sep = ','
	case "pipes":
		sep = '|'
	case "ssv":
		sep = ' '
	case "tsv":
		sep = '\t'
	default:
		return fmt.Errorf("unknown repeated path parameter separator: %s", name)
	}
	r.repeatedPathParamSeparator = repeatedFieldSeparator{
		name: name,
		sep:  sep,
	}
	return nil
}

// SetUseGoTemplate sets useGoTemplate
func (r *Registry) SetUseGoTemplate(use bool) {
	r.useGoTemplate = use
}

// GetUseGoTemplate returns useGoTemplate
func (r *Registry) GetUseGoTemplate() bool {
	return r.useGoTemplate
}

// SetIgnoreComments sets ignoreComments
func (r *Registry) SetIgnoreComments(ignore bool) {
	r.ignoreComments = ignore
}

// GetIgnoreComments returns ignoreComments
func (r *Registry) GetIgnoreComments() bool {
	return r.ignoreComments
}

// SetDisableDefaultErrors sets disableDefaultErrors
func (r *Registry) SetDisableDefaultErrors(use bool) {
	r.disableDefaultErrors = use
}

// GetDisableDefaultErrors returns disableDefaultErrors
func (r *Registry) GetDisableDefaultErrors() bool {
	return r.disableDefaultErrors
}

// SetSimpleOperationIDs sets simpleOperationIDs
func (r *Registry) SetSimpleOperationIDs(use bool) {
	r.simpleOperationIDs = use
}

// GetSimpleOperationIDs returns simpleOperationIDs
func (r *Registry) GetSimpleOperationIDs() bool {
	return r.simpleOperationIDs
}

// SetWarnOnUnboundMethods sets warnOnUnboundMethods
func (r *Registry) SetWarnOnUnboundMethods(warn bool) {
	r.warnOnUnboundMethods = warn
}

// SetOmitPackageDoc controls whether the generated code contains a package comment (if set to false, it will contain one)
func (r *Registry) SetOmitPackageDoc(omit bool) {
	r.omitPackageDoc = omit
}

// GetOmitPackageDoc returns whether a package comment will be omitted from the generated code
func (r *Registry) GetOmitPackageDoc() bool {
	return r.omitPackageDoc
}

// SetProto3OptionalNullable set proto3OtionalNullable
func (r *Registry) SetProto3OptionalNullable(proto3OtionalNullable bool) {
	r.proto3OptionalNullable = proto3OtionalNullable
}

// GetProto3OptionalNullable returns proto3OtionalNullable
func (r *Registry) GetProto3OptionalNullable() bool {
	return r.proto3OptionalNullable
}

func (r *Registry) FieldName(f *Field) string {
	return f.GetName()
}

func (r *Registry) GetDatabaseMethodOption(qualifiedMethod string) (*database.Query, bool) {
	opt, ok := r.databaseOptions[qualifiedMethod]
	return opt, ok
}
