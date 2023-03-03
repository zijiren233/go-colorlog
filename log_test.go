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
		Infof("this is info formate: %d", L_Info)
		Warning("this is warning")
		Warningf("this is warning formate: %d", L_Warning)
		Error("this is err")
		Errorf("this is err formate: %d", L_Error)
		Fatal("this is fatal")
		Fatalf("this is fatal formate: %d", L_Fatal)
	}
}
