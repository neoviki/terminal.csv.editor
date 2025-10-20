go clean -cache        # clears the build cache
go clean -modcache     # clears the module download cache (optional, only if you want to force redownload)
rm csvgo
go build -o csvgo csvgo.go

