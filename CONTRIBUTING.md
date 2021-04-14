# Contributing

The commands below might be useful for contributors, including myself after forgetting them : ).

```
# run tests
go test -v

# install from source
go install github.com/felixge/pprofutils/cmd/...

# cut a new release
echo "v0.3.0" > version.txt
git add version.txt
git commit -m "Release $(cat version.txt)"
git tag $(cat version.txt)
git push --tags
```
