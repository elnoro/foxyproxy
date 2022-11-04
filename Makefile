releaseC:
	docker run --rm -e GITHUB_TOKEN -v `pwd`:/app -w /app goreleaser/goreleaser release
