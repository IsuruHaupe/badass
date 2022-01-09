package main

import "os"

func SetBDDEnvironmentVariable() {
	os.Setenv("DBUSER", "admin")
	os.Setenv("DBPASS", "$Rootroot1")
}
