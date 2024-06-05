package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetupLogging(t *testing.T) {
	logger := SetupLogging(true)
	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected DebugLevel, got %v", logger.GetLevel())
	}

	logger = SetupLogging(false)
	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", logger.GetLevel())
	}
}
