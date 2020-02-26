package customob

import (
	"log"  "golang.org/x/net/context"
)
// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) ListAllAccounts(ctx context.Context, in *ListAllAccountsRequest) (*ResponseBankingAllAccountList, error) {
log.Printf("Receive message %s", in.Greeting)
return &ResponseBankingAllAccountList{Greeting: "bar"}, nil
}