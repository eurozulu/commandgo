package commandgo

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

var testVarBool bool
var testVarString string
var testVarURL *url.URL

func testFunc(s string) string {
	return fmt.Sprintf("--%s--", s)
}

func testFuncInt(i int) string {
	return fmt.Sprintf("--%d--", i)
}

func testFuncUrl(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.String()
}

type testStruct struct {
	FieldInt   int
	FieldBool  bool
	FieldFloat float64
}

func (t testStruct) CapitalString(s string) string {
	return strings.ToUpper(s)
}

func (t testStruct) CapitalFields() []string {
	return []string{
		strconv.Itoa(t.FieldInt),
		strconv.FormatBool(t.FieldBool),
		strconv.FormatFloat(t.FieldFloat, 0, 2, 64),
	}
}

func TestCommands_Run_NoCommand(t *testing.T) {
	cmds := Commands{}
	_, err := cmds.Run()
	if err != ErrorNoCommandFound {
		t.Fatalf("expected %v error with empty command line, found %v", ErrorNoCommandFound, err)
	}
	_, err = cmds.Run("test")
	if err != ErrorCommandNotKnown {
		t.Fatalf("expected %v error with empty command map, found %v", ErrorCommandNotKnown, err)
	}
	cmds = Commands{"test": testFunc}
	_, err = cmds.Run()
	if err != ErrorNoCommandFound {
		t.Fatalf("expected %v error with empty command line, found %v", ErrorNoCommandFound, err)
	}
	_, err = cmds.Run("unknown")
	if err != ErrorCommandNotKnown {
		t.Fatalf("expected %v error with unknown command, found %v", ErrorCommandNotKnown, err)
	}
	_, err = cmds.Run("test")
	if err == nil || !strings.HasPrefix(err.Error(), "missing argument") {
		t.Fatalf("expected error %s with known command, no param, found %v", "missing argument", err)
	}
	_, err = cmds.Run("test", "hello")
	if err != nil {
		t.Fatalf("expected no error with known command, found %v", err)
	}
}

func TestCommands_Run_vars(t *testing.T) {
	testVarBool = false
	testVarString = ""
	testVarURL = nil

	cmds := Commands{
		"-flag1": &testVarBool,
		"-flag2": &testVarString,
		"-flag3": &testVarURL,
		"":       testFunc,
	}

	out, err := cmds.Run("-flag1", "-flag2", "hello", "-flag3", "http://www.google.com/")
	if err == nil || !strings.HasPrefix(err.Error(), "missing argument") {
		t.Fatalf("expected error %s with valid command map but no command param. found %v", "missing argument", err)
	}

	out, err = cmds.Run("teststring", "-flag1", "-flag2", "hello", "-flag3", "http://www.google.com/")
	if err != nil {
		t.Fatalf("unexpected error with valid command line and command param.  found %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("unexpected output with valid command line  %v", out)
	}
	if out[0].(string) != "--teststring--" {
		t.Fatalf("unexpected output expected %v, found %v", "--teststring--", out[0])
	}
	if !testVarBool {
		t.Fatalf("unexpected testvarbool, expected %v, found %v", true, testVarBool)
	}
	if testVarString != "hello" {
		t.Fatalf("unexpected testvarstring, expected %v, found %v", "hello", testVarString)
	}

	if testVarURL == nil {
		t.Fatalf("unexpected testvarurl, expected %v, found nil", "http://www.google.com/")
	}

	if testVarURL.String() != "http://www.google.com/" {
		t.Fatalf("unexpected testvarurl, expected %v, found %s", "http://www.google.com/", testVarURL.String())
	}
}

func TestCommands_Run_Fields(t *testing.T) {
	test := &testStruct{}

	cmd := Commands{
		"-b":  &test.FieldBool,
		"-i":  &test.FieldInt,
		"-f":  &test.FieldFloat,
		"cap": test.CapitalString,
	}
	out, err := cmd.Run("cap", "teststring", "-i", "555", "-b", "-f", "0.555")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected output. Expected %d item, found %d", 1, len(out))
	}
	if out[0].(string) != strings.ToUpper("teststring") {
		t.Fatalf("unexpected output. Expected %s, found %s", strings.ToUpper("teststring"), out[0].(string))
	}

	if !test.FieldBool {
		t.Fatalf("unexpected field value. bool expected %v, found %v", true, test.FieldBool)
	}
	if test.FieldInt != 555 {
		t.Fatalf("unexpected field value. int expected %v, found %v", 555, test.FieldInt)
	}
	if test.FieldFloat != 0.555 {
		t.Fatalf("unexpected field value. int expected %v, found %v", 0.555, test.FieldFloat)
	}
}

