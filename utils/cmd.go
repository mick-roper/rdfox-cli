package utils

import "github.com/spf13/cobra"

type rootFlags struct {
	Server   string
	Protocol string
	Role     string
	Password string
	Version  int
}

func RootCommandFlags(cmd *cobra.Command) *rootFlags {
	server := cmd.Flags().Lookup("server").Value.String()
	protocol := cmd.Flags().Lookup("protocol").Value.String()
	role := cmd.Flags().Lookup("role").Value.String()
	password := cmd.Flags().Lookup("password").Value.String()
	version, err := cmd.Flags().GetInt("server-version")
	if err != nil {
		panic(err)
	}
	return &rootFlags{server, protocol, role, password, version}
}
