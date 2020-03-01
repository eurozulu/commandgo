package mainline_test

import (
	"fmt"
	"github.com/eurozulu/mainline"
	"net/url"
	"strings"
	"testing"
	"time"
)

type SingleFieldTest struct {
	BoolFlag bool
}

type MultiFieldTest struct {
	BoolFlag   bool
	StringFlag string
	IntFlag    int
}

type ComplexFieldTest struct {
	DurationFlag time.Duration `flag:"duration, d"`
	TimeFlag     time.Time     `flag:"time, t"`
	URLFlag      *url.URL      `flag:"url, u"`
}

type CommandFuncTest struct {
	BoolFlag   bool                        `flag:"bool"`
	StringFlag string                      `flag:"str"`
	IntFlag    int                         `flag:"int"`
	Command    func(string, time.Duration) `command:"runthis"`
	Command2   func(...string)             `command:"runthat"`
}

func (c CommandFuncTest) RunThisCommand(bla string, d time.Duration) {
	fmt.Printf("bla %s for a duration of %v", bla, d)
}

func RunThatCommand(args ...string) {
	fmt.Printf("%d arguments:  %v", len(args), args)
}

func TestDecoder_CommandFunc(t *testing.T) {
	st := CommandFuncTest{}
	st.Command = st.RunThisCommand
	st.Command2 = RunThatCommand

	if err := mainline.NewDecoder([]string{"-bool", "--str", "24h", "-int", "42", "runthis bla 24m"}).Decode(&st); err != nil {
		t.Fatal(err)
	}

	if err := mainline.NewDecoder([]string{"-bool", "--str", "24h", "-int", "42", "runthat"}).Decode(&st); err != nil {
		t.Fatal(err)
	}

}

func TestDecoder_FieldTags(t *testing.T) {
	st := ComplexFieldTest{}
	if err := mainline.NewDecoder([]string{"-d", "24h", "--t", "2006-01-02T15:04:05Z"}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.DurationFlag == 0 {
		t.Fatalf("Expected duration set but found zero using tag names")
	}
	if st.TimeFlag.IsZero() {
		t.Fatalf("Expected time set but found zero time using tag names")
	}

	if err := mainline.NewDecoder([]string{"-duration", "24h", "--time", "2006-01-02T15:04:05Z"}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.DurationFlag == 0 {
		t.Fatalf("Expected duration set but found zero using tag names")
	}
	if st.TimeFlag.IsZero() {
		t.Fatalf("Expected time set but found zero time using tag names")
	}

}

func TestDecoder_BadArgs(t *testing.T) {
	st := SingleFieldTest{}
	// empty should pass OK
	if err := mainline.NewDecoder([]string{}).Decode(&st); err != nil {
		t.Fatal(err)
	}

	err := mainline.NewDecoder([]string{"--abadbadflag", "and its bad bad value"}).Decode(&st)
	if err == nil {
		t.Fatal(fmt.Errorf("Expected err passing unknown flag"))
	}

	if !strings.Contains(err.Error(), "abadbadflag") {
		t.Fatal(fmt.Errorf("Expected err to mention unknown flag name"))
	}
}

func TestDecoder_DecodeSingleField(t *testing.T) {
	st := SingleFieldTest{}
	if err := mainline.NewDecoder([]string{}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.BoolFlag {
		t.Errorf("expected ampty Args to have false bool, found true")
	}

	if err := mainline.NewDecoder([]string{"--boolflag"}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if !st.BoolFlag {
		t.Errorf("expected Args with boolflag set to have true bool, found false")
	}
}

func TestDecoder_DecodeMultiField(t *testing.T) {
	st := MultiFieldTest{}
	if err := mainline.NewDecoder([]string{}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.BoolFlag {
		t.Errorf("expected empty Args to have false bool, found true")
	}
	if st.IntFlag != 0 {
		t.Errorf("expected empty Args to have zero int, found %d", st.IntFlag)
	}
	if st.StringFlag != "" {
		t.Errorf("expected empty Args to have empty string, found %s", st.StringFlag)
	}

	if err := mainline.NewDecoder([]string{"--boolflag", "--stringflag", "stringvalue", "--intflag", "99"}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if !st.BoolFlag {
		t.Errorf("expected Args with bool flag set to have true bool, found false")
	}
	if st.IntFlag != 99 {
		t.Errorf("expected Args to have 99 int, found %d", st.IntFlag)
	}
	if st.StringFlag != "stringvalue" {
		t.Errorf("expected Args to have string value 'stringvalue', found %s", st.StringFlag)
	}

}

func TestDecoder_DecodeComplexField(t *testing.T) {
	st := ComplexFieldTest{}
	if err := mainline.NewDecoder([]string{}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.DurationFlag != 0 {
		t.Errorf("expected empty Args to have zero duration, found %v", st.DurationFlag)
	}
	if !st.TimeFlag.IsZero() {
		t.Errorf("expected empty Args to have zero time, found %v", st.TimeFlag)
	}
	if st.URLFlag != nil {
		t.Errorf("expected empty Args to have nil URL, found %v", st.URLFlag)
	}

	if err := mainline.NewDecoder([]string{
		"--durationflag", "30h", "-timeflag", "2006-01-02T15:04:05Z", "--urlflag", "http://www.spoofer.org/haha?one=1"}).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.DurationFlag == 0 {
		t.Errorf("expected Args to have set duration, found %v", st.DurationFlag)
	}
	if st.TimeFlag.IsZero() {
		t.Errorf("expected Args to have set time, found %v", st.TimeFlag)
	}
	if st.URLFlag == nil {
		t.Errorf("expected Args to have set url, found nil")
	}
	if st.URLFlag.Host != "www.spoofer.org" {
		t.Errorf("expected Args to have set url host 'www.spoofer.org', found %s", st.URLFlag.Host)
	}

	fmt.Println(st.TimeFlag)
}
