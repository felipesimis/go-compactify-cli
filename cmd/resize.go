/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	directory string
	width     int
	height    int
)

func resizeRun(cmd *cobra.Command, args []string) {
	fmt.Println("resize called")
}

var resizeCmd = &cobra.Command{
	Use:   "resize",
	Args:  cobra.NoArgs,
	Short: "Resize an image to specified dimensions",
	Long: `Resize an image to a specific width and height.
This command allows you to change the dimensions of an image, which can be useful for optimizing images for 
different uses, such as web, mobile, or print. You can specify the desired width and height, 
and the image will be resized accordingly.`,
	Run: resizeRun,
}

func init() {
	rootCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing the image to resize")
	resizeCmd.Flags().IntVarP(&width, "width", "w", 0, "Desired width of the image")
	resizeCmd.Flags().IntVarP(&height, "height", "H", 0, "Desired height of the image")

	resizeCmd.MarkFlagRequired("directory")
	resizeCmd.MarkFlagRequired("width")
	resizeCmd.MarkFlagRequired("height")
}
