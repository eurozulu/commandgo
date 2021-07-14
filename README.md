# CommandGo
### Command line parser - object mapper

* Simplifies writing command line tools by mapping the command line arguments directly into functions and method.  
* Performs automatic type detection and parsing for all variables, fields and parameters with extendable framework for customised data types.
* Keeps all flags and commands in one place.  Simple to maintain and IDE friendly.
* Clean and simple to use.  No interfaces, struct's or types required to learn, just a simple map.
  
Unlike many command line parsers, Commandgo does not rely on any interface or predefined structures.  
It will map into almost any function/method or variable/field already in your application, so you don't have to maintain
additional code for your command line interface.  This works well for existing libraries and tools that require a simple
command line interface applied to it.  Of course you can create one of your own, using generic struct's, functions and variable
and map into these, but the choice is yours.

##### Goal
The goal of this parser is to simplfy the boiler work of mapping command line arguments into your application.  
The prime goals are:  
  
+ Simplicity.  Interface must be clean, simple and obvious to anyone looking at it.
No additional types or interfaces to be imposed on the developer.  They define the code, the framework maps into it.  
  
+ Reduce code maintenance
Keep code relating to the command line to a minimum. Keep the framework in the background.
Single point of maintenance.  All commands and flags in one, easily editable place, with direct links (via IDE) to the point thay map to.
Provide single point of help, easily maintained. (No double changes to text prompts and comments, use the comments!)  
  
Partly insprired by [Spring IOC]("https://www.baeldung.com/inversion-control-and-dependency-injection-in-spring") from
the Java world, structures and parameters can have their values 'injected' into them from an external
source, in this case, the command line arguments.
  
A large part of the design was to keep the API as simple as possible, so you, the developer, doesn't have a steep learning curve to use the framework.  
No complex structures or special interfaces to implement, Commandgo tries to be intuitive and obvious, so anyone should be able to look at your code
and understand exactly what its doing without having to know or learn the framework behind it.  
All the complexities of parsing value types and mapping functions is hidden away leaving the developer with a simple map.  
This maps the strings expected in the command line arguments, into their respective locations within your application,
those locations being functions, methods, variables or fields.  (known as mapped points)
  
### Usage
To show an example Commands map, first, here is a fictional "deployment" application struct, which performs deployments of code.  
This would be your own application code, just an example of something to map into.  
```

type EnvironmentType int
const (
	DEV EnvironmentType = iota
	TEST
	PROD
)

type Deployer struct {
    // Config file holds the data to connect to the jenkins servers etc.
    Config        io.ReaderCloser 
    Environment   EnvironmentType
}
// DoDeploy performs a deployment of the given project, to the configured servers and returns the version deployed 
func (d Deployer) DoDeploy(repoURL *url.URL) (float32, error) {
   ...
}
// DeployedVersion gets the version of the project at the given repo
func (d Deployer) DeployedVersion(repoURL *url.URL) (float32, error) {
   ...
}

```  
To map a command line interface into this application we might use a map something like this:
```
func main() {
    // create our 'app' with its default settings
    dp := &Deployer{Environment: DEV }
    
    cmds := commandgo.Commands{
        "deploy": dp.DoDeploy,
        "ver":    dp.DeployedVersion,
         
        "--config": &dp.Config,
        "--env":    &dp.Environment,
    }
    
    v, err := cmds.RunArgs()
    if err != nil {
        log.FatalLn(err)
    }
    
    fmt.Println(v)
}
```

The corresponding command line, for a binary named 'myapp', would be:  
  
``` 
myapp ver https://github.com/eurzulu/blabla -config ~/.deploy/config.yaml 
1.2

```  
``` 
myapp deploy https://github.com/eurzulu/blabla -config ~/.deploy/config.yaml -env 1
1.3

```

Both methods require a single URL parameter, returning a float, error so the framework parses the first arg after the 'command' as a url.  
The two flags are also both parsed into their respective types The 'EnvironmentType' parses as an int.
The Config field is defined as an interface `io.ReaderCloser`.  By default, the framework has a custom type of *File, which can handle this type.
The *File type parses an argument as a local file path and opens the corrisponding file for that path.  
  
And that's it!  You have a command line application.
  
This is the simplest example to show how it works.  
There's loads more features to help you maintain the most complex of command line parsing needs, from sub maps, default keys and functional flags to 
custom type mapping and automatic help generation.  Read on to discover these and see how commandgo can take the grind work out of writing command line tool
and let you focus on what the application does.  
  
  
  
### Commands
The Commands map is the heart of the framework.  It's a simple map using string keys which map to `interface{}` values.  
A very simple example:   
```
cmds := commandgo.Commands { 
      "mycommand": MyFunction,  
      "-myflag":   &SomeVariable,
    }
 ```
