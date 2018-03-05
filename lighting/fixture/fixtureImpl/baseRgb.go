// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import (
	"channellight"
	hkAccessory "github.com/brutella/hc/accessory"
	hkService "github.com/brutella/hc/service"
	"github.com/op/go-logging"
	"lighting/channelUpdater"
	"lighting/lights"
	"lighting/store"
	"time"
)

var log = logging.MustGetLogger("baseRGBFixture")

type baseRGBFixture struct {
	colorFixtureChannels colorFixtureChannels
	lightbulb            *hkService.Lightbulb
	accessory            *hkAccessory.Accessory
	lightModel           channellight.ChannelLight
}

type colorFixtureChannels struct {
	fader    *lights.Address
	red      *lights.Address
	green    *lights.Address
	blue     *lights.Address
	white    *lights.Address
	amber    *lights.Address
	uv       *lights.Address
}

func newBaseRGBFixture(fixture FixtureImpl, channels colorFixtureChannels) *baseRGBFixture {
	baseRGBFixture := &baseRGBFixture{
		colorFixtureChannels: channels,
	}

	if channels.fader != nil {
		var lightModel channellight.ChannelLight

		lightbulb := hkService.NewLightbulb()
		accessory := hkAccessory.New(hkAccessory.Info{Name: fixture.GetName()}, hkAccessory.TypeLightbulb)

		accessory.AddService(lightbulb.Service)

		fixture.SetHomeKitAccessory(accessory)

		if channels.red != nil &&
		   channels.green != nil &&
		   channels.blue != nil {
			if channels.white != nil {
				lightModel = &channellight.SevenChannelLight{}
			} else {
				lightModel = &channellight.FourChannelLight{}
			}
		}

		baseRGBFixture.lightModel = lightModel
		baseRGBFixture.lightbulb = lightbulb
		baseRGBFixture.accessory = accessory

		lightbulb.On.OnValueRemoteUpdate(func(on bool) {
			baseRGBFixture.syncColorsForLight()
		})

		lightbulb.Brightness.OnValueRemoteUpdate(func(b int) {
			baseRGBFixture.syncColorsForLight()
		})

		lightbulb.Hue.OnValueRemoteUpdate(func(h float64) {
			baseRGBFixture.syncColorsForLight()
		})

		lightbulb.Saturation.OnValueRemoteUpdate(func(s float64) {
			baseRGBFixture.syncColorsForLight()
		})

		store.Subscribe(baseRGBFixture.ValueChange)
	}

	return baseRGBFixture
}

func (this *baseRGBFixture) syncColorsForLight() {
	if this.lightModel == nil {
		return
	}
	this.lightModel.SetColor(this.lightbulb)

	switch v := this.lightModel.(type) {
	case *channellight.SevenChannelLight:
		outputColor := v.GetOutputColor()

		this.SetFaderValue(lights.Value(outputColor.Brightness), time.Duration(0))
		this.SetRedValue(lights.Value(outputColor.Red), time.Duration(0))
		this.SetGreenValue(lights.Value(outputColor.Green), time.Duration(0))
		this.SetBlueValue(lights.Value(outputColor.Blue), time.Duration(0))
		this.SetWhiteValue(lights.Value(outputColor.White), time.Duration(0))
		if this.IsAmberAvailable() {
			this.SetAmberValue(lights.Value(outputColor.Amber), time.Duration(0))
		}
		if this.IsUvAvailable() {
			this.SetUvValue(lights.Value(outputColor.Uv), time.Duration(0))
		}
	case *channellight.FourChannelLight:
		outputColor := v.GetOutputColor()

		this.SetFaderValue(lights.Value(outputColor.Brightness), time.Duration(0))
		this.SetRedValue(lights.Value(outputColor.Red), time.Duration(0))
		this.SetGreenValue(lights.Value(outputColor.Green), time.Duration(0))
		this.SetBlueValue(lights.Value(outputColor.Blue), time.Duration(0))
	}
}

func (this *baseRGBFixture) doGetValue(description string, channel *lights.Address) lights.Value {
	if channel == nil {
		log.Infof("Requested '%s' channel when not implemented",description)
		return 0
	}
	return store.GetValue(*channel)
}

