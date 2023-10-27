package main

import "argparse"

func main() {
	parser := argparse.NewParser("My Program", "This is a description of my program", "My Program v1.0")
	parser.AddArg("-c", "--cidr-arg", "cidr range to scan", "cidrrange", "192.168.0.0/16", true)
	parser.AddArg("-s", "--string-arg", "a string argument", "string", "", false)
	parser.AddArg("-I", "--integer-arg", "an integer argument", "integer", "1", true)
	parser.AddArg("-i", "--ip-address-arg", "an ip address argument", "ipaddress", "1", true)
	parser.AddArg("-b", "--bool-arg", "a boolean argument", "bool", "false", true)
	parser.Parse()
	s := parser.Get("string-arg").(string)
	fmt.Println(s)
}
