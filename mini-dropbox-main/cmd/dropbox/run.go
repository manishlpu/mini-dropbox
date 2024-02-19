package main

import (
	"time"

	"github.com/manishlpu/assignment/api"
	"github.com/manishlpu/assignment/utils"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Starts running the application server",
		Run: func(cmd *cobra.Command, args []string) {
			srv, err := NewServer()
			if err != nil {
				utils.ErrorLog("Error getting new server:", err)
				return
			}

			s := gocron.NewScheduler(time.Local)
			_, _ = s.Cron("30 1 * * *").Do(func() {
				utils.InfoLog("Cron runs at 1:30 AM every night asynchronously")

				err = api.DeleteInactiveRecords()
				if err != nil {
					utils.ErrorLog("unable to delete records through cron job:", err)
					return
				}

			})
			s.StartAsync()

			StartServer(srv)
		},
	}

	rootCmd.AddCommand(runCmd)
}
