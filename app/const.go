package app

const (
	Name = "cheqd-node"
	Home = ".cheqdnode"

	// Set ProtocolVersion to app's major version number
	ProtocolVersion = 1
	UpgradeName     = "v1"

	// allowed msg types of ica host
	authzMsgExec                        = "/cosmos.authz.v1beta1.MsgExec"
	authzMsgGrant                       = "/cosmos.authz.v1beta1.MsgGrant"
	authzMsgRevoke                      = "/cosmos.authz.v1beta1.MsgRevoke"
	bankMsgSend                         = "/cosmos.bank.v1beta1.MsgSend"
	bankMsgMultiSend                    = "/cosmos.bank.v1beta1.MsgMultiSend"
	distrMsgSetWithdrawAddr             = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	distrMsgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	distrMsgFundCommunityPool           = "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	distrMsgWithdrawDelegatorReward     = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	feegrantMsgGrantAllowance           = "/cosmos.feegrant.v1beta1.MsgGrantAllowance"
	feegrantMsgRevokeAllowance          = "/cosmos.feegrant.v1beta1.MsgRevokeAllowance"
	govMsgVoteWeighted                  = "/cosmos.gov.v1beta1.MsgVoteWeighted"
	govMsgSubmitProposal                = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	govMsgDeposit                       = "/cosmos.gov.v1beta1.MsgDeposit"
	govMsgVote                          = "/cosmos.gov.v1beta1.MsgVote"
	stakingMsgEditValidator             = "/cosmos.staking.v1beta1.MsgEditValidator"
	stakingMsgDelegate                  = "/cosmos.staking.v1beta1.MsgDelegate"
	stakingMsgUndelegate                = "/cosmos.staking.v1beta1.MsgUndelegate"
	stakingMsgBeginRedelegate           = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	stakingMsgCreateValidator           = "/cosmos.staking.v1beta1.MsgCreateValidator"
	vestingMsgCreateVestingAccount      = "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
	ibcMsgTransfer                      = "/ibc.applications.transfer.v1.MsgTransfer"

	// cheqd namespace
	didMsgCreateDidDoc        = "/cheqd.did.v2.MsgCreateDidDoc"
	didMsgUpdateDidDoc        = "/cheqd.did.v2.MsgUpdateDidDoc"
	didMsgDeactivateDidDoc    = "/cheqd.did.v2.MsgDeactivateDidDoc"
	resourceMsgCreateResource = "/cheqd.resource.v2.MsgCreateResource"
)
