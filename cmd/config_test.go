package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/harrybrwn/apizza/cmd/internal/cmdtest"
	"github.com/harrybrwn/apizza/pkg/config"
	"github.com/harrybrwn/apizza/pkg/tests"
)

var testconfigjson = `
{
	"name":"joe","email":"nojoe@mail.com",
	"address":{
		"street":"1600 Pennsylvania Ave NW",
		"cityName":"Washington DC",
		"state":"","zipcode":"20500"
	},
	"card":{"number":"","expiration":"","cvv":""},
	"service":"Carryout"
}`

var testConfigOutput = `name: "joe"
email: "nojoe@mail.com"
address:
  street: "1600 Pennsylvania Ave NW"
  cityname: "Washington DC"
  state: ""
  zipcode: "20500"
card:
  number: ""
  expiration: ""
service: "Carryout"
`

func TestConfigStruct(t *testing.T) {
	r := cmdtest.NewRecorder()
	defer r.CleanUp()
	r.ConfigSetup()
	check(json.Unmarshal([]byte(testconfigjson), r.Config()), "json")

	if r.Config().Get("name").(string) != "joe" {
		t.Error("wrong value")
	}
	if err := r.Config().Set("name", "not joe"); err != nil {
		t.Error(err)
	}
	if r.Config().Get("Name").(string) != "not joe" {
		t.Error("wrong value")
	}
	if err := r.Config().Set("name", "joe"); err != nil {
		t.Error(err)
	}
	config.SetNonFileConfig(cfg) // reset the global config for compatability
}

func TestConfigCmd(t *testing.T) {
	r := cmdtest.NewRecorder()
	c := newConfigCmd(r).(*configCmd)
	c.file = true
	if err := c.Run(c.Cmd(), []string{}); err != nil {
		t.Error(err)
	}
	c.file = false
	r.Compare(t, "\n")
	r.ClearBuf()
	c.dir = true
	if err := c.Run(c.Cmd(), []string{}); err != nil {
		t.Error(err)
	}
	r.Compare(t, "\n")
	r.ClearBuf()
	c.dir = false
	c.getall = true
	if err := c.Run(c.Cmd(), []string{}); err != nil {
		t.Error(err)
	}
	r.Compare(t, testConfigOutput)
	r.ClearBuf()
	c.getall = false
	cmdUseage := c.Cmd().UsageString()
	if err := c.Run(c.Cmd(), []string{}); err != nil {
		t.Error(err)
	}
	r.Compare(t, cmdUseage)
	r.ClearBuf()
	if err := c.Run(c.Cmd(), []string{}); err != nil {
		t.Error(err)
	}
	r.Compare(t, c.Cmd().UsageString())
}

func TestConfigEdit(t *testing.T) {
	r := cmdtest.NewRecorder()
	c := newConfigCmd(r).(*configCmd)
	err := config.SetConfig(".apizza/tests", r.Conf)
	if err != nil {
		t.Error(err)
	}
	os.Setenv("EDITOR", "cat")
	c.edit = true

	exp := `{
    "Name": "",
    "Email": "",
    "Address": {
        "Street": "",
        "CityName": "",
        "State": "",
        "Zipcode": ""
    },
    "Card": {
        "Number": "",
        "Expiration": ""
    },
    "Service": "Delivery"
}`

	tests.CompareOutput(t, exp, func() {
		if err = c.Run(c.Cmd(), []string{}); err != nil {
			t.Error(err)
		}
	})
	if err = os.RemoveAll(config.Folder()); err != nil {
		t.Error(err)
	}
	config.SetNonFileConfig(cfg) // for compatibility with old tests
}

func testConfigGet(t *testing.T) {
	c := newConfigGet()
	buf := &bytes.Buffer{}
	c.SetOutput(buf)
	if err := c.Run(c.Cmd(), []string{"email", "name"}); err != nil {
		t.Error(err)
	}
	tests.Compare(t, string(buf.Bytes()), "nojoe@mail.com\njoe\n")
	buf.Reset()
	if err := c.Run(c.Cmd(), []string{}); err == nil {
		t.Error("expected error")
	} else if err.Error() != "no variable given" {
		t.Error("wrong error message, got:", err.Error())
	}
	if err := c.Run(c.Cmd(), []string{"nonExistantKey"}); err == nil {
		t.Error("expected error")
	} else if err.Error() != "cannot find nonExistantKey" {
		t.Error("wrong error message, got:", err.Error())
	}
}

func testConfigSet(t *testing.T) {
	c := newConfigSet() //.(*configSetCmd)
	if err := c.Run(c.Cmd(), []string{"name=someNameOtherThanJoe"}); err != nil {
		t.Error(err)
	}
	if cfg.Name != "someNameOtherThanJoe" {
		t.Error("did not set the name correctly")
	}
	if err := c.Run(c.Cmd(), []string{}); err == nil {
		t.Error("expected error")
	} else if err.Error() != "no variable given" {
		t.Error("wrong error message, got:", err.Error())
	}
	if err := c.Run(c.Cmd(), []string{"nonExistantKey=someValue"}); err == nil {
		t.Error("expected error")
	}
	if err := c.Run(c.Cmd(), []string{"badformat"}); err == nil {
		t.Error(err)
	} else if err.Error() != "use '<key>=<value>' format (no spaces), use <key>='-' to set as empty" {
		t.Error("wrong error message, got:", err.Error())
	}
}
