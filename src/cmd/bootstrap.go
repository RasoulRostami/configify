/*
Copyright © 2023 Rasoul Rostami rasoul.rostami.dev@gmail.com
*/
package cmd

import (
	"configify/databases"
	"configify/helpers"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstarp all settings.",
	Long:  `Delete all existing configuration file and write again.`,

	Run: func(cmd *cobra.Command, args []string) {
		helpers.Logger.Info("Bootstrap is called.")
		messages := make(chan databases.Message, 100)
		db := databases.PublisherFactory(
			viper.GetString("publisher"),
			viper.GetStringMap(fmt.Sprintf("%s_config", viper.GetString("publisher"))),
			databases.DecoderFactory(viper.GetString("decoder")))

		var wg sync.WaitGroup
		go Performer(messages, &wg)
		db.Keys("*", messages, &wg)
		wg.Wait()
		helpers.Logger.Info("Bootstrap successfully finished.")
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
