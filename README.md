# Mainline
### Command line arguments object mapper  
  
  
Library for developing command line based tools requiring the parsing of command line flags and command.  
  
Maps named flags into structure fields, coercing the string arguments into their respective data types in the fields using go reflection.
  
  
 
#### Usage
Create a struct which contains all the argument flags you require for your application.
This struct acts like a configuration object, containing the properties as they are set by the command line.  
e.g. If we have three flags:  
- name  A string name
- timeout  A duration
- debug  A bool flag  
  
We could create a struct:
```
type MyAppConfig struct {  
   Name  string  
   Timeout  time.Duration
   Debug    bool  
}
```


Then in the application `main()` the arguments are parsed:  
```
func main() {
    var cfg MyAppConfig
    if err := mainline.NewDecoder(os.Args[1:]).Decode(&cfg); err != nil {
    	t.Fatal(err)
    }
    
    // Then cfg is ready to use...
    ctx, cnl := context.WithTimeout(context.Background(), cfg.Timeout)
    ...
    if cfg.Debug {
        log.Loglevel(log.Debug)
    }
    fmt.Printf("Hello %s", cfg.Name)
}   
```
  
  
Using the command line flags:
`yourapp --name john --timeout 3h --debug`  
  
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
Fields may be tagged to to specify alternative names for the flag using standard go tagging.
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

There is no distiction beween the double dash and single dash for flags.  "-" is the same as "--"  
 
