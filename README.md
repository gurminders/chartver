### Overview
Helm Charts often have git-commit-id as part of the version and it is difficult to find the latest version of a Helm Chart.

**chartver** is a utility to find the Helm Chart versions.

### Install
```shell
# change to HOME directory
cd

# install chartver executable
go install github.com/gurminders/chartver/chartver
```

### Usage

```shell
# get all the Helm Chart from a repository
chartver

# get a particular chart
chartver elasticsearch

# get list of charts
chartver elasticsearch influxdb prometheus-operator

# command line options
chartver -h
```