Where `MyFunction` is an exported function and `SomeVariable` is a variable.
The keys are the command line arguments and the `interface{}` value is either a func/method, or pointer to variable/field.    
  
Notice the function mapping requires no reference to the functions parameters, return values or their types, as with
the variable and its type.  These are all detected and handled for you by the framework.  
  
Should `MyFunction` look something like this:  
`func MyFunction(name string, id int, click *url.URL) (string, error) { ...}`    
Then the framework will expect three command line arguments following the "mycommand" argument, and those arguments must be parsable as:  
string, int and a url respectively.  If the wrong amount of arguments is given or they cannot be parsed into these types,
the framework reports the error back to the user, informing them of what is required.  
This frees the developer of the tedious task of parsing and basic validation.  
Mappings can be to almost any type of function/method or to any variable/struct field.  

The interface values can be one of three types:
+ func / method
+ variable/field pointer
+ another 'sub' Commands map

A helpful bonus of having the function / variables directly referenced in the map is that most IDE's allow you to click through to the implementation or definition of that item.  
Viewing the Commands map, you see the command being mapped, and a click takes you to the function or field it maps to.


#### Mapping points
Mapping points are the values, other than submaps, which keys are mapped to. i.e. func/methods, variable/fields.  
These points can be ANY exported point of a library or application or point local to the same package.  
For simple tools, your map can often call directly into the library or data model you wish to use,
without the need for an intermedeate structure or interface.  For more complex needs, an intermediate structure can be used, grouping your flags and commands into 
induvidual structs and performing 'pre-call' or 'post-call' operations, perhaps using flags set as fields etc.  
  

#### default key
The default key is an empty string, when present, gives the Commands a default location to call when no arg->key mapping is found.  
Without the default key, (the out of the box state) the framework will report a "missing command" or "unknown command" error when no key is found.  
When present, this key is used when no other matches.  
An important note about the default key is how parameters are handled.  Considure a cmd line of:  
`hello world`  and a map with no matching key but a default key.  The framework will first search for the `hello` key as the command,
and on failing, revert back to the default.  The `hello` will then revert to being a parameter (rather than a command).  
With no default key, it will report `hello` as an unknown command.
The same rules apply to the default mapping as any other, in that the arguments must match the mapping points type or signature.

#### Command alias
To specifiy more than one command name or flag, simply map two or more entiries with the same value.
```
cmd := commandgo.Commands{
  "mylongcommandname" : MyCommands.LongName,
  "mlcn"              : MyCommands.LongName,
}
```
The Help system detects these multi entries and groups all keys into the same help subject.

#### Submaps
In addition to func and vars etc, values may also be other Commands maps, containing their own set of flags and command keys.  
Using sub maps commands can be 'chained' into sequences, forming a hierarchy of commands.  
e.g. `get users accounts`, the `get` is in the 'top' map as a key to another map containing the `users`
key mapping to yet another map with the `accounts` key finally mapped to a function/method.  
Each level of commands has its own set of flags (and alternative commands), only being parsed when that sequence if followed.  
e.g. if we want a flag to apply to all `users` commands, we place them in the second map.  
Flags only applying to `account` command, are placed in the last map and so on.  Flags and commands are 'scoped' to the map they're defined in.  
When the final function map is found, all flags should have been assigned at their respective levels.  Remaining flags are considured an error.



An example of submaps, here, a mapping of a "getter" and "putter" structs being used to map commands to get and put resources
to either external urls, or local files.
```
Commmands{
  "-verbose": &Verbose,
  "-log":     &logs.Level,
  "about":    showAbout,
  "get": Commands{
      "":        getter.DoGet,
      "later":   getter.DelayedGet,
      "-format": &getter.OutputFormat,
  },
  "put": Commands{
      "":         putter.DoPut
      "-encrypt": &putter.Encrypt,
      "-key":     &putter.Key,
      "-user":    &putter.User,
      "new": Commands{ 
          "":        builder.Build,
          "-name":   &builder.Name, 
          "-id":     &builder.Id, 
          "-status": &builder.Status,
  },
}
```    
  
In this example there are five mappings in the 'root' map, two flags (-verbose, -log) and three commands, "get", "put" & "about".  
The top level flag mappings are called global flags, as they are always available to all commands.  These usually map to global variable.   
Regardless of the command being used, these flags will be parsed from the command line first.  
Of the three top level commands, two, `get` and `put` map into methods and `about` maps to a global function `showAbout()`.
The method mappings are using submaps to define some additional flags that are specific to those commands only.
In addition, `put` has a third level command `new`  command which maps yet another submap and to a Builder object for creating new instances.  
e.g.  
```myapp put http://myserver/theputtedstuff "This is the data" -encrypt -user john -key ~/.ssh/id_rsa.pub```  
or  
```mycmd put new -name mynewfile -id "blabla" -status draft```

