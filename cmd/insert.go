/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gpioblink.com/app/makemyfat/mkmyfat"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert <imagePath> <filePath> <entryNum>",
	Short: "Insert a new file to the FAT32 image",
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		filePath := args[1]
		entryNum, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("imagePath %s, filePath %s, entryNum %d\n", imagePath, filePath, entryNum)
		err = mkmyfat.Insert(imagePath, filePath, entryNum)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// insertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
