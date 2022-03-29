Contributing
============

I'd appreciate your help making `attest` better! Please keep in mind that I'm
trying to keep this package small - I don't want it to have the same breadth of
assertions as [testify](https://pkg.go.dev/github.com/stretchr/testify). If
you'd like to add some new APIs, please open an issue to discuss before making
a PR. :heart_eyes_cat:

Most importantly, please remember to treat your fellow contributors with
respect!

## Build and test

`attest` is simple: you can build and test with the usual `go test .`, or you
can use the Makefile to match CI more closely (`make help` lists the available
targets). If you're opening a PR, you're most likely to get my attention if
you:

* Add tests to cover your changes.
* Write a [good commit message][commit-message].
* Maintain backward compatibility.
* Stay patient. This isn't my day job. :wink:

[commit-message]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
