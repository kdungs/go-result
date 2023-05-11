# then

The _lazy_ sibling of [result](https://github.com/kdungs/go-result/result). Why would this be useful? Take a look at [example_test.go](example_test.go).


## Caveats

 - Right now, there's no good way to handle resource management (defer Close); maybe we need some sort of tee + resource manager...
