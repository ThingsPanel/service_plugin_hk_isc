package main

import (
	"fmt"
	"sync"
)

// Service 接口定义了所有服务应该实现的方法
type Service interface {
	Run() error
	Close()
}

// Server 结构体现在包含了多个服务
type Server struct {
	services []Service
	wg       sync.WaitGroup
}

// NewServer 创建一个新的 Server 实例
func NewServer() *Server {
	return &Server{
		services: make([]Service, 0),
	}
}

// AddService 添加一个新的服务到 Server
func (s *Server) AddService(service Service) {
	s.services = append(s.services, service)
}

// Run 启动所有的服务
func (s *Server) Run() error {
	for _, service := range s.services {
		s.wg.Add(1)
		go func(srv Service) {
			defer s.wg.Done()
			if err := srv.Run(); err != nil {
				fmt.Printf("服务运行出错: %v\n", err)
			}
		}(service)
	}
	s.wg.Wait()
	return nil
}

// Close 关闭所有的服务
func (s *Server) Close() {
	for _, service := range s.services {
		service.Close()
	}
}
