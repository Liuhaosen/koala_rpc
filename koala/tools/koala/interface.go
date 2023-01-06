package main

import "github.com/emicklei/proto"

type Generator interface {
	Run(opt *Option, metaData *ServiceMetaData) (err error)
}

type ServiceMetaData struct {
	Service  *proto.Service
	Messages []*proto.Message
	Rpc      []*proto.RPC
	Package  *proto.Package
	Prefix   string
}
