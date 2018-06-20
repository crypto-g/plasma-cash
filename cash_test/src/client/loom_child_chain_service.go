package client

import (
	"fmt"
	"log"

	loom "github.com/loomnetwork/go-loom"
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
	"github.com/loomnetwork/go-loom/types"
)

// LoomChildChainService child client to reference server
type LoomChildChainService struct {
	url          string
	ChainID      string
	WriteURI     string
	ReadURI      string
	contractAddr string
	loomcontract *LoomContract
}

func (c *LoomChildChainService) CurrentBlock() (Block, error) {
	return c.Block(0) //asking for block zero gives latest
}

func (c *LoomChildChainService) BlockNumber() int64 {
	request := &pctypes.GetCurrentBlockRequest{}
	result := &pctypes.GetCurrentBlockResponse{}

	if err := c.loomcontract.StaticCallContract("GetCurrentBlockRequest", request, &result); err != nil {
		log.Fatalf("failed getting Block number - %v\n", err)

		return 0
	}

	log.Printf("get block height %v '\n", result.BlockHeight.Value.String())
	return result.BlockHeight.Value.Int64()
}

func (c *LoomChildChainService) Block(blknum int64) (Block, error) {
	fmt.Printf("trying to get Block data\n")
	blk := &types.BigUInt{*loom.NewBigUIntFromInt(blknum)}

	var result pctypes.GetBlockResponse
	params := &pctypes.GetBlockRequest{
		BlockHeight: blk,
	}

	if err := c.loomcontract.StaticCallContract("GetBlockRequest", params, &result); err != nil {
		log.Fatalf("failed getting Block data - %v\n", err)

		return err, nil
	}

	log.Printf("get block value %v '\n", result)

	return NewClientBlock(result.Block), nil
}

func (c *LoomChildChainService) Proof(blknum int64, uid int64) (*Proof, error) {

	return nil, nil
}

func (c *LoomChildChainService) SubmitBlock() error {
	fmt.Printf("submitting block\n")
	request := &pctypes.SubmitBlockToMainnetRequest{}
	//	params := &pctypes.GetBlockRequest{}

	if err := c.loomcontract.StaticCallContract("SubmitBlockToMainnet", request, nil); err != nil {
		log.Fatalf("failed submitting block - %v\n", err)

		return err
	}

	log.Println("succeeded submitting a block ")

	return nil
}

type LoomTx struct{}

func (c *LoomChildChainService) SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) (Tx, error) {
	fmt.Printf("trying to get send transaction\n")

	address := loom.MustParseAddress(newOwner)
	tx := &pctypes.PlasmaTx{
		Slot:          uint64(slot),
		PreviousBlock: &types.BigUInt{*loom.NewBigUIntFromInt(prevBlock)},
		Denomination:  &types.BigUInt{*loom.NewBigUIntFromInt(denomination)},
		NewOwner:      address.MarshalPB(),
	}

	params := &pctypes.PlasmaTxRequest{
		Plasmatx: tx,
	}

	if err := c.loomcontract.StaticCallContract("PlasmaTxRequest", params, nil); err != nil {
		log.Fatalf("failed trying to send transaction - %v\n", err)

		return err, nil
	}

	log.Printf("Transaction succeeded")

	return &LoomTx{}, nil
}

func NewLoomChildChainService(readuri, writeuri string) ChainServiceClient {
	return &LoomChildChainService{loomcontract: NewLoomContract(readuri, writeuri, "plasmacash")}
}
