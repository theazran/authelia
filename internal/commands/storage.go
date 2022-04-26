package commands

import (
	"github.com/spf13/cobra"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
)

// NewStorageCmd returns a new storage *cobra.Command.
func NewStorageCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:               "storage",
		Short:             "Manage the Authelia storage",
		Args:              cobra.NoArgs,
		PersistentPreRunE: storagePersistentPreRunE,
	}

	cmdWithConfigFlags(cmd, true, []string{"configuration.yml"})

	cmd.PersistentFlags().String("encryption-key", "", "the storage encryption key to use")

	cmd.PersistentFlags().String("sqlite.path", "", "the SQLite database path")

	cmd.PersistentFlags().String("mysql.host", "", "the MySQL hostname")
	cmd.PersistentFlags().Int("mysql.port", 3306, "the MySQL port")
	cmd.PersistentFlags().String("mysql.database", "authelia", "the MySQL database name")
	cmd.PersistentFlags().String("mysql.username", "authelia", "the MySQL username")
	cmd.PersistentFlags().String("mysql.password", "", "the MySQL password")

	cmd.PersistentFlags().String("postgres.host", "", "the PostgreSQL hostname")
	cmd.PersistentFlags().Int("postgres.port", 5432, "the PostgreSQL port")
	cmd.PersistentFlags().String("postgres.database", "authelia", "the PostgreSQL database name")
	cmd.PersistentFlags().String("postgres.schema", "public", "the PostgreSQL schema name")
	cmd.PersistentFlags().String("postgres.username", "authelia", "the PostgreSQL username")
	cmd.PersistentFlags().String("postgres.password", "", "the PostgreSQL password")
	cmd.PersistentFlags().String("postgres.ssl.mode", "disable", "the PostgreSQL ssl mode")
	cmd.PersistentFlags().String("postgres.ssl.root_certificate", "", "the PostgreSQL ssl root certificate file location")
	cmd.PersistentFlags().String("postgres.ssl.certificate", "", "the PostgreSQL ssl certificate file location")
	cmd.PersistentFlags().String("postgres.ssl.key", "", "the PostgreSQL ssl key file location")

	cmd.AddCommand(
		newStorageMigrateCmd(),
		newStorageSchemaInfoCmd(),
		newStorageEncryptionCmd(),
		newStorageUserCmd(),
	)

	return cmd
}

func newStorageUserCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "user",
		Short: "Manages user settings",
	}

	cmd.AddCommand(
		newStorageUserIdentifiersCmd(),
		newStorageTOTPCmd(),
	)

	return cmd
}

func newStorageUserIdentifiersCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "identifiers",
		Short: "Manages user opaque identifiers",
	}

	cmd.AddCommand(
		newStorageUserIdentifiersExportCmd(),
		newStorageUserIdentifiersImportCmd(),
		newStorageUserIdentifiersGenerateCmd(),
		newStorageUserIdentifiersAddCmd(),
	)

	return cmd
}

func newStorageUserIdentifiersExportCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "export",
		Short: "Export the identifiers to a YAML file",
		RunE:  storageUserIdentifiersExport,
	}

	cmd.Flags().StringP("file", "f", "user-opaque-identifiers.yml", "The file name for the YAML export")

	return cmd
}

func newStorageUserIdentifiersImportCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "import",
		Short: "Import the identifiers from a YAML file",
		RunE:  storageUserIdentifiersImport,
	}

	cmd.Flags().StringP("file", "f", "user-opaque-identifiers.yml", "The file name for the YAML import")

	return cmd
}

func newStorageUserIdentifiersGenerateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate opaque identifiers in bulk",
		RunE:  storageUserIdentifiersGenerate,
	}

	cmd.Flags().StringSlice("users", nil, "The list of users to generate the opaque identifiers for")
	cmd.Flags().StringSlice("services", []string{identifierServiceOpenIDConnect}, "The list of services to generate the opaque identifiers for, valid values are: openid")
	cmd.Flags().StringSlice("sectors", []string{""}, "The list of sectors to generate identifiers for")

	return cmd
}

func newStorageUserIdentifiersAddCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "add [username]",
		Short: "Add an identifiers to the database",
		Args:  cobra.ExactArgs(1),
		RunE:  storageUserIdentifiersAdd,
	}

	cmd.Flags().String("identifier", "", "The optional version 4 UUID to use, if not set a random one will be used")
	cmd.Flags().String("service", identifierServiceOpenIDConnect, "The service to add the identifier for, valid values are: openid")
	cmd.Flags().String("sector", "", "The sector identifier to use (should usually be blank)")

	return cmd
}

func newStorageEncryptionCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "encryption",
		Short: "Manages encryption",
	}

	cmd.AddCommand(
		newStorageEncryptionChangeKeyCmd(),
		newStorageEncryptionCheckCmd(),
	)

	return cmd
}

func newStorageEncryptionCheckCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "check",
		Short: "Checks the encryption key against the database data",
		RunE:  storageSchemaEncryptionCheckRunE,
	}

	cmd.Flags().Bool("verbose", false, "enables verbose checking of every row of encrypted data")

	return cmd
}

func newStorageEncryptionChangeKeyCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "change-key",
		Short: "Changes the encryption key",
		RunE:  storageSchemaEncryptionChangeKeyRunE,
	}

	cmd.Flags().String("new-encryption-key", "", "the new key to encrypt the data with")

	return cmd
}

func newStorageTOTPCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "totp",
		Short: "Manage TOTP configurations",
	}

	cmd.AddCommand(
		newStorageTOTPGenerateCmd(),
		newStorageTOTPDeleteCmd(),
		newStorageTOTPExportCmd(),
	)

	return cmd
}

func newStorageTOTPGenerateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "generate [username]",
		Short: "Generate a TOTP configuration for a user",
		RunE:  storageTOTPGenerateRunE,
		Args:  cobra.ExactArgs(1),
	}

	cmd.Flags().String("secret", "", "Optionally set the TOTP shared secret as base32 encoded bytes (no padding), it's recommended to not set this option unless you're restoring an TOTP config")
	cmd.Flags().Uint("secret-size", schema.TOTPSecretSizeDefault, "set the TOTP secret size")
	cmd.Flags().Uint("period", 30, "set the TOTP period")
	cmd.Flags().Uint("digits", 6, "set the TOTP digits")
	cmd.Flags().String("algorithm", "SHA1", "set the TOTP algorithm")
	cmd.Flags().String("issuer", "Authelia", "set the TOTP issuer")
	cmd.Flags().BoolP("force", "f", false, "forces the TOTP configuration to be generated regardless if it exists or not")
	cmd.Flags().StringP("path", "p", "", "path to a file to create a PNG file with the QR code (optional)")

	return cmd
}

func newStorageTOTPDeleteCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "delete username",
		Short: "Delete a TOTP configuration for a user",
		RunE:  storageTOTPDeleteRunE,
		Args:  cobra.ExactArgs(1),
	}

	return cmd
}

func newStorageTOTPExportCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "export",
		Short: "Performs exports of the TOTP configurations",
		RunE:  storageTOTPExportRunE,
	}

	cmd.Flags().String("format", storageExportFormatURI, "sets the output format")
	cmd.Flags().String("dir", "", "used with the png output format to specify which new directory to save the files in")

	return cmd
}

func newStorageSchemaInfoCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "schema-info",
		Short: "Show the storage information",
		RunE:  storageSchemaInfoRunE,
	}

	return cmd
}

// NewMigrationCmd returns a new Migration Cmd.
func newStorageMigrateCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Perform or list migrations",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newStorageMigrateUpCmd(), newStorageMigrateDownCmd(),
		newStorageMigrateListUpCmd(), newStorageMigrateListDownCmd(),
		newStorageMigrateHistoryCmd(),
	)

	return cmd
}

func newStorageMigrateHistoryCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "history",
		Short: "Show migration history",
		Args:  cobra.NoArgs,
		RunE:  storageMigrateHistoryRunE,
	}

	return cmd
}

func newStorageMigrateListUpCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list-up",
		Short: "List the up migrations available",
		Args:  cobra.NoArgs,
		RunE:  newStorageMigrateListRunE(true),
	}

	return cmd
}

func newStorageMigrateListDownCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list-down",
		Short: "List the down migrations available",
		Args:  cobra.NoArgs,
		RunE:  newStorageMigrateListRunE(false),
	}

	return cmd
}

func newStorageMigrateUpCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   storageMigrateDirectionUp,
		Short: "Perform a migration up",
		Args:  cobra.NoArgs,
		RunE:  newStorageMigrationRunE(true),
	}

	cmd.Flags().IntP("target", "t", 0, "sets the version to migrate to, by default this is the latest version")

	return cmd
}

func newStorageMigrateDownCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   storageMigrateDirectionDown,
		Short: "Perform a migration down",
		Args:  cobra.NoArgs,
		RunE:  newStorageMigrationRunE(false),
	}

	cmd.Flags().IntP("target", "t", 0, "sets the version to migrate to")
	cmd.Flags().Bool("pre1", false, "sets pre1 as the version to migrate to")
	cmd.Flags().Bool("destroy-data", false, "confirms you want to destroy data with this migration")

	return cmd
}