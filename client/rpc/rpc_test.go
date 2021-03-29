package rpc_test

import (
	"fmt"
	"testing"

	ctypes "github.com/line/ostracon/rpc/core/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/v2/client/rpc"
	"github.com/line/lbm-sdk/v2/codec/legacy"
	clitestutil "github.com/line/lbm-sdk/v2/testutil/cli"
	"github.com/line/lbm-sdk/v2/testutil/network"
	"github.com/line/lbm-sdk/v2/types/rest"
)

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), network.DefaultConfig())
	s.Require().NotNil(s.network)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestStatusCommand() {
	val0 := s.network.Validators[0]
	cmd := rpc.StatusCommand()

	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, []string{})
	s.Require().NoError(err)

	// Make sure the output has the validator moniker.
	s.Require().Contains(out.String(), fmt.Sprintf("\"moniker\":\"%s\"", val0.Moniker))
}

func (s *IntegrationTestSuite) TestLatestBlocks() {
	val0 := s.network.Validators[0]

	res, err := rest.GetRequest(fmt.Sprintf("%s/blocks/latest", val0.APIAddress))
	s.Require().NoError(err)

	var result ctypes.ResultBlock
	err = legacy.Cdc.UnmarshalJSON(res, &result)
	s.Require().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}