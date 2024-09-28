/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gpioblink.com/app/makemyfat/mkmyfat"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize <imagePath>",
	Short: "Visualiz any files as jpeg",
	Run: func(cmd *cobra.Command, args []string) {
		volumePath := args[0]
		outPath := args[1]
		err := mkmyfat.SaveVisualizeBinary(volumePath, outPath)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// visualizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// visualizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
