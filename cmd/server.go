package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/shield/config"
	"github.com/odpf/shield/internal/store/postgres/migrations"
	"github.com/odpf/shield/pkg/db"
	shieldlogger "github.com/odpf/shield/pkg/logger"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

func ServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server <command>",
		Aliases: []string{"s"},
		Short:   "Server management",
		Long:    "Server management commands.",
		Example: heredoc.Doc(`
			$ shield server init
			$ shield server start
			$ shield server start -c ./config.yaml
			$ shield server migrate
			$ shield server migrate -c ./config.yaml
			$ shield server migrate-rollback
			$ shield server migrate-rollback -c ./config.yaml
		`),
	}

	cmd.AddCommand(serverInitCommand())
	cmd.AddCommand(serverStartCommand())
	cmd.AddCommand(serverMigrateCommand())
	cmd.AddCommand(serverMigrateRollbackCommand())

	return cmd
}

func serverInitCommand() *cobra.Command {
	var configFile string
	var resourcesURL string
	var rulesURL string

	c := &cli.Command{
		Use:   "init",
		Short: "Initialize server",
		Long: heredoc.Doc(`
			Initializing server. Creating a sample of shield server config.
			Default: ./config.yaml
		`),
		Example: "shield server init",
		RunE: func(cmd *cli.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			defaultResourcesURL := fmt.Sprintf("file://%s", path.Join(cwd, "resources_config"))
			defaultRulesURL := fmt.Sprintf("file://%s", path.Join(cwd, "rules"))

			if resourcesURL == "" {
				resourcesURL = defaultResourcesURL
			}
			if rulesURL == "" {
				rulesURL = defaultRulesURL
			}

			if err := config.Init(resourcesURL, rulesURL, configFile); err != nil {
				return err
			}

			fmt.Printf("server config created: %v\n", configFile)
			return nil
		},
	}

	c.Flags().StringVarP(&configFile, "output", "o", "./config.yaml", "Output config file path")
	c.Flags().StringVarP(&resourcesURL, "resources", "r", "", heredoc.Doc(`
		URL path of resources. Full path prefixed with scheme where resources config yaml files are kept
		e.g.:
		local storage file "file:///tmp/resources_config"
		GCS Bucket "gs://shield-bucket-example"
		(default: file://{pwd}/resources_config)
	`))
	c.Flags().StringVarP(&rulesURL, "rule", "u", "", heredoc.Doc(`
		URL path of rules. Full path prefixed with scheme where ruleset yaml files are kept
		e.g.:
		local storage file "file:///tmp/rules"
		GCS Bucket "gs://shield-bucket-example"
		(default: file://{pwd}/rules)
	`))

	return c
}

func serverStartCommand() *cobra.Command {
	var configFile string

	c := &cli.Command{
		Use:     "start",
		Short:   "Start server and proxy default on port 8080",
		Example: "shield server start",
		RunE: func(cmd *cli.Command, args []string) error {
			appConfig, err := config.Load(configFile)
			if err != nil {
				panic(err)
			}
			logger := shieldlogger.InitLogger(appConfig.Log)

			return StartServer(logger, appConfig)
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	return c
}

func serverMigrateCommand() *cobra.Command {
	var configFile string

	c := &cli.Command{
		Use:     "migrate",
		Short:   "Run DB Schema Migrations",
		Example: "shield migrate",
		RunE: func(c *cli.Command, args []string) error {
			appConfig, err := config.Load(configFile)
			if err != nil {
				panic(err)
			}

			return db.RunMigrations(db.Config{
				Driver: appConfig.DB.Driver,
				URL:    appConfig.DB.URL,
			}, migrations.MigrationFs, migrations.ResourcePath)
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	return c
}

func serverMigrateRollbackCommand() *cobra.Command {
	var configFile string

	c := &cli.Command{
		Use:     "migration-rollback",
		Short:   "Run DB Schema Migrations Rollback to last state",
		Example: "shield migration-rollback",
		RunE: func(c *cli.Command, args []string) error {
			appConfig, err := config.Load(configFile)
			if err != nil {
				panic(err)
			}

			return db.RunRollback(db.Config{
				Driver: appConfig.DB.Driver,
				URL:    appConfig.DB.URL,
			}, migrations.MigrationFs, migrations.ResourcePath)
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	return c
}
