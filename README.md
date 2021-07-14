# CommandGo

### Command line parser / object mapper

* Simplifies writing command line tools by mapping command line arguments directly into functions and method.  
* Keeps all flags and commands in one place.
* Maps command line either  into your own "command" structures or directly into your application model.
* Performs automatic type detection and parsing for all variables, field and parameters with extendable framework for customised data type.  
(Supports optional, varidac parameters for any of the supported data types.)  
  

##### Goal
The goal of this parser is to simplfy the boiler work of mapping command line arguments into commands, parameters and
flags.  
Partly insprired by [Spring IOC]("https://www.baeldung.com/inversion-control-and-dependency-injection-in-spring") from
the Java world, generic variables, structures and parameters can have their values 'injected' into them from an external
source, in this case, the command line arguments.

To use a function, simply write the function with the parameters you need and the framework will parse the given command line arguments
into the correct types for your function.
Global variables and struct fields can be mapped to flags, and the framework parses the flag values into the required type.

### Commands map
Commands map is the core aspect of the framework.  It maps elements from the command line (arguments) to variable, fields, functions and methods or sub commands.  
The Command map consists of at least one map but can also have 'sub maps', mapping to additional sub commands and flag, to form a hierachey of commands.  
Each map contains its own command set and flags for that particualr mapping.  
This structure allows for a flexible command structure to be formed giving a logical context to each group of commands.  
Commands which share flags can be groupped into one command map, isolating those flags from other commands.  
These often map logically to command structs, each struct having a field set for flags and a respective comamnd map
mapping its fields and methods to command line commands.
  
e.g.  
```
Commmands{
  "-verbose": &Verbose,
  "-log":     &logs.Level,
  "about":    showAbout,
  "get": Commands {
      "":        getter.DoGet,
      "later":   getter.DelayedGet,
      "-format": &getter.OutputFormat,
  },
  "put" : Commands {
      "" :       putter.DoPut
      "encrypt" :&putter.Encrypt,
      "-key":    &putter.key,
      "-user":   &putter.user,
      "new" : Commands{ 
          "":      builder.New,
          "-name": &builder.Name, 
          "-id":   &builder.Id, 
          "-status": &builder.Status,
  },
}
```  
  
In this example there are five mappings in the 'root' map, two flags (-verbose, -log) and three commands, "get", "put" & "about".  
The top level flag mappings are called global flags, as they are always available to all commands.  These usually map to global variable.   
Regardless of the command being used, these flags will be parsed from the command line first.  
Of the three top level commands, two, `get` and `put` map into methods and `about` maps to a global function `showAbout()`.
The method mappings are using submaps to define some additional flags that are specific to those commands only.
In addition, put has a third level command `new` which maps into a Builder object for creating new instances.  
e.g.
```mycmd put http://myserver/theputtedstuff "This is the data" -encrypt -user john -key ~/.ssh/id_rsa.pub```
or
```mycmd put new -name mynewfile -id "blabla" -status draft```


An example of what this map is mapping into:  
```
type MyGetter struct {
  Outputformat string
}

func (mg MyGetter) DoGet(url *url.Url) (string, error) {
  ...
}

func (mg MyGetter) DelayedGet(url *url.Url) (string, error) {
  ...
}
```
Note the parameters of these function are the types required (URL), which is parsed for you by the framework, from the given command line.  
Command lines missing arguments or using malformed arguments for that type are reported as errors automatically by the framework.  


### Execution
Once a map is created, it can be called using the arguments to be parsed.  
Commands has two points to call, `Run(args ...string)` and a convienience method `runArgs()` which simply uses the os.Args.  


#### Execution order
On calling `Run` or `RunArgs` the command line is parsed in the following order:  
- The Flags are located along with their following values.  
- Flags are each run, first the assignment mappings, followed by any func mappings.  
- Finally the command mapping is found and run, using any remaining args (not consumed by flags) as parameters (or values) for the command.
  
The command line is parsed by passing it to each map and sub map, which 'consumes' arguments from it.  Consume meaning they are no longer
made available to the following func or assignments.  Think of each mapping taking its arguments from the command line until all thats left
is a command followed by its parameters.  Each flag takes what it needs and finally the command uses the remaining to set its parameters.
  


