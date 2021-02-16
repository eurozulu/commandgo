# CommandGo

### Command line arguments object mapper

Simplifies writing command line tools by mapping command line arguments into functions and method. Maps string "
commands" into methods and/or functions, parsing the remaining arguments into the parameters for that method/func call.
Simply define the method and the parameters it requires, and arguments will be mapped into those parameters.

##### Goal

The goal of this parser is to simplfy the boiler work of mapping command line arguments into commands, parameters and
flags.  
Partly insprired by [Spring IOC]("https://www.baeldung.com/inversion-control-and-dependency-injection-in-spring") from
the Java world, generic variables, structures and parameters can have their values 'injected' into them from an external
source, in this case, the command line arguments.

To use a function, simply write the paramters you need and the framework ensures the paramters given can be parsed into
those values. For optional, command line flags (-myflag) values can be mapped to fields in a specific structure or to
global variables.

Offer structure to the functional aspect of the command line tool. Each command can be global, or 'wrapped' in its own
specific strucutre. Flags can be global or specific to those methods in the structure.

##### Usage

A simple, two command "tool" which prints out user "data" or server data.

```
func main() {
    cmds := commandgo.Commands{
		"user":       UserInfo,
		"server":     ServerInfo,
		"":           ShowHelp
	}

	err := cmds.Run(os.Args...)
	if err != nil {
		fmt.Println(err)
	}
}

func UserInfo(name string) error {
  fmt.Printf("Hello %s\n", name)
}

func ServerInfo(url *url.Url) error {
  r, err := http.Get(url.String())
  if err != nil {
      return err
  }
  
  fmt.Printf("%v\n", r.Header)
}
func ShowHelp() error {
  fmt.Println("use 'user' or 'server' to output data")
}

```

### Flags

Flags are any command line argument beginning with a dash '-'. Any number of dashes are treated as a single dash.  
Following the dash is the flag name. Each flag is a named value, mapped to a variable or field in the applications.  
The value of the flag is dependant on the type of the variable or field its mapped to. For most flags, the argument
following the named flag is taken as the value.  
e.g. `myapp mycmd -format json`  looks for a field or variable mapped to 'format' and assigns the 'json' value to it.

Most data types are supported, all the base types, int64, float32/64, bool, string etc, as well as    
Slices, Maps, URL, Time and some other structs.

In the command line, Flags can appear in any order. All flags, with the exception of bool types must have a following
argument as its value.  
This value is converted to the relevant data type for the Field. Booleans MAY have a value, if it is parsable as a bool.
If they have a following argument which is not parsable as bool, that value is ignored by the bool flag. Bool flag are
True when they are present, unless they are followed by a 'false' value.

certain structs are supported:

+ Those implementing the [json.UnmarshalJSON](https://golang.org/pkg/encoding/json/#example__customMarshalJSON)
  interface
+ Those supporting [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler) interface
+ Date, Duration and url.URL
  
Flags may be mapped to global variables using a pointer to that variable and assigning one or more flag names to it:  
`commandgo.AddFlag(&Verbose, "verbose","v")`
This assumes there is a global variable called Verbose:
`var Verbose bool`  

Being global, all command have access to these values.  
Fields may also be mapped to struct fields.  For commands wishing to have flags specific to that command,
and not global.  Such commands can map to a method in a generic struct, which makes all the fields in 
that struct available as flags.

#### Tags

Fields may be tagged to specify alternative names for the flag using standard go tagging. e.g.

```
type MyAppConfig struct {  
   Name     string           `flag:"nom, naam, n"`  
   Timeout  time.Duration    `flag:"timeout,t"`
   Debug    bool             `flag:"debug,db"`
}
```  

If tags are specified, only the names in the tag will map to flags. If no tags are specified, the fieldname, in lower
case, is used.

Tagging a field with a '-' `flag:"-"` will hide that field from the argument parsing.

There is no distiction between the double dash and single dash for flags.  "-" is the same as "--"

#### Command alias

To specifiy more than one command name for the same function, simply map two or more entiries with the same value.

```
cmd := commandgo.Commands{
  "mylongcommandname" : MyCommands.LongName,
  "mlcn"              : MyCommands.LongName,
}
```

#### Unnamed arguments / Parameters

All arguments which are not flags or values of flags are classed as unnamed arguments or parameters.  
`... --myflag myflagvalue unnamed1 -v unnamed2 unnamed3`  
In this example there are 3 unnamed values, (Assuming -v maps to a bool)

Once all flags and their values are removed from the command line, the remaining, unnamed arguments are used to map to
the parameters of the method being called, in the order they are presented.  
e.g. should a method or function have a signature such as:  
`func (ma Myargs) MyCommand(s string, t time.Time, count int)`  
It will require three unnamed parameters in the command line and those parameters must be valid for the data types in
that position.  
`mycommand hello "1/1/2001T12:00:00" 4`  
will parse correctly.  
`mycommand "1/1/2001T12:00:00" hello 4`  
will throw an error of invalid date.

Check the fields description for the data types supported as parameters.

#### Variadic Parameters
Variadic parameters are supported.  When present, the command line arguments 
from the final position, onwards, are all parsed into a slice of the Variadic type.


