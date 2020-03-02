# Mainline
### Command line arguments object mapper  
  
Emulates the decoder package by 'decoding' / Unmarshalling the cmd line arguments into a generic struct.  
Maps named flags into structure fields, coercing the string arguments into their respective data types using go reflection.  
    
 
##### Usage
To parse an example command line:  
`... --name john -timeout 24h -d param1 param2 param3`

Create a struct containing fields for the named argument flags required.
e.g. We have three named flags:    
* name  string
* timeout time.Duration
* debug  bool  
  
and an expected unnamed parameters 'param1', 'param2', 'param3'  
  
We could create a struct:
```
type MyArgs struct {  
    Name    string  
    Timeout time.Duration   `flag:"t"`
    Debug   bool            `flag:"d"`
    Params  []string        `flag:"*"`
}
```
Once decoded, the MyArgs instance will have the arguments mapped
to the fields, with the Params slice containing the 'param1', 'param2'...
  
To decode the struct, in the application `main()`:  
```
func main() {
    var args MyArgs
    err := mainline.NewDecoder(os.Args[1:]).Decode(&args)
    if err != nil {
    	t.Fatal(err)
    }
    
    // That's it!, args is ready to use... e.g.
    ctx, cnl := context.WithTimeout(context.Background(), args.Timeout)
    ...
    if args.Debug {
        log.Loglevel(log.Debug)
    }
    fmt.Printf("Hello %s", args.Name)
}   
```
  
Flags can appear in any order.  All flags, with the exception of bool types must have a following argument as its value.  
This value is converted to the relevant data type for the Field.
Booleans MAY have a value, if it is parsable as a bool.  If they have a following argument which is not parsable as bool, that value is ignored by the bool flag.
Bool flag are True when they are present, unless they are followed by a 'false' value.  

Most other data types are supported, all the base types, int64, float32/64, bool, string etc.  
Slices are supported. Maps not yet.  
certain structs are supported:
+ Those implementing the [json.UnmarshalJSON](https://golang.org/pkg/encoding/json/#example__customMarshalJSON) interface
+ Those supporting [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler) interface
+ Date, Duration and url.URL


####Tags
Fields may be tagged to specify alternative names for the flag using standard go tagging.
e.g.  
```
type MyAppConfig struct {  
   Name     string           `flag:"nom", "naam", "n"`  
   Timeout  time.Duration    `flag:"t"`
   Debug    bool             `flag:"db"`
}
``` 
Using these flags, the `Name` field could be set with any of the following command line argiments:  
+ --name john
+ -n john
+ -nom "alice cooper"

Tagging a field with a '-' `flag:"-"` will hide that field from the argument parsing.

Command func(args ...string)  `command:"dothis"`

There is no distiction between the double dash and single dash for flags.  "-" is the same as "--"  
 
####Unnamed arguments
All arguments which are not flags or values of flags are classed as unnamed arguments or parameters.  
`... --no novalue unnamed1 -v unnamed2 unnamed3`  
In this example there are 3 unnamed values, (Assuming -v maps to a bool)  
  
Once all flags and their values are removed from the command line, the remaining, unnamed arguments may be mapped to a field
using the flag tag with a "*" name.  e.g. to map the unnamed to MyCmdArgs:  
``` MyCmdArgs []string  `flag:"*"` ```  
In this example, `MyCmdArgs` would contain the three values:  unnamed1, unnamed2, unnamed3  
  
  

