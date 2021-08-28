package process

import (
	"context"
	"errors"
	"onewayframe/pkg/plugin"
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

func testCoreFunc(msg map[string]string, isChnOpen bool) (map[string]string, error) {
	switch expectedStatus {
	case success:
		return msg, nil
	case failed:
		return msg, errors.New("error: failed")
	}

	return msg, nil
}

func TestInitPluginMap(t *testing.T) {
	assert.NotEqual(t, nil, Plugin, "failed to init process plugin map")
}

func TestWrapWithProcessLoopCloseSuccess(t *testing.T) {
	plugin.T2PChan = make(chan map[string]string, 5)
	plugin.P2OChan = make(chan map[string]string, 5)

	wg := &sync.WaitGroup{}
	ctx, _ := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3 * time.Second)
		close(plugin.T2PChan)
	}()

	expectedStatus = success
	WrapWithProcessLoop(ctx, wg, testCoreFunc)()

	time.Sleep(1 * time.Second)
	assert.Equal(t, true, plugin.ProcessStatus.Completed, "failed to close T2PChan by process loop")
}
