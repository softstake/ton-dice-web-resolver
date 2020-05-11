package resolver

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	api "github.com/tonradar/ton-api/proto"
	"github.com/tonradar/ton-dice-web-resolver/config"
	store "github.com/tonradar/ton-dice-web-server/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	ResolveQueryFileName = "resolve-query.fif"
)

const (
	// bet lifecycle states
	UNSAVED = iota - 1
	SAVED
	SENT
	RESOLVED
)

type Resolver struct {
	conf          *config.TonWebResolverConfig
	apiClient     api.TonApiClient
	storageClient store.BetsClient
}

func NewResolver(conf *config.TonWebResolverConfig) *Resolver {
	log.Println("Fetcher init...")
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		withClientUnaryInterceptor(),
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.StorageHost, conf.StoragePort), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	storageClient := store.NewBetsClient(conn)

	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", conf.TonAPIHost, conf.TonAPIPort), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	apiClient := api.NewTonApiClient(conn)

	return &Resolver{
		conf:          conf,
		apiClient:     apiClient,
		storageClient: storageClient,
	}
}

func (r *Resolver) ResolveQuery(betId int, seed string) error {
	log.Printf("Resolving bet with id %d...", betId)
	fileNameWithPath := ResolveQueryFileName
	fileNameStart := strings.LastIndex(fileNameWithPath, "/")
	fileName := fileNameWithPath[fileNameStart+1:]

	bocFile := strings.Replace(fileName, ".fif", ".boc", 1)

	_ = os.Remove(bocFile)

	var out bytes.Buffer
	cmd := exec.Command("fift", "-s", fileNameWithPath, r.conf.KeyFileBase, r.conf.ContractAddr, strconv.Itoa(betId), seed)

	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("cmd.Run() failed with: %v\n", err)
		return err
	}

	if FileExists(bocFile) {
		data, err := ioutil.ReadFile(bocFile)
		if err != nil {
			log.Println(err)
		}

		sendMessageRequest := &api.SendMessageRequest{
			Body: data,
		}

		sendMessageResponse, err := r.apiClient.SendMessage(context.Background(), sendMessageRequest)
		if err != nil {
			log.Printf("failed send message: %v\n", err)
			return err
		}

		log.Printf("send message status: %v\n", sendMessageResponse.Ok)

		return nil
	}

	return fmt.Errorf("file not found, maybe fift compile failed")
}

func (r *Resolver) Start() {
	log.Println("Resolver start")
	for {
		ctx := context.Background()

		getActiveBetsReq := &api.GetActiveBetsRequest{}
		getActiveBetsResp, err := r.apiClient.GetActiveBets(ctx, getActiveBetsReq)
		if err != nil {
			log.Printf("failed to get active bets: %v\n", err)
			continue
		}
		bets := getActiveBetsResp.GetBets()

		log.Printf("%d active bets received from smart-contract", len(bets))

		for _, bet := range bets {
			var state int32
			storedBet, err := r.getBet(ctx, bet.Id)
			if err != nil {
				if strings.Contains(err.Error(), "sql: no rows in result set") {
					state = UNSAVED
				} else {
					log.Println(err)
					continue
				}
			}

			var updatedAt time.Time
			if storedBet != nil {
				state = storedBet.State
				updatedAt, err = ptypes.Timestamp(storedBet.UpdatedAt)
			}

			if state == UNSAVED {
				req, err := BuildSaveBetRequest(bet)
				if err != nil {
					log.Printf("failed to build create bet request: %v\n", err)
					continue
				}
				_, err = r.storageClient.SaveBet(ctx, req)
				if err != nil {
					log.Printf("save bet in DB failed: %v\n", err)
					continue
				}
			} else {
				log.Println("the bet is already in storage")
			}

			if state == SENT && updatedAt.Add(15*time.Second).Before(time.Now()) {
				log.Println("last resolve attempt was less than 15s ago")
				continue
			} else {
				log.Println("try to resolve...")
			}

			err = r.ResolveQuery(int(bet.Id), bet.Seed)
			if err != nil {
				log.Printf("failed to resolve bet: %v\n", err)
				continue
			}

			req := &store.UpdateBetRequest{
				Id:    bet.Id,
				State: SENT,
			}
			_, err = r.storageClient.UpdateBet(ctx, req)
		}
	}
}
