/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gpioblink.com/app/makemyfat/mkmyfat"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <imagePath> <file list (ex: a.txt b.txt c.txt ...)>",
	Short: "Create a new FAT32 image",
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		fileList := args[1:]
		fmt.Printf("imagePath %s, fileList %s\n", imagePath, fileList)

		err := mkmyfat.Create(imagePath, fileList)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
