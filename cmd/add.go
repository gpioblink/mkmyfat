/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <imagePath> <file list (ex: a.txt b.txt c.txt ...)>",
	Short: "Add files to new FAT32 image",
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		filePaths := args[1:]

		fmt.Printf("imagePath %s, filePaths %s \n", imagePath, filePaths)
		fmt.Printf("Not implemented. \n")

		// err := mkmyfat.Add(imagePath, filePaths)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
