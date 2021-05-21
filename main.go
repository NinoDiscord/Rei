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

	logrus.Info(string(os.PathSeparator))
	cwd, _ := os.Getwd()
	flag := convertCmd.Flags()
	flag.StringP("uri", "u", "", "The connection URI string")
	flag.StringP("db", "n", "", "The database")
	flag.StringP("collections", "c", "", "List of collections to convert")
	flag.StringP("directory", "d", cwd + string(os.PathSeparator) + "data", "A directory to place it in, it'll default to `$ROOT/data/<db>/<collection>.json`")
	rootCmd.AddCommand(convertCmd)
}

func Execute(cmd *cobra.Command, _ []string) error {
	cwd, _ := os.Getwd()
	collections, _ := cmd.Flags().GetString("collections")
	uri, _ := cmd.Flags().GetString("uri")
	db, _ := cmd.Flags().GetString("db")
	dir, _ := cmd.Flags().GetString("directory")

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

	base := cwd + string(os.PathSeparator) + "data"
	if dir != "" {
		base = dir
	}
	_, err = os.Stat(base)
	if os.IsNotExist(err) {
		logrus.Infof("Create %s/", base)
		err := os.Mkdir(base, 0644)
		if err != nil {
			logrus.Fatalf("Failed to create data directory: %v", err)
		}
		_, err = os.Stat(base + string(os.PathSeparator) + db)
		if os.IsNotExist(err) {
			logrus.Infof("Create %s/%s", base, db)
			err := os.Mkdir(base + string(os.PathSeparator) + db, 0644)
			if err != nil {
				logrus.Fatalf("Failed to create the directory data/%s: %v", db, err)
			}
		}
	}

	mongoDb := client.Database(db)
	cols := strings.Split(collections, ",")
	logrus.Info(collections, cols)

	for i, c := range cols {
		fmt.Printf("[%d/%d | %s] Now retrieving documents...\n", i + 1, len(cols), c)

		collection := mongoDb.Collection(c)
		if err := writer.WriteDocumentsToFile(base + string(os.PathSeparator) + db + string(os.PathSeparator) + c + ".json", collection); err != nil {
			return err
		}
	}

	return nil
}
