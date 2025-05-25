# tekton PEG (parsing expression grammar) trial

A parser which uses PEG using Pigeon to parse PEG based grammar in golang to list Tasks in Tekton and then run them.

### Setup Environment

Create a Kubernetes Kind cluster and install Tekton Pipelines.

```
kind create cluster
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```

### Create test resources 

Create test resources to test the application against.

```
kubectl create example-pipelinerun.yaml
```

### Build project

Let's build the project and then start using it.

```
pigeon tekton.peg > main.go
go build .
```

After runnning the above you will have a binary called `tekton-pigeon` you can use to parse your commands and run your own langauge based off of PEG.

```
./tekton-pigeon <your-command>
```

### List Tasks

Run the following command to list tasks installed.

```
./tekton-pigeon "list task"
= Tasks: echo, run-js
```

### Create TaskRun based on an existing Task

Run the command below to start your TaskRun from and existing Task
```
./tekton-pigeon "run task echo"                    
= Started TaskRun for task: echo
```