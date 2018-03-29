/*
Copyright Ziggurat Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package node

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/inklabsfoundation/inkchain/peer/common"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func statusCmd() *cobra.Command {
	return nodeStatusCmd
}

var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Returns status of the node.",
	Long:  `Returns the status of the running node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return status()
	},
}

func status() (err error) {

	adminClient, err := common.GetAdminClient()
	if err != nil {
		logger.Warningf("%s", err)
		fmt.Println(&pb.ServerStatus{Status: pb.ServerStatus_UNKNOWN})
		return err
	}

	status, err := adminClient.GetStatus(context.Background(), &empty.Empty{})
	if err != nil {
		logger.Infof("Error trying to get status from local peer: %s", err)
		err = fmt.Errorf("Error trying to connect to local peer: %s", err)
		fmt.Println(&pb.ServerStatus{Status: pb.ServerStatus_UNKNOWN})
		return err
	}
	fmt.Println(status)
	return nil
}
