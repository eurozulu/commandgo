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

To use a function, simply write the parameters you need and the framework ensures the command line arguments given are parsed into
the correct data types for those parameters.  
For optional, command line flags (-myflag) values can be mapped to fields in a specific structure or to
global variables.

Offer structure to the functional aspect of the command line tool. Each command can be global, or 'wrapped' in its own
specific strucutre. Flags can be global or specific to those methods in the structure.




### Commands map
Commands map is the core aspect of the framework.  It maps the command line arguments to variable, fields, fucntions and method.  
The Command map consists of at least one 'Commands' map but can also have 'sub maps', mapping to additional command maps,
to form a hierachey of commands.  Each map contains its own command set and flags for that particualr mapping.  
This structure allows for a flexible command structure to be formed giving a logical context to each group of commands.  
Commands which share flags can be groupped into one command map, isolating those flags from other commands.  
These often map logically to command structs, each struct having a field set for flags and a respective comamnd map
mapping its fields and methods to command line commands.
  
e.g.  
```
Commmands{
  "-flag1" : &Flag1,
  "-flag2" : &myapp.Name,
  "about" : showAbout,
  "get" : Commands {
      "" :&getter.DoGet,
      "later" :&getter.DelayedGet,
      "-format": &getter.OutputFormat,
  },
  "put" : Commands {
      "" :&putter.DoPut
      "encrypt" :&putter.Encrypt,
      "-key": &putter.key,
      "-user": &putter.user,
      "new" : Commands{ 
          "" : &builder.New,
          "-name": &builder.Name, 
          "-id": &builder.Id, 
          "-status": &builder.Status,
  },
}
```  
  
In this example there are four mappings in the 'root' map, two flags and two commands.  
The flag mappings, ` -flag1` and `flag2` are called global flags, as they are always available.  
Regardless of the command being used, these flags will be parsed from the command line first.  
This leaves the two commands `get` and `put`.  These both map to their own sub maps, both having a unique set of flags.  
get uses the default, "", mapping to map to the 'DoGet' method on the 'getter' object.  

### Execution order
On calling `Run` or `RunArgs` the command line is parsed in the following order:  
- The Flags are located along with their following values.  
- Flags are each run, first the assignment mappings, followed by any func mappings.  
- Finally the command mapping is found and run, using any remaining args (not consumed by flags) as parameters (or values) for the command.
  
The command line is parsed by passing it to each stage, which 'consumes' arguments from it.  Consume meaning they are no longer
made available to the following func or assignments.  Think of each mapping taking its arguments from the command line until all thats left
is a command followed by its parameters.  Each flag takes what it needs and finally the command uses the remaining to set its parameters.
  


### Mappings
Mappings map the command line command and flag names to their respective points in the application or additional mappings.
Mappings can be seen as two types, Commands and Flags.  The prime distinction between these is during execution,
A single command is executed whereas ALL flags from the command line are executed, prior to the command execution.  
Usually flags are mapped to variables or fields as a means of assigning a 'setting' the method/func call will use.  
Again, usually, commands map to a func or method which is called and the result being the final output.  
However this is only a concept which clarifies the usualy behaviour of a cmdline execution.  
Internally mappings are viewed as assignments and calls.  Assignments being pointers to fields/variable and
calls being methods/func.  Both flags and commands can map to either.  
In addtion, mappings can also map to a sub command map, containing its own set of commands and flags.  
  
 
The mapping consist of a unique string key, mapped to either:  
- A pointer to a variable or field
- A function or method
- A 'sub map' of additional mappings.  
  
The key in the mapping should be any UTF8 string which could be reasonable input from the command line, with the exception of whitespace.  
No whitespace is allowed in keys.  

#### Default Key
The map map contain a single empty key which is treated as the 'defualt' mapping for the map.  
Default mapping is invoked when no command is found in the command line, after all the falgs and values have been removed.  
  

#### Flags
A Key may be marked as a 'Flag' by preceeding it with one or more '-' dash characters.  
Flags are usually optional arguments which can alter the behaviour of the 'main' command.  
Flag keys are treated with priority when executing the command line.  Non flags, which are not parameter values, are treated as
a command, and executed once.  Flags are ALL executed before the main command is invoked.
Any name can be mapped to any of these three mappings.


Note on flags mapped to functions.  
Flags are usually optional 'settings' which naturally map to variable/field pointers, however they can be mapped to functions, which
are executed prior to the 'main' command.  (The help system uses such a mapping for the -? and --help flags, to execute the help function.)  
When mapping flags to functions care should be taken on how that flag will 'consume' the arguments given on the command line.  
With no func flags (or fixed param func flags), the command line will be interpretted in any given order.  That is to say, flags 
and commands can be placed in any order or even mixed (e.g. -flag1 hello mycommand "do this" -verbose true) still parses as:
mycommand "do this" -flag1 hello -verbose true  
  
This behaviour will be changed if flags are mapped to functions using variadic parameters.
Functions with variadic (optional) parameters are offered all of the available, non flag args following the mapped argument name.
whereas non variadic functions are only offered the number of args they require, leaving any following those args, in the command line as commands.
This effects how the command line is interpreted at runtime, therefore care should be taken when choosing to map flags to functions which
use variatic paramters.  Avoiding these allows the command line to be interpreted more freely, allowing any order of flags and comamnds to be parsed correctly.  
e.g. if the above example, -flag1 was variadic, it would be parsed as:
-flag1 hello mycommand "do this"
-verbose true
and no command! where as placing -flag1 on the end of the command line, it would parse as expected.
Unless absoluetly required, avoid mapping flags to func/methods using variadic parameters.  
This doesn't apply to commands as they are exeuted last and therefore will have all of the remaining command line, not already consumed,
and with all flags (and flag values) removed. (Any remaining, unconsumed flags prior to command execution throw an error of unknown flag)


Flags mapping to function with non variadic (fixed) parameters can be placed in front of commands, whereas those mapping to
variadic parameter functions will 'consume' the following command as a parameter.
e.g. from the example above, if -flag2 has a non variadic, single bool parameter, it could be safely placed in front
of the command.
-flag2 false cmd dothisthing -flag1 hello world
(With bool only, even the 'false' could be ommited safely)
however, if, for example -flag1 was variadic, placing it in front of the command will prevent the 'cmd' mapping ever being executed,
as it will be 'consumed' by the -flag1 function.
-flag2 false -flag1 hello world cmd dothisthing
will call -flag1's func with 3





========================== old 
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


