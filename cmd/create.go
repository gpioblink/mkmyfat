/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"gpioblink.com/app/makemyfat/mkmyfat"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <imagePath> <fileSize> [<fileExt> <numOfFiles> <eachFileSize>]",
	Short: "Create a new FAT32 image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 2 {
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
		} else {
			imagePath := args[0]
			fileSizeText := args[1]
			fileExt := args[2]
			numOfFilesText := args[3]
			eachFileSizeText := args[4]

			fileSize, err := humanize.ParseBytes(fileSizeText)
			if err != nil {
				fmt.Println(err)
				return
			}

			numOfFiles, err := strconv.Atoi(numOfFilesText)
			if err != nil {
				fmt.Println(err)
				return
			}

			eachFileSize, err := humanize.ParseBytes(eachFileSizeText)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("imagePath %s, fileSize %d, fileExt %s, numOfFiles %d, eachFileSize %d \n", imagePath, fileSize, fileExt, numOfFiles, eachFileSize)

			err = mkmyfat.CreateWithEmptyFiles(imagePath, int(fileSize), fileExt, int(numOfFiles), int(eachFileSize))
			if err != nil {
				fmt.Println(err)
				return
			}
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
