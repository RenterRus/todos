package main

import (
	"fmt"
	"github.com/spf13/cobra"
	grpc2 "google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"todoes/internal/grpc"
	"todoes/internal/http"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

var rootCmd = &cobra.Command{
	Use:   "todos",
	Short: "Root command",
	Long: "This command is the main one, that is, it is entered as an entry point to the CLI application, for example, " +
		"as the main command in Git (git merge..., git pull..., etc)",
	Run: RunOrchestrator,
}


func init() {
	rootCmd.PersistentFlags().String("http", "127.0.0.1:9999", "HTTP Server addr")
	rootCmd.PersistentFlags().String("grpc", ":9998", "GRPC Server addr")
	rootCmd.PersistentFlags().String("db", "todos", "DB name")
	rootCmd.PersistentFlags().String("table", "todo", "DB table name")

	rootCmd.Execute()
}

func main(){

}


func RunOrchestrator(cmd *cobra.Command, args []string){
	var wg sync.WaitGroup
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	addrGrpc, _ := cmd.Flags().GetString("grpc")
	dbname, _ := cmd.Flags().GetString("db")
	dbtable, _ := cmd.Flags().GetString("table")
	l, err := net.Listen("tcp", addrGrpc)
	if err != nil {
		panic(err.Error())
	}

	var httpServer *http.HTTPServer
	grpcServer := grpc2.NewServer()
	todoes.RegisterTodoesServer(grpcServer, &grpc.GrpcServer{})


	sqlite.DBClient = sqlite.Initial(dbname, dbtable)
	wg.Add(3)
	go func() {
		fmt.Println("GRPC Server starting")
		defer wg.Done()
		if err := grpcServer.Serve(l); err != nil {
			panic(err.Error())
		}
	}()
	go func() {
		defer wg.Done()
		fmt.Println("HTTP Server starting")
		addr, _ := cmd.Flags().GetString("http")
		httpServer = http.NewServer(addr, addrGrpc)
		httpServer.Start()
	}()
	go func() {
		defer wg.Done()
		<-sign
		grpcServer.GracefulStop()
		fmt.Println("GRPC Shutdown")
		httpServer.GraceShutdown()
		fmt.Println("HTTP Shutdown")
		sqlite.DBClient.DisableConnect()
		fmt.Println("DB Connection is closed")
	}()

	wg.Wait()
}