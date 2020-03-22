package customob

import (
	"io"
	"log"

	ag "github.com/emxqm-11/grpc-ob/aggregator"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) ListAllAccounts(stream ag.BankingAggregatorService_ListAllAccountsServer) error {

	log.Println("***Received ListAllAccounts Request***")
	countOpenAccounts := int32(0)
	countClosedAccounts := int32(0)
	productCategoryMap := make(map[string]int32)
	accountMap := make(map[string]string)

	for {
		account, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&ag.ResponseBankingAllAccountList{
				ProductCategoryMap:  productCategoryMap,
				TotalOpenAccounts:   countOpenAccounts,
				TotalClosedAccounts: countClosedAccounts,
				Accounts:            accountMap,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		switch accountStatus := account.GetOpenStatus().String(); accountStatus {
		case "LIST_ACCOUNTS_REQUEST_OPEN_STATUS_OPEN":
			countOpenAccounts++
		case "LIST_ACCOUNTS_REQUEST_OPEN_STATUS_CLOSED":
			countClosedAccounts++
		}
		accountMap[account.AccountId] = account.GetDisplayName()
		productCategoryMap[account.ProductCategory.String()]++
	}
}
