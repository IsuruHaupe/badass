package main

import "os"

func SetBDDEnvironmentVariable() {
	os.Setenv("DBUSER", "root")
	os.Setenv("DBPASS", "root")
}
