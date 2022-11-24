package colorlog

import (
	"testing"
	"time"
)

func TestLog(test *testing.T) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		Info("this is info")
		Infof("this is info formate: %d", 1)
		Error("this is err")
		Errorf("this is err formate: %d", 2)
		Fatal("this is fatal")
		Fatalf("this is fatal formate: %d", 3)
	}
}
