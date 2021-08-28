package output

import (
	"context"
	"errors"
	"fmt"
	"onewayframe/pkg/plugin"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedStatus string

var (
	success = "success"
	failed  = "failed"
)

func testCoreFunc(message map[string]string) error {
	switch expectedStatus {
	case success:
		return nil
	case failed:
		return errors.New("error: failed")
	}

	return nil
}

func TestInitPluginMap(t *testing.T) {
	assert.NotEqual(t, nil, Plugin, "failed to init output plugin map")
}

func TestWrapWithOutputLoopOpenSuccess(t *testing.T) {
	plugin.P2OChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	expectedStatus = success

	defer cancel()
	go WrapWithOutputLoop(ctx, wg, testCoreFunc)()

	plugin.P2OChan <- map[string]string{}
	time.Sleep(1 * time.Second)
	assert.LessOrEqual(t, int32(1), plugin.Metrics.OutputOK, "failed to handle incoming message")
}

func TestWrapWithOutputLoopCancelSuccess(t *testing.T) {
	plugin.P2OChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	alreadyCanceled := false
	go func() {
		time.Sleep(2 * time.Second)
		r := assert.Equal(t, true, alreadyCanceled, "failed to cancel the output loop")
		if r != true {
			fmt.Println("failed to cancel the WrapWithOutputLoop, force exit testing")
			os.Exit(1)
		}
	}()

	cancel()
	WrapWithOutputLoop(ctx, wg, testCoreFunc)()
	alreadyCanceled = true
}

func TestWrapWithOutputLoopCloseSuccess(t *testing.T) {
	plugin.P2OChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, _ := context.WithCancel(context.Background())

	expectedStatus = success
	plugin.P2OChan <- map[string]string{}

	close(plugin.P2OChan)
	WrapWithOutputLoop(ctx, wg, testCoreFunc)()

	assert.Equal(t, true, plugin.OutputStatus.Completed, "failed to run on oneTimeExec mode")
}

func TestWrapWithOutputLoopFailed(t *testing.T) {
	plugin.P2OChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	expectedStatus = failed
	plugin.P2OChan <- map[string]string{}

	defer cancel()
	go WrapWithOutputLoop(ctx, wg, testCoreFunc)()

	time.Sleep(1 * time.Second)
	assert.LessOrEqual(t, int32(1), plugin.Metrics.OutputErr, "failed to make record for output error")
}
