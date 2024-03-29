# ArgParse
## A simple argument parsing library for Golang

<div align="center">
    <img src="https://img.wanman.io/fUSu0/SaCUyEMe87.png/raw" style="width: 350px; height: auto" alt="Ar1ste1a" title="Ar1ste1a Offensive Security">
</div>

# Getting Started

To begin, import the argparse library
``` go
import "github.com/Ar1ste1a/ArgParse"
```

Next,  create a new argparse object with the following information:
- Program Name
- Program Description
- Program Details (Typically Version Information)

This will allow argparse to dynamically generate a help banner for your convenience
  
``` go
parser := argparse.NewParser("Program name", "Program Description", "Program Details")
```

Next, add your argument(s) to the argparse object. Argparse can currently interpret the following primitive and composite data types using the accompanied string value. Argparse utilizes casting and interfaces to provide each stored value as it was originally intended to the users program.
- Ints (integer)
- Strings (string)
- Booleans (bool)
- IP Address (ipaddress)
- CIDR Network (cidrrange)

``` go
# Adds the flag "cidr" to the argparse object allowing for the parameters "-c" and "--cidr".
# Stores the object as a CIDR Network with the default value of "192.168.0.0/16"
# Adds the argument as a required parameter
parser.AddArg("-c", "--cidr", "cidr range to scan", "cidrrange", "192.168.0.0/16", true)
```

Parse the arguments passed to the program

``` go
parser.Parse()
```

You can now access any parameters passed by the user by using "Get"
``` go
cidr := parser.Get("cidr").(*net.IPNet)
```


Todo: 
1. Add dependencies based on other parameters
