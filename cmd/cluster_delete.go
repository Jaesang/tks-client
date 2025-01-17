/*
Copyright © 2022 SK Telecom <https://github.com/openinfradev>

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
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

// clusterDeleteCmd represents the delete command
var clusterDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a TKS Cluster.",
	Long: `Delete a TKS Cluster to AWS.
	
Example:
tks cluster delete <CLUSTER_ID>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify cluster name.")
			fmt.Println("Usage: tks cluster delete <CLUSTER_ID>")
			os.Exit(1)
		}

		var conn *grpc.ClientConn
		tksClusterLcmUrl = viper.GetString("tksClusterLcmUrl")
		if tksClusterLcmUrl == "" {
			fmt.Println("You must specify tksClusterLcmUrl at config file")
			os.Exit(1)
		}
		conn, err := grpc.Dial(tksClusterLcmUrl, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterLcmServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		data := pb.IDRequest{}
		data.Id = args[0]
		m := protojson.MarshalOptions{
			Indent:        "  ",
			UseProtoNames: true,
		}
		jsonBytes, _ := m.Marshal(&data)
		verbose, err := rootCmd.PersistentFlags().GetBool("verbose")
		if verbose {
			fmt.Println("Proto Json data...")
			fmt.Println(string(jsonBytes))
		}
		r, err := client.DeleteCluster(ctx, &data)
		fmt.Println(r)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("The request to delete cluster ", args[0], " was accepted.")
		}
	},
}

func init() {
	clusterCmd.AddCommand(clusterDeleteCmd)
}
