package main

import (
	"dev.floofy.nino/rei/pkg/connection"
	"dev.floofy.nino/rei/pkg/writer"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
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
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
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
	flag := convertCmd.Flags()
	flag.StringP("uri", "u", "", "The connection URI string")
	flag.StringP("db", "d", "", "The database")
	flag.StringP("collections", "c", "", "List of collections to convert")
	rootCmd.AddCommand(convertCmd)
}

func Execute(cmd *cobra.Command, args []string) error {
	collections, _ := cmd.Flags().GetString("collections")
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

	client, err := connection.CreateMongoClient(uri); if err != nil {
		return err
	}

	mongoDb := client.Database(db)
	cols := strings.Split(collections, ",")
	logrus.Info(collections, cols)

	for i, c := range cols {
		fmt.Printf("[%d/%d | %s] Now retrieving documents...\n", i + 1, len(cols), c)

		collection := mongoDb.Collection(c)
		if err := writer.WriteDocumentsToFile(c + ".json", collection); err != nil {
			return err
		}
	}

	return nil
}
