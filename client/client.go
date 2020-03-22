package main

import (
	"context"
	"flag"
	"log"
	"time"

	ag "github.com/emxqm-11/grpc-ob/aggregator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:7777", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := ag.NewBankingAggregatorServiceClient(conn)

	runRecordRoute(client)

}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client ag.BankingAggregatorServiceClient) {
	// create 3 accounts
	accounts := []*ag.BankingAccount{
		{
			DisplayName:     "ANZ Adventures",
			IsOwned:         true,
			OpenStatus:      ag.BankingAccount_BANKING_ACCOUNT_OPEN_STATUS_OPEN,
			ProductCategory: ag.BankingProductCategory_MARGIN_LOANS,
			AccountId:       "1234",
		},
		{
			DisplayName:     "Commonweatlh Frequent Saver",
			IsOwned:         false,
			OpenStatus:      ag.BankingAccount_BANKING_ACCOUNT_OPEN_STATUS_OPEN,
			ProductCategory: ag.BankingProductCategory_CRED_AND_CHRG_CARDS,
			AccountId:       "12973910728378",
		},
		{
			DisplayName:     "Commonweatlh Low Rate Margin Loan",
			IsOwned:         true,
			OpenStatus:      ag.BankingAccount_BANKING_ACCOUNT_OPEN_STATUS_CLOSED,
			ProductCategory: ag.BankingProductCategory_MARGIN_LOANS,
			AccountId:       "11111111",
		},
		{
			DisplayName:     "NAB Overdraft Account",
			IsOwned:         true,
			OpenStatus:      ag.BankingAccount_BANKING_ACCOUNT_OPEN_STATUS_CLOSED,
			ProductCategory: ag.BankingProductCategory_OVERDRAFTS,
			AccountId:       "0000000sdq",
		},
		{
			DisplayName:     "eToro Trading Account",
			IsOwned:         true,
			OpenStatus:      ag.BankingAccount_BANKING_ACCOUNT_OPEN_STATUS_CLOSED,
			ProductCategory: ag.BankingProductCategory_TRADE_FINANCE,
			AccountId:       "zzz",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListAllAccounts(ctx)
	if err != nil {
		log.Fatalf("%v.ListAllAccounts(_) = _, %v", client, err)
	}
	for _, account := range accounts {
		if err := stream.Send(account); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, account, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Accounts List Response: %v", reply)
}
