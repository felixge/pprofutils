# Contributing

The commands below might be useful for contributors, including myself : ).

```
# run tests
go test -v

# cut a new release
echo "v0.3.0" > version.txt
git add version.txt
git commit -m "release $(cat version.txt)"
git tag $(cat version.txt)
git push --tags
```
