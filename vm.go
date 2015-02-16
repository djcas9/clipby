package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-mruby"
)

type VM struct {
	Mrb           *mruby.Mrb
	PluginClasses []string
}

func (self *VM) Close() {
	self.Close()
}

func (self *VM) PluginNames() []string {
	result, err := self.Mrb.LoadString("Clipby.constants.select {|c| Class === Clipby.const_get(c)}")

	if err != nil {
		log.Fatal(err.Error())
	}

	var classes []string
	data := strings.Replace(strings.Replace(result.String(), "[", "", 1), "]", "", 1)
	array := strings.Split(data, ", ")

	for _, class := range array {
		c := class[1:]
		classes = append(classes, c)
	}

	self.PluginClasses = classes

	return classes
}

func (self *VM) Init() {
	addFunc := func(m *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		args := m.GetArgs()
		return mruby.Int(args[0].Fixnum() + args[1].Fixnum()), nil
	}

	logOutput := func(m *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		args := m.GetArgs()

		var pluginName string

		if len(self.String()) <= 0 {
			pluginName = "N/A"
		} else {
			pluginName = self.String()
		}

		output := fmt.Sprintf("%s STDOUT: %s", pluginName, args[0].String())
		log.Debug(output)

		return mruby.Int(1), nil
	}

	class := self.Mrb.DefineClass("Plugin", nil)
	class.DefineClassMethod("add", addFunc, mruby.ArgsReq(2))
	class.DefineClassMethod("log", logOutput, mruby.ArgsReq(1))
}

func (self *VM) Run(cb CBType) {

	for _, class := range self.PluginClasses {

		//
		// TODO: mruby has a max token len. We can split the cb data
		// and send as a multi-part array or try to compress the data.
		//
		b := hex.EncodeToString([]byte(cb.Data))
		code := fmt.Sprintf(`Clipby::%s.run('%s', ["%s"])`, class, cb.Type, b)

		result, err := self.Mrb.LoadString(code)

		if err != nil {
			log.Fatal(err.Error())
		}

		data := result.String()

		if len(data) > 0 && *Debug {
			log.Debugf("Clipby::%s Return: %s", class, result.String())
		}

		self.Mrb.FullGC()
	}

}

func (self *VM) Load(plugin *Plugin) {
	context := mruby.NewCompileContext(self.Mrb)
	defer context.Close()

	parser := mruby.NewParser(self.Mrb)
	defer parser.Close()

	if _, err := parser.Parse(plugin.Code, context); err != nil {
		panic(err.Error())
	}

	code := parser.GenerateCode()

	if _, err := self.Mrb.Run(code, nil); err != nil {
		panic(err.Error())
	}
}
