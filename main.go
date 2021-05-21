package main

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "rei",
	Long: "CLI script made in Go to serialize all documents in multiple collections to JSON and prints them out in a file",
	Short: "CLI script to serialize documents from multiple collections into a file.",
	Version: "1.0.0",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func init() {
	// Setup commands
	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Short: "Returns the version of Rei",
		RunE: func (cmd *cobra.Command, args []string) error {
			println(rootCmd.Version)
			return nil
		},
	})

	convertCmd := &cobra.Command{
		Use: "convert",
		Short: "Starts the conversion process",
		RunE: Execute,
	}

	convertCmd.Flags().StringP("uri", "u", "", "The connection URI string")
	convertCmd.Flags().StringP("db", "d", "", "The database")
	convertCmd.Flags().StringArrayP("collections", "c", []string{}, "List of collections to convert")

	rootCmd.AddCommand(convertCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	collections, _ := cmd.Flags().GetStringArray("collections")
	uri, _ := cmd.Flags().GetString("uri")
	db, _ := cmd.Flags().GetString("db")

	if len(collections) == 0 {
		return errors.New("missing collections array")
	}

	if uri == "" {
		return errors.New("missing connection uri")
	}

	if db == "" {
		return errors.New("missing database name")
	}

	println("Using connection string:", uri, "with database", db, "...")
	return nil
}
