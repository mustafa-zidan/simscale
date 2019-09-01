# Simscale Coding Challenge

A golang coding task for simscale

# Prequisites
- [Go 1.12+](https://golang.org/)
- [Docker 2.1+](https://www.docker.com/)


## Getting Started
### Install build tools.
Although go modules is taking care of this through indirect dependency if you want to have it manually installed
please run
```bash
make deps #Just a neat way to install github.com/mitchellh/gox
```
> **NOTE**: if you are having the repo cloned under `$GOPATH` you might want to set `GO111MODULE` to `on`

```bash
export GO111MODULE=on
```

### Run
To Run the app locally use

```bash
make run i=<PATH/TO/INPUT/FILE> o=<PATH/TO/OUTPUT/FILE>
```
The app itself use `--input-file`, `-i` for input file and `--output-file`, `-o` for output file

### Build the app.

To be able to cross compile the app for differet platforms I user [gox](https://github.com/mitchellh/gox)

```bash
make build
```

### Docker
The app can be fully dockerized and this can be done by running

```bash
make run/docker
```

This command will build the app, image and run the docker container.
To build the docker image only run

```bash
make build/docker
```
> **NOTE**: This command expects 3 env variables to be set:
> - `DOCKER_MOUNT_DIR`: directory where the input file resides
> - `DOCKER_INPUT_FILE`: name of input file e.g., `small-log.txt`
> - `DOCKER_OUTPUT_FILE`: name of output file e.g., `small-trace.txt`

### Test

To run unit tests

```bash
make test
```

To include benchmark tests

```bash
make test/benchmark
```
## More Info
Check Out the overview to understand more about the datastructure in [Overview](/OVERVIEW.md) doc
