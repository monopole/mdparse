module github.com/monopole/mdparse

go 1.21

require (
	github.com/gomarkdown/markdown v0.0.0-20231115200524-a660076da3fd
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.8.4
	github.com/yuin/goldmark v1.6.0
)

replace (
	github.com/gomarkdown/markdown => ../../gomarkdown/markdown
	github.com/yuin/goldmark => ../../yuin/goldmark
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
