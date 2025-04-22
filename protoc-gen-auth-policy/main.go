package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/meimeitou/grpc-common/common/auth"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	// Version is the version of the plugin.
	Version = "1.0.0"
)

var (
	file      = flag.String("file", "-", "where to load data from")
	outPrefix = flag.String("output_prefix", "-", "where to save rbac auth data file")
	version   = flag.Bool("version", false, "print the version and exit")
	_         = flag.Bool("logtostderr", false, "Legacy glog compatibility. This flag is a no-op, you can safely remove it")
)

func main() {
	flag.Parse()
	f := os.Stdin
	if *file != "-" {
		var err error
		f, err = os.Open(*file)
		if err != nil {
			grpclog.Fatal(err)
		}
	}
	if *version {
		fmt.Printf("protoc-gen-auth-policy version: %s\n", Version)
		return
	}
	req, err := ParseRequest(f)
	if err != nil {
		grpclog.Fatal(err)
	}
	if grpclog.V(1) {
		grpclog.Info("Parsed code generator request")
	}

	reg := NewRegistry(req.GetParameter())
	reg.ParseProto(req)
	reg.Output()
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

type Registry struct {
	// Services is a map from service name to service.
	outPrefix string // output file name
	resp      *pluginpb.CodeGeneratorResponse
	authRbac  map[string]*bytes.Buffer
}

func NewRegistry(param string) *Registry {
	r := &Registry{
		resp:     &pluginpb.CodeGeneratorResponse{},
		authRbac: make(map[string]*bytes.Buffer),
	}
	r.parseReqParam(param)
	return r
}

func (r *Registry) parseReqParam(param string) {
	if param == "" {
		return
	}
	for _, p := range strings.Split(param, ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 2 {
			if spec[0] == "output_prefix" {
				r.outPrefix = spec[1]
			}
		}
	}
	return
}

func (r *Registry) ParseProto(req *pluginpb.CodeGeneratorRequest) {
	fileMap := make(map[string]*descriptorpb.FileDescriptorProto)
	for _, file := range req.GetProtoFile() {
		fileMap[file.GetName()] = file
	}
	for _, fileName := range req.GetFileToGenerate() {
		file := fileMap[fileName]
		if len(file.Service) == 0 {
			continue
		}
		// Process services
		for _, service := range file.Service {
			// fmt.Fprintf(tw, "Service: %s\n", service.GetName())

			// Process methods
			for _, method := range service.GetMethod() {
				// fmt.Fprintf(tw, "  Method: %s\n", method.GetName())
				opts, err := extractAPIOptions(method)
				if err != nil {
					grpclog.Errorf("Failed to extract HttpRule from %s.%s: %v", service.GetName(), method.GetName(), err)
					continue
				}
				if opts == nil {
					continue
				}
				// fmt.Fprintf(tw, "    Options:\n")
				// fmt.Fprintf(tw, "      Roles: %v\n", opts.GetRole())
				for _, p := range opts.GetRoles() {
					key := fmt.Sprintf("%s_%s", service.GetName(), p.GetCategory())
					if _, ok := r.authRbac[key]; !ok {
						buf := new(bytes.Buffer)
						r.authRbac[key] = buf
					}
					for _, role := range p.GetNames() {
						fmt.Fprintf(r.authRbac[key], "p, %s, %s\n",
							role, fmt.Sprintf("/%s.%s/%s",
								file.GetPackage(), service.GetName(), method.GetName()))
					}
				}
			}
		}
		// tw.Flush()
		for category, buf := range r.authRbac {
			fileName := fmt.Sprintf("%s.csv", category)
			if r.outPrefix != "" {
				fileName = fmt.Sprintf("%s_%s.csv", r.outPrefix, category)
			}
			r.resp.File = append(r.resp.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(fileName),
				Content: proto.String(buf.String()),
			})
		}
	}
}

func (r *Registry) Output() {
	output, err := proto.Marshal(r.resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to serialize response: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
		os.Exit(1)
	}
}

func extractAPIOptions(meth *descriptorpb.MethodDescriptorProto) (*auth.RBACAuth, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, auth.E_RbacAuth) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, auth.E_RbacAuth)
	opts, ok := ext.(*auth.RBACAuth)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}
	return opts, nil
}
