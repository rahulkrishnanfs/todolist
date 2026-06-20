//go:build mage

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// Build the binary locally
func Build() error {

	return sh.Run("go", "build", "-o", "./bin/app", "./cmd")
}

// login to docker registry wiith the username and password from the env variable
func DockerLogin() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return sh.RunWith(
		map[string]string{
			"DKUSERNAME": os.Getenv("DKUSERNAME"),
			"DKPASSWORD": os.Getenv("DKPASSWORD"),
		}, "sh", "-c",
		`echo $DKPASSWORD | docker login -u $DKUSERNAME --password-stdin`)
}

// Build the docker image and tagged with version
func DockerBuild() error {

	return sh.Run("docker", "build", "-t", "rahulkrishnanfs/todolist:v1", ".")
}

// Push the docker image to the docker image registry
func DockerPush() error {
	return sh.Run("docker", "push", "rahulkrishnanfs/todolist:v1")

}

func All() {
	mg.Deps(Build, DockerLogin, DockerBuild, DockerPush)
}
