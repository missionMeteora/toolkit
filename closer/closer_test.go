package closer

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	c := New()

	go func() {
		c.Close(errors.New("HALP"))
	}()

	if err := c.Wait(); err.Error() != "HALP" {
		t.Error("Invalid error provided:", err)
	}
}

func TestManual(t *testing.T) {
	c := New()
	end := make(chan error, 1)

	go func() {
		c.Wait()
		end <- nil
	}()

	go func() {
		fmt.Println("Please attempt to kill the program (CTRL+C/CMD+C)")
		time.Sleep(time.Second * 10)
		fmt.Println("Sleep done..")
		end <- errors.New("interrupt signal not detected")
	}()

	if err := <-end; err != nil {
		t.Error(err)
	}
}