An example of what this map is mapping into:  
```
type MyPutter struct {
  Encrypt bool
  Key *os.File
  User string
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

In the command line, Flags can appear in any order. All flags, with the exception of bool types must have a following
argument as its value.  
This value is converted to the relevant data type for the Field. Booleans MAY have a value, if it is parsable as a bool.
If they have a following argument which is not parsable as bool, that value is ignored by the bool flag. Bool flag are
True when they are present, unless they are followed by a 'false' value.



#### Flags
A Key may be marked as a 'Flag' by preceeding it with one or more '-' dash characters.  
Flags are usually optional arguments which can alter the behaviour of the 'main' command.  
Flag keys are treated with priority when executing the command line.  Non flags, which are not parameter values, are treated as
a command, and executed once.  Flags are ALL executed before the main command is invoked.
Any name can be mapped to any of these three mappings.

#### Parameters

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
  
#### Variadic Parameters
Variadic parameters are supported.  When present, the command line arguments
from the final position, onwards, are all parsed into a slice of the Variadic type.  
These allow the command line to accept zero or more arguments optionaly.  
Variadic parameters can be of any of the supported data types.  
A func with a signature such as:
```
func DoThis(s ...string) {}
```
Can accept any command line passed to it as all args are strings already, and all are optional.




  
### Data Types
When parsing the command line argument strings, the destination of the argument is examinied to determine its type.  
e.g. "-myflag": &SomeFloat  pointing to a float var attempts to parse the string following the -myflag agument, into a float.  
This is the same process for function parameters, struct fields and variables.
When mappng to flags, the order of the command line is unimportant, flags can be anywhere in the command line.  
When mapping to functions or methods, the order of the arguments and the count* must match that of the function parameters being called.  
*unless function/method uses variadic parameters, in which case argument count must only match the non variadic parameters. 

Most data types are supported with the exception of channels:  
+ int, int16, int32 int64, uint8, uint16, uint32, uint64
+ float32 float64 
+ bool
+ string    
+ slices 
+ maps
+ struct
  
`bool` types are exceptional as they are the only type not requiring a value.  They default to true when no value is provided.  
`structs` are parsed as json from the command line or, if the struct supports encoding, will be Unmarshalled using that.
+ Those supporting [encoding.BinaryUnmarshaler](https://golang.org/pkg/encoding/#BinaryUnmarshaler) interface
+ Those supporting [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler) interface

Additional types can be added via the 'CustomType' system.  
  
  
#### Custom Data type  
Data type support can be extended using custom data types to allow mapping into variables or functions which use a specific type.  
These types define a specific data type and provide a custom function to parse the string argument into that type.  
The framework include three custom types out of the box:  
+ \*url.Url
+ time.Time, time.Duration
+ \*os.File  
  
Custom types apply to both fields/variable values and func/method parameters.  
By specifying a custom type, function parameters and variables of any type which can be mapped directly from the command line and parsed in the required type.
The type only has to be "assignable", therefore a parameter type of say `io.Reader` will be parsed by the \*os.File custom type.  

To define a new custom type use the `values.NewCustomType` method, which accepts a reflect.Type and a ArgValue function.  
The ArgValue function is passed a string argument and a reflect.Type of the type required.  
It returns an `interface{}, error` of the value it parses the argument into.  
  
As an example of a custom type, take the first example "deployment" app above, the EnvironmentType is an `int` const.  
Rather than accept numbers on the command line, we could map this from a more natural string name:  
```
ett := reflect.TypeOf(EnvironmentType(0))
values.NewCustomType(ett, func(s string, t reflect.Type) (interface{}, error) {
	switch s {
	case "DEV":
		return DEV, nil
	case "TEST":
		return TEST, nil
	case "PROD":
		return PROD, nil
	default:
	    return 0, fmt.Errorf("%s is not a known environment", s)
	}
})
```  
Then the flag would be more natural `-env DEV` rather than `-env 0`  
Once registered this way, the type will automatically be called on all values which are assignable to `EnvironmentType`.  



### Help System
All command line parsers require a help system to guide the final user about the commands and flags.  
The help system is current designed to be open and freely editable by the application designer.
Custom help can be added or replace any existing help.  

In line with minimal effort, the help system aims to use Godoc comments to form the help system.  
This is still currently under development, but is aimed as a pre-build process, extracting the key mappings from source
and matching them to the comments they map to.
