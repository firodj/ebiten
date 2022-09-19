//go:build !ios && !nintendosdk
// +build !ios,!nintendosdk

package gamepad

import (
	"hash/crc32"
	"math"
	"strings"
	"time"
	"unsafe"
)

type dualShock4 struct {
	gamepadId uint32
	device    _IOHIDDeviceRef
	busType   string
	timer     *time.Timer
}

const (
	kSonyProduct05c4 uint32 = 0x054c05c4
	kSonyProduct09cc uint32 = 0x054c09cc
	kScufProduct7725 uint32 = 0x2e957725

	kBusTypeUSB       = "usb"
	kBusTypeBluetooth = "bluetooth"

	kReportId05 = 0x05
	kReportId11 = 0x11

	kRumbleMagnitudeMax = 0xff
)

func isDualShock4(vendor, product uint32) bool {
	switch (vendor <<  16) | product {
	case kSonyProduct05c4:
		return true
	case kSonyProduct09cc:
		return true
	case kScufProduct7725:
		return true
	}
	return false
}

func newDualShock4(device _IOHIDDeviceRef, vendor, product uint32) *dualShock4 {
	return &dualShock4{
		device: device,
		gamepadId: (vendor <<  16) | product,
		busType: queryBusType(device),
	}
}

func queryBusType(device _IOHIDDeviceRef) string {
	var transport string
	if prop := _IOHIDDeviceGetProperty(_IOHIDDeviceRef(device), _CFStringCreateWithCString(kCFAllocatorDefault, kIOHIDTransportKey, kCFStringEncodingUTF8)); prop != 0 {
		var cstr [256]byte
		_CFStringGetCString(_CFStringRef(prop), cstr[:], kCFStringEncodingUTF8)
		transport = strings.TrimRight(string(cstr[:]), "\x00")

		if (transport == kIOHIDTransportUSBValue) {
      return kBusTypeUSB
		}
    if (transport == kIOHIDTransportBluetoothValue ||
        transport == kIOHIDTransportBluetoothLowEnergyValue) {
      return kBusTypeBluetooth
    }
	}
	return ""
}

func computeDS4checksum(report []byte) {
	data := []byte{0xa2}
	data = append(data, report[:len(report)-4]...)
	crc := crc32.ChecksumIEEE(data)

	report[len(report)-4] = (byte)(crc & 0xff);
	report[len(report)-3] = (byte)((crc >> 8) & 0xff);
	report[len(report)-2] = (byte)((crc >> 16) & 0xff);
	report[len(report)-1] = (byte)((crc >> 24) & 0xff);
}

func (ds4 *dualShock4) vibrate(duration time.Duration, strongMagnitude float64, weakMagnitude float64) {
	if duration != 0 && (strongMagnitude != 0 || weakMagnitude != 0) {
		if ds4.timer != nil {
			ds4.timer.Stop()
		}
		ds4.timer = time.AfterFunc(duration, func() {
			ds4.vibrate(0, 0, 0)
		})
	}

	if ds4.busType == kBusTypeBluetooth && ds4.gamepadId != kScufProduct7725 {
		ds4.vibrateBluetooth(strongMagnitude, weakMagnitude)
		return
	}
	ds4.vibrateUSB(strongMagnitude, weakMagnitude)
}

func (ds4 *dualShock4) vibrateUSB(strongMagnitude float64, weakMagnitude float64) bool {
	control_report := make([]byte, 32)
	control_report[0] = kReportId05
	control_report[1] = 0x01
	control_report[4] = byte(math.Round(weakMagnitude * kRumbleMagnitudeMax))
	control_report[5] = byte(math.Round(strongMagnitude * kRumbleMagnitudeMax))
	return _IOHIDDeviceSetReport(ds4.device, kIOHIDReportTypeOutput, int(control_report[0]), unsafe.Pointer(&control_report[0]), len(control_report)) == kIOReturnSuccess
}

func (ds4 *dualShock4) vibrateBluetooth(strongMagnitude float64, weakMagnitude float64) bool {
	control_report := make([]byte, 78)
	control_report[0] = kReportId11;
	control_report[1] = 0xc0;  // unknown
	control_report[2] = 0x20;  // unknown
	control_report[3] = 0xf1;  // motor only, don't update LEDs
	control_report[4] = 0x04;  // unknown
	control_report[6] = byte(math.Round(weakMagnitude * kRumbleMagnitudeMax))
	control_report[7] = byte(math.Round(strongMagnitude * kRumbleMagnitudeMax))
	control_report[21] = 0x43;  // volume left
	control_report[22] = 0x43;  // volume right
	control_report[24] = 0x4d;  // volume speaker
	control_report[25] = 0x85;  // unknown

	computeDS4checksum(control_report)
	return _IOHIDDeviceSetReport(ds4.device, kIOHIDReportTypeOutput, int(control_report[0]), unsafe.Pointer(&control_report[0]), len(control_report)) == kIOReturnSuccess
}