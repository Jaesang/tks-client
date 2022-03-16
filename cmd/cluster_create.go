/*
Copyright © 2021 SK Telecom <https://github.com/openinfradev>

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
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	address = "tks-cluster-lcm.taco-cat.xyz:9110"
)

// clusterCreateCmd represents the create command
var clusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a TACO Cluster.",
	Long: `Create a TACO Cluster to AWS.
	
Example:
tks cluster create <CLUSTERNAME> --contract-id <CONTRACTID> --csp-id <CSPID>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify cluster name.")
			fmt.Println("Usage: tks cluster create <CLUSTERNAME> --contract-id <CONTRACTID> --csp-id <CSPID>")
			os.Exit(1)
		}
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterLcmServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Minute)
		defer cancel()
		data := make([]pb.CreateClusterRequest, 1)
		conf := &pb.ClusterConf{}
		ContractId, _ := cmd.Flags().GetString("contract-id")
		CspId, _ := cmd.Flags().GetString("csp-id")
		ClusterName := args[0]
		data[0].ContractId = ContractId
		data[0].CspId = CspId
		data[0].Name = ClusterName
		conf.MasterFlavor = "hello"
		conf.MasterReplicas = 10
		conf.MasterRootSize = 50
		conf.WorkerFlavor = "hello"
		conf.WorkerRootSize = 50
		conf.WorkerReplicas = 10
		conf.K8SVersion = "Hello"
		data[0].Conf = conf

		m := protojson.MarshalOptions{
			Indent:        "  ",
			UseProtoNames: true,
		}
		jsonBytes, _ := m.Marshal(&data[0])
		fmt.Println("Proto Json data...")
		fmt.Println(string(jsonBytes))
		r, err := client.CreateCluster(ctx, &data[0])
		fmt.Println(r)
	},
}

func init() {
	clusterCmd.AddCommand(clusterCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clusterCreateCmd.Flags().String("contract-id", "", "Contract ID")
	clusterCreateCmd.MarkFlagRequired("contract-id")
	clusterCreateCmd.Flags().String("csp-id", "", "CSP ID")
	clusterCreateCmd.MarkFlagRequired("csp-id")
}