package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

func ExeRealTimeOut(ctx context.Context, writer io.Writer, arg ...string) error {
	c := strings.Join(arg, " ")
	cmd := exec.CommandContext(ctx, "bash", "-c", c)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go read(&wg, writer, stdout)
	go read(&wg, writer, stderr)
	wg.Wait()
	return nil
}

func read(wg *sync.WaitGroup, writer io.Writer, std io.ReadCloser) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		readString, err := reader.ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprint(writer, err)
			return
		}
		_, _ = fmt.Fprint(writer, readString)
	}
}
