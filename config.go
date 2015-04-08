package main

type Flags struct {
	isJson   bool
	isForm   bool
	verbose  bool
	download bool
	output   string
	auth     string
	authType string
	proxy    string
	cert     string
	certKey  string
	debug    bool
}
