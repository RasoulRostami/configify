/*
Copyright © 2023 Rasoul Rostami rasoul.rostami.dev@gmail.com
*/
package cmd

import (
	"configify/publishers"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstarp all settings.",
	Long:  `Delete all existing configuration file and write again.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bootstrap called")
		red := publishers.PublisherFactory(
			viper.GetString("publisher"),
			viper.GetStringMap(fmt.Sprintf("%s_config", viper.GetString("publisher"))),
			publishers.DecoderFactory(viper.GetString("decoder")))

		result, error := red.Get("mykey")
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println(result["key1"])

	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bootstrapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bootstrapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