func TestCommands_Run_Func(t *testing.T) {
	cmd := Commands{
		"dash": testFunc,
	}
	_, err := cmd.Run("dash")
	if err == nil || !strings.HasPrefix(err.Error(), "missing argument") {
		t.Fatalf("expected error %s, found %s", "missing argument", err)
	}
	out, err := cmd.Run("dash", "abc")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if len(out) != 1 || out[0].(string) != "--abc--" {
		t.Fatalf("unexpected output, expected %s, found %v", "--abc---", out[0])
	}

	cmd = Commands{
		"url": testFuncUrl,
	}
	out, err = cmd.Run("url", "nonsense://bla bla")
	if err == nil {
		t.Fatalf("expected error, found none")
	}

	out, err = cmd.Run("url", "http://www.google.com")
	if err != nil {
		t.Fatalf("unexpected error, %s", err)
	}
	if len(out) != 1 || out[0].(string) != "http://www.google.com" {
		t.Fatalf("unexpected output, expected %s, found %v", "http://www.google.com", out[0])
	}
}

func TestCommands_Run_Meth(t *testing.T) {
	ts := &testStruct{}
	cmd := Commands{
		"meth": ts.CapitalString,
	}
	_, err := cmd.Run("meth")
	if err == nil || !strings.HasPrefix(err.Error(), "missing argument") {
		t.Fatalf("expected error %s, found %s", "missing argument", err)
	}
	out, err := cmd.Run("meth", "abc")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if len(out) != 1 || out[0].(string) != "ABC" {
		t.Fatalf("unexpected output, expected %s, found %v", "ABC", out[0])
	}
}

func TestCommands_Run_Default(t *testing.T) {

	cmd := Commands{
		"":    testFunc,
		"num": testFuncInt,
	}
	_, err := cmd.Run()
	if err == nil || !strings.HasPrefix(err.Error(), "missing argument") {
		t.Fatalf("expected error %s, found %s", "missing argument", err)
	}

	_, err = cmd.Run("num")
	if err == nil || !strings.HasPrefix(err.Error(), "missing argument") {
		t.Fatalf("expected error %s, found %s", "missing argument", err)
	}

	out, err := cmd.Run("teststring")
	if err != nil {
		t.Fatalf("unexpected error %s,", err)
	}
	if len(out) != 1 || out[0].(string) != "--teststring--" {
		t.Fatalf("unexpected output.  expected %s, found %s", "--teststring--", out[0])
	}

	_, err = cmd.Run("num", "teststring")
	if err == nil || !strings.HasSuffix(err.Error(), "could not be read as a int") {
		t.Fatalf("expected error %s, found %v", "could not be read as a int", err)
	}
	out, err = cmd.Run("num", "123")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if len(out) != 1 || out[0].(string) != "--123--" {
		t.Fatalf("unexpected output.  expected %v, found %v", "--123--", out[0])
	}
}

func TestCommands_Run_Submaps(t *testing.T) {
	ts := &testStruct{}
	cmds := Commands{
		"-b":    &testVarBool,
		"0func": testFunc,
		"one": Commands{
			"-u":    &testVarURL,
			"1func": testFunc,
			"two": Commands{
				"-s":    &testVarString,
				"2func": ts.CapitalString,
			},
		},
	}

	// level zero, no flags
	out, err := cmds.Run("0func", "teststring")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if len(out) != 1 && out[0].(string) != "teststring" {
		t.Fatalf("unexpected output testing subcommands.  Expected %s, found %v", "teststring", out[0])
	}

	// level zero, valid flag
	testVarBool = false
	out, err = cmds.Run("0func", "teststring", "-b", "true")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if len(out) != 1 && out[0].(string) != "teststring" {
		t.Fatalf("unexpected output testing subcommands.  Expected %s, found %v", "teststring", out[0])
	}
	if !testVarBool {
		t.Fatalf("unexpected flag value testing testvarbool.  Expected %v, found %v", true, testVarBool)
	}

	// level zero, flag out of scope
	testVarString = ""
	out, err = cmds.Run("0func", "teststring", "-s", "hello")
	if err == nil || !strings.HasPrefix(err.Error(), "unexpected flag found") {
		t.Fatalf("expected error unexpected flag found, %v", err)
	}

	// level one, no flags
	out, err = cmds.Run("one", "1func", "teststring")
	if err != nil {
		t.Fatalf("unexpected error testing subcommands, %v", err)
	}

	// level one, valid flags
	testVarURL = nil
	testVarBool = false
	out, err = cmds.Run("one", "1func", "teststring", "-u", "http://www.google.com", "-b", "true")
	if err != nil {
		t.Fatalf("unexpected error testing subcommands, %v", err)
	}
	if testVarURL == nil {
		t.Fatalf("no value assigned to test flag variable testvarurl")
	}
	if testVarURL.String() != "http://www.google.com" {
		t.Fatalf("unexpected value assigned to test flag variable testvarurl.  expected %s, found %s", "http://www.google.com", testVarURL.String())
	}
	if !testVarBool {
		t.Fatalf("unexpected value assigned to test flag variable testvarbool")
	}

	// level one, invalid flags
	out, err = cmds.Run("one", "1func", "teststring", "-s", "teststring")
	if err == nil || !strings.HasPrefix(err.Error(), "unexpected flag found") {
		t.Fatalf("expected error unexpected flag found, %v", err)
	}

}
