package descriptor

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/meimeitou/grpc-common/common/database"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// loadServices registers services and their methods from "targetFile" to "r".
// It must be called after loadFile is called for all files so that loadServices
// can resolve names of message types and their fields.
func (r *Registry) loadServices(file *File) error {
	glog.V(1).Infof("Loading services from %s", file.GetName())
	var svcs []*Service
	for _, sd := range file.GetService() {
		glog.V(2).Infof("Registering %s", sd.GetName())
		svc := &Service{
			File:                   file,
			ServiceDescriptorProto: sd,
			ForcePrefixedName:      r.standalone,
		}
		for _, md := range sd.GetMethod() {
			glog.V(2).Infof("Processing %s.%s", sd.GetName(), md.GetName())

			opt, err := extractDatabaseOptions(md)
			if err != nil {
				glog.Errorf("Failed to extract HttpRule from %s.%s: %v", svc.GetName(), md.GetName(), err)
				return err
			}
			if opt == nil {
				glog.V(2).Infof("No database option for %s.%s", svc.GetName(), md.GetName())
				continue
			}
			meth, err := r.newMethod(svc, md, opt)
			if err != nil {
				return err
			}
			svc.Methods = append(svc.Methods, meth)
		}
		if len(svc.Methods) == 0 {
			continue
		}
		glog.V(2).Infof("Registered %s with %d method(s)", svc.GetName(), len(svc.Methods))
		svcs = append(svcs, svc)
	}
	file.Services = svcs
	return nil
}

func (r *Registry) newMethod(svc *Service, md *descriptorpb.MethodDescriptorProto, opt *database.Query) (*Method, error) {
	requestType, err := r.LookupMsg(svc.File.GetPackage(), md.GetInputType())
	if err != nil {
		return nil, err
	}
	responseType, err := r.LookupMsg(svc.File.GetPackage(), md.GetOutputType())
	if err != nil {
		return nil, err
	}
	meth := &Method{
		Service:               svc,
		MethodDescriptorProto: md,
		RequestType:           requestType,
		ResponseType:          responseType,
	}

	bind := &Binding{
		Method:  meth,
		ChQuery: opt.ChQuery,
		Query:   opt.Query,
	}
	meth.Binding = bind
	return meth, nil
}

func extractDatabaseOptions(meth *descriptorpb.MethodDescriptorProto) (*database.Query, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, database.E_Query) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, database.E_Query)
	opts, ok := ext.(*database.Query)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}
	return opts, nil
}
