package rasaxctl

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func (r *RasaXCTL) ConnectRasa() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	args := "run --enable-api --debug"
	cmd := exec.Command("rasa", strings.Split(args, " ")...)
	stdout, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err = cmd.Start(); err != nil {
		panic(err)
	}

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// Get real-time output from the pipe to the terminal and print
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Printf("(rasa production) %s", string(tmp))
		if err != nil {
			break
		}
	}

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
	if err = cmd.Wait(); err != nil {
		panic(err)
	}
	fmt.Println("done", cmd.ProcessState)
	return nil
}
