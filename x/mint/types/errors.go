package types

import "cosmossdk.io/errors"

var (
	ErrInvalidSigner                                   = errors.Register(ModuleName, 1, "expected authority account as only signer for proposal message")
	ErrNegativeCirculatingSupply                       = errors.Register(ModuleName, 2, "negative circulating supply")
	ErrNegativeTargetEmissionPerToken                  = errors.Register(ModuleName, 3, "negative target emission per token")
	ErrInvalidPreviousRewardEmissionPerUnitStakedToken = errors.Register(ModuleName, 4, "invalid previous reward")
	ErrInvalidEcosystemTokensMinted                    = errors.Register(ModuleName, 5, "invalid ecosystem tokens minted")
	ErrMaxSupplyReached                                = errors.Register(ModuleName, 6, "max supply reached")
	ErrNilReputersStaked                               = errors.Register(ModuleName, 7, "nil reputers staked")
	ErrZeroDenominator                                 = errors.Register(ModuleName, 8, "zero denominator")
)
