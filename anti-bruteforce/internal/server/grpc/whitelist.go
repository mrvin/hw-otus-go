package grpcserver

import (
	"context"
	"net"
	"slices"

	"github.com/mrvin/hw-otus-go/anti-bruteforce/internal/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) AddNetworkToWhitelist(ctx context.Context, req *api.ReqNetwork) (*emptypb.Empty, error) {
	_, network, err := net.ParseCIDR(req.GetNetwork())
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if err := s.storage.AddNetworkToWhitelist(ctx, network); err != nil {
		return &emptypb.Empty{}, err
	}
	s.storage.CacheWhitelist = append(s.storage.CacheWhitelist, network)

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteNetworkFromWhitelist(ctx context.Context, req *api.ReqNetwork) (*emptypb.Empty, error) {
	_, network, err := net.ParseCIDR(req.GetNetwork())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	if err := s.storage.DeleteNetworkFromWhitelist(ctx, network); err != nil {
		return &emptypb.Empty{}, err
	}
	for i := 0; i < len(s.storage.CacheWhitelist); i++ {
		if network.String() == s.storage.CacheWhitelist[i].String() {
			s.storage.CacheWhitelist = slices.Delete(s.storage.CacheWhitelist, i, i+1)
			break
		}
	}

	return &emptypb.Empty{}, nil
}

func list(list []*net.IPNet) *api.ResListNetworks {
	pbNetworks := make([]string, len(list))
	for i, network := range list {
		pbNetworks[i] = network.String()
	}

	return &api.ResListNetworks{Networks: pbNetworks}
}

func (s *Server) Whitelist(_ context.Context, _ *emptypb.Empty) (*api.ResListNetworks, error) {
	return list(s.storage.CacheWhitelist), nil
}
