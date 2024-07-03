package grpcserver

import (
	"context"
	"net"
	"slices"

	"github.com/mrvin/hw-otus-go/anti-bruteforce/internal/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) AddNetworkToBlacklist(ctx context.Context, req *api.ReqNetwork) (*emptypb.Empty, error) {
	_, network, err := net.ParseCIDR(req.GetNetwork())
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if err := s.storage.AddNetworkToBlacklist(ctx, network); err != nil {
		return &emptypb.Empty{}, err
	}
	s.storage.CacheBlacklist = append(s.storage.CacheBlacklist, network)

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteNetworkFromBlacklist(ctx context.Context, req *api.ReqNetwork) (*emptypb.Empty, error) {
	_, network, err := net.ParseCIDR(req.GetNetwork())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	if err := s.storage.DeleteNetworkFromBlacklist(ctx, network); err != nil {
		return &emptypb.Empty{}, err
	}
	for i := 0; i < len(s.storage.CacheBlacklist); i++ {
		if network.String() == s.storage.CacheBlacklist[i].String() {
			s.storage.CacheBlacklist = slices.Delete(s.storage.CacheBlacklist, i, i+1)
			break
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Blacklist(_ context.Context, _ *emptypb.Empty) (*api.ResListNetworks, error) {
	return list(s.storage.CacheBlacklist), nil
}
