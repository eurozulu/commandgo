Commando
--------
**Overview**  
A command line argument parser with a twist.  
Partly inspired by Java Spring's, container injection, 
Commando takes the concept of a 'bean' and translates it into a generic go struct.
Using the command line arguments as the 'bean definition' it creates an instance of the strut and injects the values from the command line into it.  
  
Methods on the struct can be mapped to string commands and the remaining command line arguments, following that command, will be parsed into the correct data types for the method parameters it is mapped to.  

Field values in the struct automatically become argument 'flags', (arguments starting with '-' or '--') 
which, like the method parameters, can also be mapped into most data types.  
An optional tag is available on the fields to mark them with alternative names or to ignore them.


**Usage**  
Using a generic 'MyConfig' type from our application model as our 'bean':   
`type MyConfig struct {`  
`}`  
`func (c MyConfig) SetHost(u *url.Url) { ... }`  
`func (c MyConfig) SetDatabase(db string, collection string, readOnly bool) { ... }`  

- Note: This can be ANY struct.  One in your code or even a 3rd party lib struct.  
There is no constraint or interface required. This is simply an example.  

To use that struct we map any one of its exported methods to a string command, say 'host':  

`commando.AddCommand("host", Config.SetHost)`  

And then call commando from `main` to kick it off:  
`commando.RunCommandLine()`  

And that's it!  
The application (mycode) can be started:  
`>mycode host http://www.spoofer.org/`
And the `SetHost` method will be called with a URL pointer containg the given url, as a parameter.  

*Fields*  
To add flags to the command line, simply add a field to your struct:  
`type MyConfig struct {`  
>>`Debug bool`  

`}`  
`func (c MyConfig) SetHost(u *url.Url) { ... }`  
`func (c MyConfig) SetDatabase(db string, collection string, readOnly bool) { ... }`  

By default a field will use its name as the flag name so the application could be started with:  
`>mycode --debug host http://www.spoofer.org/`  
And the field value of `Debug` will be set to true.

Flag position is unimportant in the command line, however NON flags, parameter position is important as these must align with the parameters on the method being called.  
`>mycode host http://www.spoofer.org/ --debug`  
is the same as the above.  

Fields can be marked with a tag to give them names, other than their field name, as well as alternative names.    
>> `Debug bool 'flag:"debug, d"'`  
(Note incorrect outer quotes, but can't get tag /raw string quotes into mark down /-)  

This marks the field with two name, 'debug' and 'db', both `--debug` and `--d` will set the Debug field.  

There is no limit to the number of parameters on a method however for every parameter there must be a corrisponding argument given on the command line.  
Slices are supported therefore go, optional paramters are supported as a result.  
Most data types can be mapped to from sensible values, even some structs.
URLs are supported as seen in the example, but also any struct 
which can parse json or encoded text, which they decode form the argument.  
Obviously the base types, string, int, bool and float as well as slices.  
Not sure about maps yet. 


Part of the 'magic' is the mapping of data types from the string arguments, into the relevant type for the struct it maps to.  Using go's `reflect` package, commando matches the signatures of each field or method with the arguments available and attempts to convert the string arguments into data type which match those signatures.
This allows the developer to define the exact parameter types for each command and the types of each flag for that command.


Most data types are supported with some limitations:  
+ `struct` must implement the `encoding.Marshaller` or `json.Marshaller` interface
with the exception of URL's, which are supported also.
+ `channel`, and `func` types are NOT supported
+ `slice` is supported
+ `map` is not yet supported (As I think json encoder will do the same thing?)
+ all the base types, float, int, bool and string are supported.

The result is a clean, simple interface, which offers a powerfull way to build simple, command line tools.





**Notes**  
The first argument dictates the command name  

All arguments marked as flags (- or --) are stripped from the argument list, along with their values.  
Bool flags only take a value if the following argument can be parsed as bool.  Otherwise they default to true.  
The remaining, unnamed, arguments are used as parameters to the method call.  These remaining arguments are matched to the available
methods in the structure to find the one with the signature whic matches the available arguments.  
If no methods match, an error is thrown.  If more than one method matches, the first method in the struct is called.  
  
The order of flags, within the command line, is unimportant.  During parsing all flags and data are removed, leaving just the unnamed parameters which are passed to the method.  
The order of the remaining parameters IS important, and is used to map to the correct method.  
i.e. Method signatures (i int, s string) and (s string, i int) are two distinct method signatures.  


TODO:  
Integrate Godoc generation on the command struct into a help system for the command line.  
Use commants from each method and from the struct itself to generate output for --help command
