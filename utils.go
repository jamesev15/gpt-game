package main

import "os"

func readEnvVar(env string) string {
	return os.Getenv(env)
}
