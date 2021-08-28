package input

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jdjgya/service-frame/pkg/oneway/plugin"
	"github.com/stretchr/testify/assert"
)

var expectedStatus string

var (
	success  = "success"
	failed   = "failed"
	interval = 0 * time.Second
)

func testCoreFunc() ([]byte, error) {
	switch expectedStatus {
	case success:
		return []byte(success), nil
	case failed:
		return []byte(failed), errors.New("error: failed")
	}

	return []byte{}, nil
}

func TestInitPluginMap(t *testing.T) {
	assert.NotEqual(t, nil, Plugin, "failed to init input plugin map")
}

func TestWrapWithInputLoopOpenSuccess(t *testing.T) {
	plugin.I2TChan = make(chan []byte, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	expectedStatus = success

	defer cancel()
	go WrapWithInputLoop(ctx, wg, testCoreFunc, interval)()

	msg := <-plugin.I2TChan
	assert.LessOrEqual(t, int32(1), plugin.Metrics.InputOK, "failed to handle inputted message")
	assert.Equal(t, "success", string(msg), "")
}

func TestWrapWithInputLoopCancelSuccess(t *testing.T) {
	plugin.I2TChan = make(chan []byte, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	alreadyCanceled := false
	go func() {
		time.Sleep(2 * time.Second)
		r := assert.Equal(t, true, alreadyCanceled, "failed to cancel the input loop")
		if r != true {
			fmt.Println("failed to cancel the WrapWithInputLoop, force exit testing")
			os.Exit(1)
		}
	}()

	cancel()
	WrapWithInputLoop(ctx, wg, testCoreFunc, interval)()
	alreadyCanceled = true
}

func TestWrapWithInputLoopCloseSuccess(t *testing.T) {
	plugin.I2TChan = make(chan []byte, 5)

	wg := &sync.WaitGroup{}
	ctx, _ := context.WithCancel(context.Background())

	plugin.IsOneTimeExec = true

	expectedStatus = success
	WrapWithInputLoop(ctx, wg, testCoreFunc, interval)()

	time.Sleep(1 * time.Second)
	assert.Equal(t, true, plugin.InputStatus.Completed, "failed to run on oneTimeExec mode")
}

func TestWrapWithInputLoopFailed(t *testing.T) {
	plugin.I2TChan = make(chan []byte, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	expectedStatus = failed
	plugin.I2TChan <- []byte("test msg")

	defer cancel()
	go WrapWithInputLoop(ctx, wg, testCoreFunc, interval)()

	time.Sleep(1 * time.Second)
	assert.LessOrEqual(t, int32(1), plugin.Metrics.InputErr, "failed to make record for input error")
}
