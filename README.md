Mainline
--------
**Overview**  
A command line argument parser with a twist.  
The goal of mainline is to simplify creating tools which use a command line interface by getting rid of all the boiler plate code usually associated with libraires of this kind.  
Maintaining a layer of 'Command' objects, with their help text and flag config to map into functions in your application model adds a layer of complexity and work which can be mostly automated.  
Mainline aims to automate the boiler plate mapping of command line text, into function calls in your applications.  

Partly inspired by Java Spring's, container injection, 
Mainline takes the concept of a 'bean', translated to a generic go struct, and applied a 'bean definition', the arguments on the command line,
injecting values and parameters into it, from the command line.  
  
Methods on the struct are mapped to string commands and the remaining command line arguments, following that command, will be parsed into the correct data types for the method parameters it is mapped to.  

Field values in the struct automatically become argument 'flags', (arguments starting with '-' or '--') 
which, like the method parameters, can also be mapped into most data types.  
An optional tag is available on the fields to mark them with alternative names or to ignore them.


**Usage**  
Using a generic (and somewhat contrived) 'Client' type from our application model as our 'bean':   
`type Client struct {`  
`}`  
`func (c Client) ReadMessages(u *url.Url) []string { ... }`  
`func (c Client) SendMessages(u *url.Url, msgs []string, sendBy time.Time) { ... }`  

- Note: This can be ANY struct.  One in your code or even a 3rd party lib struct.  
There is no constraint or interface required. This is simply an example.  

To use that struct we map any of its exported methods to string commands:  

`mainline.AddCommand("read", Client.ReadMessages)`  
`mainline.AddCommand("send", Client.SendMessages)`  

And then call Mainline from `main()` function to kick it off:  
`mainline.RunCommandLine()`  

And that's it!  
The application (myclient) can be started:  
`>myclient read http://www.spoofer.org/messages`  
or  
`>myclient send http://www.spoofer.org/messages, "'We meet at one', 'one is the time'" 13:00:00`  

And the `ReadMessages` or `SendMessages` methods will be called with the correct parameter types for each method.    
Any number of commands can be mapped in this way, to any number of methods and structs or functions.  There are few limits on the functions and methods which can be mapped into, mostly limited by the data types of the parameters.  See further on for data type limitations.  

*Fields*  
To add flags to the command line, simply add a field to your struct:  

`type Client struct {`  
>`Debug bool`  

`}`  
`func (c Client) ReadMessages(u *url.Url) []string { ... }`  
`func (c Client) SendMessages(u *url.Url, msgs []string, sendBy time.Time) { ... }`  


By default a field will use its name as the flag name so the application could be started with:  
`>myclient --debug read http://www.spoofer.org/`  
And the field value of `Debug` will be set to true.

Flag position is unimportant in the command line, however NON flags, parameter position is important as these must align with the parameters on the method being called.  
`>myclient read http://www.spoofer.org/ --debug`  
is the same as the above.  

*Tags*  
Fields can be marked with a tag to give them names, other than their field name, as well as alternative names.    
``Debug bool ` `` ``flag:"debug, d"` `` 

This marks the field with two names, 'debug' and 'd', both `--debug` and `--d` will set the Debug field.  

A tag may also indicate when a field should be ignored and not used as a flag.  
Marking a field name with a dash will mark the field as a non flag e.g. 'flag:"-"'  


*Data types*  
There is no limit to the number of parameters on a method however for every parameter there must be a corrisponding argument given on the command line.  
Slices are supported vario/optional parameters will be supported in the future.  
Most data types can be mapped to, from sensible values, even some structs.
`url.URL` is supported as seen in the example, as is `time.Time`, but also any struct 
which can parse json or encoded text, which they decode from the argument string value.  
e.g. The json argument can be passed as a string json form.  
  
Obviously the base types, string, int, bool and float as well as slices.  
Not sure about maps yet, as they could be just json objects. 


Part of the 'magic' is the mapping of data types from the string arguments, into the relevant type for the struct it maps to.  Using go's `reflect` package, Mainline matches the signatures of each field or method with the arguments available and attempts to convert the string arguments into data type which match those signatures.
This allows the developer to define the exact parameter types for each command and the types of each flag for that command.


Most data types are supported with some limitations:  
+ `struct` must implement the `encoding.Unmarshaller` or `json.Unmarshaller` interface
with the exception of URL's, which are supported also.
+ `channel`, and `func` types are NOT supported
+ `slice` is supported, with the same limitations applied to its element type.
+ `map` is not yet supported (As I think json encoder will do the same thing?)
+ all the base types, float, int, bool and string are supported.

The result is a clean, simple interface, which offers a powerfull way to build simple, command line tools.


**Help**  
No command line tool would be complete without a help system to describe the commands and their parameters and flags.  
Following the concept of no boiler plate code, mainline generates the help system for you by using the comments of the functions it is calling into as the help text for the command mapped to that function.  
This removes the need to maintain a second set of documentation for a command set and utilises existing resources.  
Mainline uses `go doc` to locate and read the comments of the functions and builds a help system from that text.  
To use the help, the developer runs a build tool called `helpmaker` in the directory of their source code.  helpmaker will generate a new function named `HelpCommand` which the developer can then map into with:  
`mainline.AddCommand("help", _help.HelpCommand)`



**Notes**  
The first argument dictates the command name  

All arguments marked as flags (- or --) are stripped from the argument list, along with their values.  
Bool flags only take a value if the following argument can be parsed as bool.  Otherwise they default to true.  
The remaining, unnamed, arguments are used as parameters to the method call.  
These remaining arguments are matched to the available parameters in the method by position, so first remaining arg is first parameter and so on.  
If the number of arguments and parameters do not match, an error is thrown as is one thrown if the string argument can not be parsed into that particular data type.  


TODO:  
Add control of output.  possibly use 'hidden' commands, commands created by the library, starting with '_', which don't show up in help.
These can control aspects of the library, without interfering with the code using it.  

Create a Token argument to specify new object, with no argument.  
e.g. --dateformat .  
So the type of the field Dateformat, dictates the actual instance.  
`Dateformat time.RFC3339`  

Also provide tokens for std.out and in.

Add sub commands.  
mapping commands in a path form:  

`mainline.AddCommand("basecmd", blabla.Bla)`  
`mainline.AddCommand("basecmd/subone", blabla.Bling)`  
`mainline.AddCommand("basecmd/subtwo", rahrah.Rah)`  
`mainline.AddCommand("basecmd/subtwo/grandsubone", rahrah.Bla)`  
