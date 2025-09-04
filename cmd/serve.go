/*
Copyright © 2025 Adharsh Manikandan <debugslayer@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"servicehub_api/infra"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve the servicehub api",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Println("error getting port", err)
			return
		}

		logger := logrus.New()

		config := infra.Config{
			Port: port,
		}

		srv := infra.NewServer(logger, config)

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("error serving the api: %v", err)
			}
		}()

		log.Println("api running on port", port)

		// block until the signal is received
		<-ch
		log.Println("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("error shutting down server: %v", err)
		}

		log.Println("server shutdown...")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// flag to set the port
	serveCmd.Flags().IntP("port", "p", 8080, "port to serve the api")
}
