package gen

import (
	"errors"
	"go/format"
	"log"

	"github.com/meimeitou/grpc-common/internal/descriptor"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var errNoTargetService = errors.New("no target service defined in the file")

type generator struct {
	reg         *descriptor.Registry
	baseImports []descriptor.GoPackage
	standalone  bool
	sqlFile     string
}

func New(reg *descriptor.Registry, standalone bool, sqlFile string) *generator {
	var imports []descriptor.GoPackage
	return &generator{
		reg:         reg,
		baseImports: imports,
		standalone:  standalone,
		sqlFile:     sqlFile,
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	var files []*descriptor.ResponseFile
	for _, file := range targets {
		if grpclog.V(1) {
			grpclog.Infof("Processing %s", file.GetName())
		}

		code, err := g.generate(file)
		if errors.Is(err, errNoTargetService) {
			if grpclog.V(1) {
				grpclog.Infof("%s: %v", file.GetName(), err)
			}
			continue
		}
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(code))
		if err != nil {
			grpclog.Errorf("%v: %s", err, code)
			return nil, err
		}
		files = append(files, &descriptor.ResponseFile{
			GoPkg: file.GoPkg,
			CodeGeneratorResponse_File: &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(file.GeneratedFilenamePrefix + ".db.go"),
				Content: proto.String(string(formatted)),
			},
		})
	}
	return files, nil
}

func (g *generator) generate(file *descriptor.File) (string, error) {
	for _, svc := range file.Services {
		for _, m := range svc.Methods {
			// log.Printf("method %s", m.GetName())
			// log.Printf("input type %s", m.RequestType.GetName())
			for _, f := range m.RequestType.Fields {
				log.Printf("input fields: %s, type: %s", f.GetName(), f.GetType())
			}
			log.Printf("output type %s", m.ResponseType.GetName())
			for _, f := range m.ResponseType.Fields {
				log.Printf("output fields: %s, type: %s, repeated: %v, common type: %v", f.GetName(), f.GetType(), f.Repeated, f.IsCommonType())
				if f.Repeated {
					log.Printf("repeated field %s", f.GetTypeName())
				}
			}
			log.Printf("binding %v", m.Binding)
		}
	}
	return "", errNoTargetService
}
