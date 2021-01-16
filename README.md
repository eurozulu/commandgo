# Mainline

### Command line arguments object mapper
Simplifies writing command line tools by mapping command line arguments into structure values and method parameters.
Maps "commands" (The first argument") into struct methods, parsing the ramaining arguments into the parameters for that
method. Simply define the method and the parameters it requires, and arguments will be mapped into those parameters.

##### Goal
The goal of this parser is to simplfy the boiler work of mapping arguments into commands, parameters and flags.  
Leading Command line arguments are mapped to methods on generic structs, with any following arguments being parsed into
the parameters for that method.  
Flags from the command line are mapped to public Fields in the same struct.

##### Usage

To create a simple, two command interface:  
`user  param1 param2 param3 --name john -timeout 24h -d`
`server param1 -timeout 24h -d`

Create a struct containing methods to be run when one of these commands are called  
UserInfo and ServerInfo

Add fields for the named argument flags required.
* name string
* timeout time.Duration
* debug bool

e.g.

```
type MyArgs struct {  
    Name    string  
    Timeout time.Duration   `flag:"t,to"`
    Debug   bool            `flag:"d"`
}
func (ma MyArgs) UserInfo(s ...string) error {
    ...
}
func (ma MyArgs) ServerInfo(s string) error {
    ...
}

```

To use the struct, in the application `main()`:

```
func main() {
    cmds := mainline.Commands{
		"user":                   MyArgs.UserInfo,
		"server":             MyArgs.ServerInfo,
	}

	out, err := cmds.Run(os.Args...)
	if err != nil {
		fmt.Println(err)
	}
	if out != nil {
		for _, o := range out {
			fmt.Printf("%v\n", o)
		}
	}
```
  
Flags can appear in any order.  All flags, with the exception of bool types must have a following argument as its value.  
This value is converted to the relevant data type for the Field.
Booleans MAY have a value, if it is parsable as a bool.  If they have a following argument which is not parsable as bool, that value is ignored by the bool flag.
Bool flag are True when they are present, unless they are followed by a 'false' value.

Most other data types are supported, all the base types, int64, float32/64, bool, string etc.  
Slices are supported. Maps not yet.  
certain structs are supported:

+ Those implementing the [json.UnmarshalJSON](https://golang.org/pkg/encoding/json/#example__customMarshalJSON)
  interface
+ Those supporting [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler) interface
+ Date, Duration and url.URL

#### Tags

Fields may be tagged to specify alternative names for the flag using standard go tagging. e.g.

```
type MyAppConfig struct {  
   Name     string           `flag:"nom, naam, n"`  
   Timeout  time.Duration    `flag:"t"`
   Debug    bool             `flag:"db"`
}
``` 

Using these flags, the `Name` field could be set with any of the following command line argiments:

+ --name john
+ -n john
+ -nom "alice cooper"

Tagging a field with a '-' `flag:"-"` will hide that field from the argument parsing.

There is no distiction between the double dash and single dash for flags.  "-" is the same as "--"

#### Command alias

By default, Commands are mapped (case insensatively) to the method. Alternative names can be provided in the same
command map by simply repeating the same emthod under an alternative key:  
`    cmds := mainline.Commands{
"userinfo, ui":             MyArgs{},
"serverinfo, si, svr":      MyArgs{}, }
`  
The first entry in the comma list should always be the method name. And following act as aliases for that method.  
The above example, 'ui' command can be used to call the UserInfo method.

##### Hidden Methods

Methods can be hidden using a preceeding dash.
`    cmds := mainline.Commands{
"user":       MyArgs.UserInfo,
"ui":             MyArgs.UserInfo,
"server":         MyArgsServerInfo,
"si":             MyArgsServerInfo, }
`
In this example, both user and server commands can also be called using ui, si commands respectively.

#### Unnamed arguments

All arguments which are not flags or values of flags are classed as unnamed arguments or parameters.  
`... --no novalue unnamed1 -v unnamed2 unnamed3`  
In this example there are 3 unnamed values, (Assuming -v maps to a bool)

Once all flags and their values are removed from the command line, the remaining, unnamed arguments are used to map to
the parameters of the method being called, in the order they are presented.  
e.g. should a method have a signature such as:  
`func (ma Myargs) MyCommand(s string, t time.Time, count int)`  
It will require three unnamed parameters in the command line and those parameters must be valid for the data types in
that position.  
`mycommand "hello", "1/1/2001T12:00:00", 4`  
will parse correctly.  
`mycommand "1/1/2001T12:00:00", "hello", 4`  
will throw an error of invalid date.

Check the fields description for the data types supported as parameters.

  

