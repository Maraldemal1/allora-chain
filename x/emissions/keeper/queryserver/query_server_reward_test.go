package queryserver_test

import (
	alloraMath "github.com/allora-network/allora-chain/math"
	"github.com/allora-network/allora-chain/x/emissions/types"
)

func (s *KeeperTestSuite) TestGetPreviousReputerRewardFraction() {
	ctx := s.ctx
	keeper := s.emissionsKeeper
	topicId := uint64(1)
	reputer := "reputerAddressExample"

	// Attempt to fetch a reward fraction before setting it
	req := &types.QueryPreviousReputerRewardFractionRequest{
		TopicId: topicId,
		Reputer: reputer,
	}
	response, err := s.queryServer.GetPreviousReputerRewardFraction(ctx, req)
	s.Require().NoError(err)
	defaultReward := response.RewardFraction
	notFound := response.NotFound

	s.Require().NoError(err, "Fetching reward fraction should not fail when not set")
	s.Require().True(defaultReward.IsZero(), "Should return zero reward fraction when not set")
	s.Require().True(notFound)

	// Now set a specific reward fraction
	setReward := alloraMath.NewDecFromInt64(50) // Assuming 0.50 as a fraction example
	_ = keeper.SetPreviousReputerRewardFraction(ctx, topicId, reputer, setReward)

	// Fetch and verify the reward fraction after setting
	response, err = s.queryServer.GetPreviousReputerRewardFraction(ctx, req)
	s.Require().NoError(err)
	fetchedReward := response.RewardFraction
	notFound = response.NotFound
	s.Require().NoError(err, "Fetching reward fraction should not fail after setting")
	s.Require().True(fetchedReward.Equal(setReward), "The fetched reward fraction should match the set value")
	s.Require().False(notFound, "Should not return no prior value after setting")
}

func (s *KeeperTestSuite) TestGetPreviousInferenceRewardFraction() {
	ctx := s.ctx
	keeper := s.emissionsKeeper
	topicId := uint64(1)
	worker := "workerAddressExample"

	// Attempt to fetch a reward fraction before setting it
	req := &types.QueryPreviousInferenceRewardFractionRequest{
		TopicId: topicId,
		Worker:  worker,
	}
	response, err := s.queryServer.GetPreviousInferenceRewardFraction(ctx, req)
	s.Require().NoError(err)
	defaultReward := response.RewardFraction
	noPrior := response.NotFound

	s.Require().NoError(err, "Fetching reward fraction should not fail when not set")
	s.Require().True(defaultReward.IsZero(), "Should return zero reward fraction when not set")
	s.Require().True(noPrior, "Should return no prior value when not set")

	// Now set a specific reward fraction
	setReward := alloraMath.NewDecFromInt64(75)
	_ = keeper.SetPreviousInferenceRewardFraction(ctx, topicId, worker, setReward)

	// Fetch and verify the reward fraction after setting
	fetchedReward, noPrior, err := keeper.GetPreviousInferenceRewardFraction(ctx, topicId, worker)
	response, err = s.queryServer.GetPreviousInferenceRewardFraction(ctx, req)
	s.Require().NoError(err)
	defaultReward = response.RewardFraction
	noPrior = response.NotFound
	s.Require().NoError(err, "Fetching reward fraction should not fail after setting")
	s.Require().True(fetchedReward.Equal(setReward), "The fetched reward fraction should match the set value")
	s.Require().False(noPrior, "Should not return no prior value after setting")
}

func (s *KeeperTestSuite) TestGetPreviousForecastRewardFraction() {
	ctx := s.ctx
	keeper := s.emissionsKeeper
	topicId := uint64(1)
	worker := "forecastWorkerAddress"

	// Attempt to fetch the reward fraction before setting it, expecting default value
	req := &types.QueryPreviousForecastRewardFractionRequest{
		TopicId: topicId,
		Worker:  worker,
	}
	response, err := s.queryServer.GetPreviousForecastRewardFraction(ctx, req)
	s.Require().NoError(err)
	defaultReward := response.RewardFraction
	noPrior := response.NotFound

	s.Require().NoError(err, "Fetching forecast reward fraction should not fail when not set")
	s.Require().True(defaultReward.IsZero(), "Should return zero forecast reward fraction when not set")
	s.Require().True(noPrior, "Should return no prior value when not set")

	// Now set a specific reward fraction
	setReward := alloraMath.NewDecFromInt64(75) // Assume setting it to 0.75
	_ = keeper.SetPreviousForecastRewardFraction(ctx, topicId, worker, setReward)

	// Fetch and verify the reward fraction after setting
	response, err = s.queryServer.GetPreviousForecastRewardFraction(ctx, req)
	s.Require().NoError(err)
	fetchedReward := response.RewardFraction
	noPrior = response.NotFound
	s.Require().NoError(err, "Fetching forecast reward fraction should not fail after setting")
	s.Require().True(fetchedReward.Equal(setReward), "The fetched forecast reward fraction should match the set value")
	s.Require().False(noPrior, "Should not return no prior value after setting")
}
