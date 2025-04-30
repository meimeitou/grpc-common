package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/meimeitou/grpc-common/internal/codegenerator"
	"github.com/meimeitou/grpc-common/internal/descriptor"
	gindatabase "github.com/meimeitou/grpc-common/protoc-gen-database-query/gen"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	// file        = flag.String("file", "-", "where to load data from")
	outPrefix   = flag.String("output_prefix", "-", "where to save rbac auth data file")
	versionFlag = flag.Bool("version", false, "print the current version")
	logfile     = flag.String("log_file", "", "log file")
	standalone  = flag.Bool("standalone", false, "generates a standalone gateway package, which imports the target service package")
	sqlFile     = flag.String("sql_file", "", "sql file")
	// _           = flag.Bool("logtostderr", false, "Legacy glog compatibility. This flag is a no-op, you can safely remove it")
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	flag.Parse()
	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			grpclog.Fatalf("Failed to open log file %s: %v", *logfile, err)
			os.Exit(1)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	if *versionFlag {
		fmt.Printf("protoc-gen-database-query Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		reg := descriptor.NewRegistry()
		codegenerator.SetSupportedFeaturesOnPluginGen(gen)
		generator := gindatabase.New(reg, *standalone, *sqlFile)

		if err := reg.LoadFromPlugin(gen); err != nil {
			return err
		}

		targets := make([]*descriptor.File, 0, len(gen.Request.FileToGenerate))
		for _, target := range gen.Request.FileToGenerate {
			f, err := reg.LookupFile(target)
			if err != nil {
				return err
			}
			targets = append(targets, f)
		}
		files, err := generator.Generate(targets)
		if err != nil {
			return err
		}
		for _, f := range files {
			if grpclog.V(1) {
				grpclog.Infof("NewGeneratedFile %q in %s", f.GetName(), f.GoPkg)
			}
			genFile := gen.NewGeneratedFile(f.GetName(), protogen.GoImportPath(f.GoPkg.Path))
			if _, err := genFile.Write([]byte(f.GetContent())); err != nil {
				return err
			}
		}

		if grpclog.V(1) {
			grpclog.Info("Processed code generator request")
		}
		return nil
	})

	// f := os.Stdin
	// if *file != "-" {
	// 	var err error
	// 	f, err = os.Open(*file)
	// 	if err != nil {
	// 		grpclog.Fatal(err)
	// 	}
	// }

	// req, err := ParseRequest(f)
	// if err != nil {
	// 	grpclog.Fatal(err)
	// }
	// if grpclog.V(1) {
	// 	grpclog.Info("Parsed code generator request")
	// }

	// reg := NewRegistry(req.GetParameter())
	// reg.ParseProto(req)
	// reg.Output()
}

func ParseRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read code generator request: %w", err)
	}
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(input, req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code generator request: %w", err)
	}
	return req, nil
}
