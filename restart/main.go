package main

import (
	"context"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if false {
		showGracefulRestart()
	}
	showGracefulStop()
}

func showGracefulStop() {
	//服务器
	router := gin.Default()
	router.GET("/", func(context *gin.Context) {
		time.Sleep(5 * time.Second)
		context.String(http.StatusOK, "welcome gin server")
	})
	srv := &http.Server{Addr: ":8080",
		Handler: router}
	go func() {
		//router.Run()
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s \n", err)
		}
	}()

	//关闭信号处理器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//延时执行器
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Sever shutdown", err)
	}
	log.Println("server exiting")
}

func showGracefulRestart() {
	//服务器运行
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "hello gin!")
	})
	if err := endless.ListenAndServe("localhost:8080", router); err != nil {
		log.Fatalf("listen :%s \n", err)
	}
	log.Println("Server exiting")
}
