# Contributing

The commands below might be useful for contributors, including myself after forgetting them : ).

Unless noted otherwise, your working dir should be the root of this repo.

```
# run tests
go test -v

# install from source
go install github.com/felixge/pprofutils/cmd/...

# cut a new release
vim version.txt
git add version.txt
git commit -m "Release $(cat version.txt)"
git tag $(cat version.txt)
git push
git push --tags
```
