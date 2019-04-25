package main

import (
	"github.com/andreymgn/RSOI-category/pkg/category"
	"github.com/andreymgn/RSOI/pkg/tracer"
)

func runPost(port int, connString, jaegerAddr string) error {
	tracer, closer, err := tracer.NewTracer("category", jaegerAddr)
	if err != nil {
		return err
	}

	defer closer.Close()

	server, err := category.NewServer(connString)
	if err != nil {
		return err
	}

	return server.Start(port, tracer)
}
