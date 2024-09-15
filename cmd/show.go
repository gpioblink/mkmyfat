/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"gpioblink.com/app/makemyfat/mkmyfat"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show <imagePath>",
	Short: "Show BIOS Parameter Block",
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		// fmt.Print(mkmyfat.PrintBPBFromFile(imagePath))
		mkmyfat.ShowImageInfo(imagePath)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
