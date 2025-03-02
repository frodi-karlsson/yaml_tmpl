[![Go Reference](https://pkg.go.dev/badge/github.com/frodi-karlsson/yaml_tmpl.svg)](https://pkg.go.dev/github.com/frodi-karlsson/yaml_tmpl)

## yaml_tmpl

This project is a hopefully not useful go server that leverages yaml as a templating language.
It's just a way for me to experiment with Go.

### Installation

You can install the package with `go get github.com/frodi-karlsson/yaml_tmpl`

### Usage

- Run the example server from /main with `go run main.go`
- Visit `http://localhost:8080` in your browser
- Alternatively, you can build the example site statically with `go run main.go -static`. You'll find the output in /main/docs
- Probably actually don't use this for anything. It's just a toy.

### Templating Logic

- Any non-indented key is an html tag
- If its value is a string, that's the content of the tag. This is a short hand of children: raw: "{{ . }}"
- If its value is a map, the keys are attributes of the tag
- An attribute with a key of 'children' is a list of child tags
- An attribute with a key of 'innerText' will be the inner text of the tag. This is a short hand of children: raw: "{{ .innerText }}"
- Children of 'children' are html tags
- 'raw' as a child is parsed as a raw html string
- For development simplicity, and lack of need, there is no difference between a sequence and a mapping
- You can use YAML aliases and anchors to repeat content
- The value of an anchor is not transpiled until it's aliased. This allows you to separate definition from use
- Overrides are also possible using "<<: *anchor", although I don't quite know if they behave in a sane way

### Other

After finishing this toy project, I stumpled upon someone with a similar idea: [Yaml2Html](https://metacpan.org/release/RJE/YAML-Yaml2Html-0.5/view/lib/YAML/Yaml2Html.pm). Very cool that someone had the same idea in 2005, and took it in such a different direction syntax-wise.