### Mappings
Mappings map the command line argument to their respective points in the application or additional mappings.
Mappings can be seen as two types, Commands and Flags.  The prime distinction between these is during execution,
A single command is executed whereas ALL flags from the command line are executed, prior to the command execution.  
Usually flags are mapped to variables or fields as a means of assigning a 'setting' the method/func call will use.  
Usually, commands map to a func or method which is called and the result being the final output.  
However this is only a concept which clarifies the usualy behaviour of a cmdline execution.  
Internally mappings are viewed as assignments and calls.  Assignments being pointers to fields/variable and
calls being methods/func.  Both flags and commands can map to either.  
In addtion, mappings can also map to a sub command map, containing its own set of commands and flags.  
Sub maps are passed the remaining comand line arguments, after the parent map flags have consumed its flags.  
The process repeats on the sub map until it reaches a non sub map mapping (func or var)
  
 
The mapping consist of a unique string key, mapped to either:  
- A pointer to a variable or field
- A function or method
- A 'sub map' of additional mappings.  
  
The key in the mapping should be any UTF8 string which could be reasonable input from the command line, with the exception of whitespace.  
No whitespace is allowed in keys.  

#### Default Key
The map may contain a single empty key which is treated as the 'defualt' mapping for the map.  
Default mapping is invoked when no command is found in the command line, after all the flags and their values have been removed.  
  

#### Flags
A Key may be marked as a 'Flag' by preceeding it with one or more '-' dash characters.  
Flags are usually optional arguments which can alter the behaviour of the 'main' command.  
Flag keys are treated with priority when executing the command line.  Non flags, which are not parameter values, are treated as
a command, and executed once.  Flags are ALL executed before the main command is invoked.
Any name can be mapped to any of these three mappings.

# Data Types
When parsing the command line argument strings, the destination of the argument is examinied to determine its type.  
e.g. "-myflag": &SomeFloat  pointing to a float var attempts to parse the string following the -myflag agument, into a float.  
This is the same process for function parameters.  Remaining arguments are converted into the type according to the order they appear
and the function parameters signature.

Most data types are supported, all the base types, int64, float32/64, bool, string etc, as well as    
Slices, Maps, URL, Time and some other structs.

In the command line, Flags can appear in any order. All flags, with the exception of bool types must have a following
argument as its value.  
This value is converted to the relevant data type for the Field. Booleans MAY have a value, if it is parsable as a bool.
If they have a following argument which is not parsable as bool, that value is ignored by the bool flag. Bool flag are
True when they are present, unless they are followed by a 'false' value.

certain structs are supported:

+ Those supporting [encoding.BinaryUnmarshaler](https://golang.org/pkg/encoding/#BinaryUnmarshaler) interface
+ Those supporting [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler) interface  

### Custom Data type  
Data types support can be extended using custom data types.  These types define a specific data type and provide a custom function to parse the string argument into that type.  
The framework include three custom types out of the box:  
+ *url.Url
+ Time, Duration
+ *os.File  
  
Custom types apply to both fields/variable values and func/method parameters.  
By specifying a custom type, function parameters and variables of any type can be mapped directly from the command line and parsed in the required type.  
  
CustomTypes uses the `NewCustomType` method, which accepts a reflect.Type and a ArgValue function.


  
  

  
Flags may be mapped to global variables using a pointer to that variable and assigning one or more flag names to it:  
`commandgo.AddFlag(&Verbose, "verbose","v")`
This assumes there is a global variable called Verbose:
`var Verbose bool`  

Being global, all command have access to these values.  
Fields may also be mapped to struct fields.  For commands wishing to have flags specific to that command,
and not global.  Such commands can map to a method in a generic struct, which makes all the fields in 
that struct available as flags.

#### Command alias

To specifiy more than one command name or flag, simply map two or more entiries with the same value.
```
cmd := commandgo.Commands{
  "mylongcommandname" : MyCommands.LongName,
  "mlcn"              : MyCommands.LongName,
}
```
The Help system detects these multi entries and groups all keys into the same help subject.  


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
These allow the command line to accept zero or more arguments optionaly.  
A func with a signature such as:  
```
func DoThis(s ...string) {}
```
Can accept any command line passed to it as all args are strings already, any all are optional.  


### Help System
All command line parsers require a help system to guide the final user about the commands and flags.  
The help system is current designed to be open and freely editable by the application designer.
Custom help can be added or replace any existing help.  

In line with minimal effort, the help system aims to use Godoc comments to form the help system.  
This is still currently under development, but is aimed as a pre-build process, extracting the key mappings from source
and matching them to the comments they map to.
