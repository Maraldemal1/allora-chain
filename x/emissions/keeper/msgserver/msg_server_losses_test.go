package msgserver_test

import (
	"encoding/hex"

	cosmosMath "cosmossdk.io/math"
	alloraMath "github.com/allora-network/allora-chain/math"
	"github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestMsgInsertBulkReputerPayload() {
	ctx, msgServer := s.ctx, s.msgServer
	require := s.Require()

	reputerPrivateKey := secp256k1.GenPrivKey()
	reputerPublicKeyBytes := reputerPrivateKey.PubKey().Bytes()
	reputerAddr := sdk.AccAddress(reputerPrivateKey.PubKey().Address())

	workerPrivateKey := secp256k1.GenPrivKey()
	workerAddr := sdk.AccAddress(workerPrivateKey.PubKey().Address())

	s.emissionsKeeper.AddToTopicCreationWhitelist(ctx, reputerAddr)
	s.emissionsKeeper.AddToReputerWhitelist(ctx, reputerAddr)

	registrationInitialStake := cosmosMath.NewUint(100)

	// Create topic 0 and register reputer in it
	s.commonStakingSetup(ctx, reputerAddr, workerAddr, registrationInitialStake)

	// add in inference and forecast data
	keeper := s.emissionsKeeper
	topicId := uint64(0)
	block := types.BlockHeight(1)
	expectedInferences := types.Inferences{
		Inferences: []*types.Inference{
			{
				Value:   alloraMath.NewDecFromInt64(1), // Assuming NewDecFromInt64 exists and is appropriate
				Inferer: workerAddr.String(),
			},
		},
	}

	nonce := types.Nonce{BlockHeight: block} // Assuming block type cast to int64 if needed
	err := keeper.InsertInferences(ctx, topicId, nonce, expectedInferences)
	require.NoError(err, "InsertInferences should not return an error")

	expectedForecasts := types.Forecasts{
		Forecasts: []*types.Forecast{
			{
				TopicId:    topicId,
				Forecaster: workerAddr.String(),
			},
		},
	}

	nonce = types.Nonce{BlockHeight: int64(block)}
	err = keeper.InsertForecasts(ctx, topicId, nonce, expectedForecasts)
	s.Require().NoError(err)

	reputerValueBundle := &types.ValueBundle{
		TopicId:       0,
		Reputer:       reputerAddr.String(),
		CombinedValue: alloraMath.NewDecFromInt64(100),
		InfererValues: []*types.WorkerAttributedValue{
			{
				Worker: workerAddr.String(),
				Value:  alloraMath.NewDecFromInt64(100),
			},
		},
		ForecasterValues: []*types.WorkerAttributedValue{
			{
				Worker: workerAddr.String(),
				Value:  alloraMath.NewDecFromInt64(100),
			},
		},
		NaiveValue: alloraMath.NewDecFromInt64(100),
		OneOutInfererValues: []*types.WithheldWorkerAttributedValue{
			{
				Worker: workerAddr.String(),
				Value:  alloraMath.NewDecFromInt64(100),
			},
		},
		OneOutForecasterValues: []*types.WithheldWorkerAttributedValue{
			{
				Worker: workerAddr.String(),
				Value:  alloraMath.NewDecFromInt64(100),
			},
		},
		OneInForecasterValues: []*types.WorkerAttributedValue{
			{
				Worker: workerAddr.String(),
				Value:  alloraMath.NewDecFromInt64(100),
			},
		},
		ReputerRequestNonce: &types.ReputerRequestNonce{
			ReputerNonce: &types.Nonce{
				BlockHeight: 2,
			},
			WorkerNonce: &types.Nonce{
				BlockHeight: 1,
			},
		},
	}

	src := make([]byte, 0)
	src, err = reputerValueBundle.XXX_Marshal(src, true)
	require.NoError(err, "Marshall reputer value bundle should not return an error")

	valueBundleSignature, err := reputerPrivateKey.Sign(src)
	require.NoError(err, "Sign should not return an error")

	// Create a MsgInsertBulkReputerPayload message
	lossesMsg := &types.MsgInsertBulkReputerPayload{
		Sender:  reputerAddr.String(),
		TopicId: 0,
		ReputerRequestNonce: &types.ReputerRequestNonce{
			ReputerNonce: &types.Nonce{
				BlockHeight: 2,
			},
			WorkerNonce: &types.Nonce{
				BlockHeight: 1,
			},
		},
		ReputerValueBundles: []*types.ReputerValueBundle{
			{
				ValueBundle: reputerValueBundle,
				Signature:   valueBundleSignature,
				Pubkey:      hex.EncodeToString(reputerPublicKeyBytes),
			},
		},
	}

	_, err = msgServer.InsertBulkReputerPayload(ctx, lossesMsg)
	require.NoError(err, "InsertBulkReputerPayload should not return an error")
}

func (s *KeeperTestSuite) TestMsgInsertBulkReputerPayloadInvalidUnauthorized() {
	ctx, msgServer := s.ctx, s.msgServer
	require := s.Require()

	// Mock setup for addresses
	reputerAddr := nonAdminAccounts[0].String()
	workerAddr := sdk.AccAddress(PKS[1].Address()).String()

	// Create a MsgInsertBulkReputerPayload message
	lossesMsg := &types.MsgInsertBulkReputerPayload{
		Sender: reputerAddr,
		ReputerRequestNonce: &types.ReputerRequestNonce{
			ReputerNonce: &types.Nonce{
				BlockHeight: 10,
			},
			WorkerNonce: &types.Nonce{
				BlockHeight: 11,
			},
		},
		ReputerValueBundles: []*types.ReputerValueBundle{
			{
				ValueBundle: &types.ValueBundle{
					TopicId:       1,
					CombinedValue: alloraMath.NewDecFromInt64(100),
					InfererValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					ForecasterValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					NaiveValue: alloraMath.NewDecFromInt64(100),
					OneOutInfererValues: []*types.WithheldWorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					OneOutForecasterValues: []*types.WithheldWorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					OneInForecasterValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
				},
				Signature: []byte("Nonce + ReputerValueBundles Signature"),
			},
		},
	}

	_, err := msgServer.InsertBulkReputerPayload(ctx, lossesMsg)
	require.ErrorIs(err, types.ErrNotInReputerWhitelist, "InsertBulkReputerPayload should return an error")
}
