![Go](https://github.com/arodriguezdlc/sonatina/workflows/Go/badge.svg)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/golang-standards/project-layout)
[![Go Report Card](https://goreportcard.com/badge/github.com/arodriguezdlc/sonatina)](https://goreportcard.com/report/github.com/arodriguezdlc/sonatina)
[![codecov](https://codecov.io/gh/arodriguezdlc/sonatina/branch/master/graph/badge.svg)](https://codecov.io/gh/arodriguezdlc/sonatina)

# Sonatina

Sonatina is an infrastructure as code framework based on Terraform, that adds some useful functionality to organize and work with terraform modules and variables in an opinionated way.

The key features of Sonatina are:

- Based on Hashicorp [Terraform](https://github.com/hashicorp/terraform) (Opensource version).

Warning: private repository are not supported yet, but it will very soon.

##  Concepts

TODO

## Getting Started

### Get sonatina

Using sonatina is easy. First, clone the repository:
```sh
git clone https://github.com/arodriguezdlc/sonatina.git
```

Sonatina uses native Go Modules to manage dependencies, so you only have to execute go build and
every dependencies will be fetched for you
```sh
cd sonatina
go build 
```

As result, you will have a sonatina binary ready to work. Optionally, you can move it to a directory in your `$PATH`:
```
mv sonatina /usr/local/bin/
```

### Create your first deployment

Sonatina stores deployment state, variables and metadata on a git repository. For this example we are going to create a local repository,
but you also could create one on [github](https://github.com/new) for example. 
Let's assume we hace created this repo: https://github.com/arodriguezdlc/sonatina-example-local-docker-storage.git


Also we need a repository with the Terraform code we are going to deploy. 
We are going to use [an example repository](https://github.com/arodriguezdlc/sonatina-example-local-docker) to deploy a local Docker with a web server. 


Finally, we can create the deployment with the following sonatina command:
```sh
sonatina create deployment example-local-docker -s https://github.com/arodriguezdlc/sonatina-example-local-docker.git -c https://github.com/arodriguezdlc/sonatina-example-local-docker.git
```

You can see the created deployment with:
```
sonatina list deployments
```

Now we can start working with deployment variables:
```sh
sonatina init
sonatina edit
```

Let's assign a name for your deployment (the required variable), and deploy it:
```sh
sonatina apply
```

And this is all! You have deployed a local docker with a http server listening to the 8080 port. Using variables you can customize other configurations
like the port or the sentence that is being shown. 

This has been a global deployment but you can deploy one customized server for a specific client. Let's do that creating an user component:
```sh
sonatina create usercomponent my-special-client
sonatina init -c my-special-client
```

Let's specify a custom sentence ("Hello my special client") for him editing the `sentence` variable. In addition, modify the `port` variable to 8081 to avoid port conflict with global component.
```sh
sonatina edit
```

And finally, apply changes:
```sh
sonatina apply -c my-special-client
```

### Adding a plugin

TODO

### Cleanup

To perform the undeploy, you simply have to execute a sonatina destroy command:
```sh
sonantina destroy -c my-special-client
sonatina destroy
```

Also you can delete the deployment from Sonatina. Note that this operation won't delete
the storage repository, so you recover by cloning it. 
```sh
sonantina delete deployment example-local-docker
```
 
## Contributing
