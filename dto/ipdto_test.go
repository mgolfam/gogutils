package dto

import "testing"

func TestLangConfSetRtl(t *testing.T) {
	var l LangConf

	l.SetRtl("true")
	if !l.Rtl {
		t.Error("SetRtl(\"true\") did not set Rtl to true")
	}

	l.SetRtl("false")
	if l.Rtl {
		t.Error("SetRtl(\"false\") did not set Rtl to false")
	}

	// Invalid input leaves the current value untouched.
	l.Rtl = true
	l.SetRtl("not-a-bool")
	if !l.Rtl {
		t.Error("SetRtl(invalid) must not change Rtl")
	}
}

func TestLangConfDefault(t *testing.T) {
	l := LangConf{Lang: "fa", Rtl: true}
	l.Default()
	if l.Lang != "en" || l.Rtl {
		t.Errorf("Default() = %+v, want {Rtl:false Lang:en}", l)
	}
}
