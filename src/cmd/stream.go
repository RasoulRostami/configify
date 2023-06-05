/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"configify/databases"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Listen to db pubsub event and update configs.",
	Long:  `Listen to db pubsub event, get new messages and update config files.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("INFO stream is called.")
		messages := make(chan databases.Message, 100)
		db := databases.PublisherFactory(
			viper.GetString("publisher"),
			viper.GetStringMap(fmt.Sprintf("%s_config", viper.GetString("publisher"))),
			databases.DecoderFactory(viper.GetString("decoder")))

		var wg sync.WaitGroup
		go Performer(messages, &wg)
		db.Stream(messages, &wg)
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
