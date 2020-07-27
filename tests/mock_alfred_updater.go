package tests

import (
	"errors"
	"time"

	aw "github.com/deanishe/awgo"
)

// ensure MockAlfredUpdater implements Updater
var _ aw.Updater = (*MockAlfredUpdater)(nil)

type MockAlfredUpdater struct {
	updateIntervalCalled  bool
	checkDueCalled        bool
	checkForUpdateCalled  bool
	updateAvailableCalled bool
	installCalled         bool

	checkShouldFail   bool
	installShouldFail bool
}

// UpdateInterval implements Updater.
func (d *MockAlfredUpdater) UpdateInterval(_ time.Duration) {
	d.updateIntervalCalled = true
}

// UpdateAvailable implements Updater.
func (d *MockAlfredUpdater) UpdateAvailable() bool {
	d.updateAvailableCalled = true
	return true
}

// CheckDue implements Updater.
func (d *MockAlfredUpdater) CheckDue() bool {
	d.checkDueCalled = true
	return true
}

// CheckForUpdate implements Updater.
func (d *MockAlfredUpdater) CheckForUpdate() error {
	d.checkForUpdateCalled = true
	if d.checkShouldFail {
		return errors.New("check failed")
	}
	return nil
}

// Install implements Updater.
func (d *MockAlfredUpdater) Install() error {
	d.installCalled = true
	if d.installShouldFail {
		return errors.New("install failed")
	}
	return nil
}