func (this *baseRGBFixture) IsFaderAvailable() bool {
	return this.colorFixtureChannels.fader != nil
}

func (this *baseRGBFixture) IsWhiteAvailable() bool {
	return this.colorFixtureChannels.white != nil
}

func (this *baseRGBFixture) IsAmberAvailable() bool {
	return this.colorFixtureChannels.amber != nil
}

func (this *baseRGBFixture) IsUvAvailable() bool {
	return this.colorFixtureChannels.uv != nil
}

func (this *baseRGBFixture) GetFaderValue() lights.Value {
	return this.doGetValue("fader", this.colorFixtureChannels.fader)
}

func (this *baseRGBFixture) GetRedValue() lights.Value {
	return this.doGetValue("red", this.colorFixtureChannels.red)
}

func (this *baseRGBFixture) GetGreenValue() lights.Value {
	return this.doGetValue("green", this.colorFixtureChannels.green)
}

func (this *baseRGBFixture) GetBlueValue() lights.Value {
	return this.doGetValue("blue", this.colorFixtureChannels.blue)
}

func (this *baseRGBFixture) GetWhiteValue() lights.Value {
	return this.doGetValue("white", this.colorFixtureChannels.white)
}

func (this *baseRGBFixture) GetAmberValue() lights.Value {
	return this.doGetValue("amber", this.colorFixtureChannels.amber)
}

func (this *baseRGBFixture) GetUvValue() lights.Value {
	return this.doGetValue("uv", this.colorFixtureChannels.uv)
}

func (this *baseRGBFixture) doSetValue(description string, channel *lights.Address, value lights.Value, fade time.Duration) {
	if channel == nil {
		log.Errorf("(%s) Attempted to set '%s' channel when not implemented", this.accessory.Info.Name, description)
		return
	}

	channelUpdater.GetChannelUpdater(*channel).UpdateValueWithFade(this.doGetValue(description, channel), value, fade)
}

func (this *baseRGBFixture) SetFaderValue(value lights.Value, fade time.Duration) {
	this.doSetValue("fader", this.colorFixtureChannels.fader, value, fade)
}

func (this *baseRGBFixture) SetRedValue(value lights.Value, fade time.Duration) {
	this.doSetValue("red", this.colorFixtureChannels.red, value, fade)
}

func (this *baseRGBFixture) SetGreenValue(value lights.Value, fade time.Duration) {
	this.doSetValue("green", this.colorFixtureChannels.green, value, fade)
}

func (this *baseRGBFixture) SetBlueValue(value lights.Value, fade time.Duration) {
	this.doSetValue("blue", this.colorFixtureChannels.blue, value, fade)
}

func (this *baseRGBFixture) SetWhiteValue(value lights.Value, fade time.Duration) {
	this.doSetValue("white", this.colorFixtureChannels.white, value, fade)
}

func (this *baseRGBFixture) SetAmberValue(value lights.Value, fade time.Duration) {
	this.doSetValue("amber", this.colorFixtureChannels.amber, value, fade)
}

func (this *baseRGBFixture) SetUvValue(value lights.Value, fade time.Duration) {
	this.doSetValue("uv", this.colorFixtureChannels.uv, value, fade)
}

func (this *baseRGBFixture) SetColor(red, green, blue lights.Value, fade time.Duration) {
	this.SetRedValue(red, fade)
	this.SetGreenValue(green, fade)
	this.SetBlueValue(blue, fade)
}

func (this *baseRGBFixture) ValueChange(change store.ValuesChange) {
	if this.lightbulb == nil {
		return
	}

	//switch change.Channel {
	//case *this.colorFixtureChannels.fader, *this.colorFixtureChannels.red, *this.colorFixtureChannels.green, *this.colorFixtureChannels.blue, *this.colorFixtureChannels.white, *this.colorFixtureChannels.amber, *this.colorFixtureChannels.uv:
	//	log.Infof("Colour channel %d changed to %d (HomeKit update would happen here)", change.Channel, change.Value)
	//}
}
