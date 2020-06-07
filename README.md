CSHARP-REFLECT 0.1.0
-

## What is it?
This tool allows for a developer to set build C# classes for use elsewhere.

## Example
```go
import csharp "github.com/ptdave20/csharp-reflect"
```

```go
    options := csharp.New()
    options.OutputPath = "./Output"
    csharp.ConvertObject(YourClass,options)
```

## Todo
- Add JsonProperty support
- Add BsonProperty support
- Code templating
- Create CLI
- Better error handling