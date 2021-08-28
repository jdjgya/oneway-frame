package transit

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

func testCoreFunc(msgs []byte) ([]map[string]string, error) {
	returnMsged := []map[string]string{}
	switch expectedStatus {
	case success:
		return returnMsged, nil
	case failed:
		return returnMsged, errors.New("error: failed")
	}

	return []map[string]string{}, nil
}

func TestInitPluginMap(t *testing.T) {
	assert.NotEqual(t, nil, Plugin, "failed to init transit plugin map")
}

func TestWrapWithTransitLoopCancelSuccess(t *testing.T) {
	plugin.I2TChan = make(chan []byte, 5)
	plugin.T2PChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	alreadyCanceled := false
	go func() {
		time.Sleep(2 * time.Second)
		r := assert.Equal(t, true, alreadyCanceled, "failed to close transit loop")
		if r != true {
			fmt.Println("failed to cancel the WrapWithTransitLoop, force exit testing")
			os.Exit(1)
		}
	}()

	cancel()
	plugin.I2TChan <- []byte("test msg")
	WrapWithTransitLoop(ctx, wg, testCoreFunc)()
	alreadyCanceled = true
}

func TestWrapWithTransitLoopCloseSuccess(t *testing.T) {
	plugin.I2TChan = make(chan []byte, 5)
	plugin.T2PChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, _ := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3 * time.Second)
		close(plugin.I2TChan)
	}()

	expectedStatus = success
	WrapWithTransitLoop(ctx, wg, testCoreFunc)()

	time.Sleep(1 * time.Second)
	assert.Equal(t, true, plugin.TransitStatus.Completed, "failed to close I2TChan by transit loop")
}
