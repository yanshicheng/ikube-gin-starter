package cmd

import (
	"github.com/spf13/cobra"
)

var (
	db      string
	migrate bool
)

var dbCommand = &cobra.Command{
	Use:   "db",
	Short: "db console",
	Long:  "db console",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 初始化全局变量
		return nil
	},
}

func init() {
	rootCommand.AddCommand(dbCommand)
	dbCommand.Flags().StringVarP(&db, "database", "d", "default", "database")
	dbCommand.Flags().BoolVarP(&migrate, "migrate", "m", false, "force syncdb")
}
