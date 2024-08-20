/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"gpioblink.com/app/makemyfat/mkmyfat"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <imagePath> <fileSize> [<file list (ex: a.txt b.txt c.txt ...)>]",
	Short: "Create a new FAT32 image",
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		fileSizeText := args[1]

		fileSize, err := humanize.ParseBytes(fileSizeText)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("imagePath %s, fileSize %d \n", imagePath, fileSize)

		err = mkmyfat.Create(imagePath, int(fileSize))
		if err != nil {
			fmt.Println(err)
			return
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
