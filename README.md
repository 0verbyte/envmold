# envmold
**This project is under active development. [Discussions are open](https://github.com/0verbyte/envmold/discussions) and should be used to discuss topics related to this project.**

CLI to create environment variables from a predefined template, molding your environment into use.

# Usage

Create mold template yaml file. The syntax is to create environment variables with the following format for each.

Note: mold expects an array of environment variables with the following structure.
```yaml
  # environment variable name.
- name: foo
  # environment variable value. Type should match the `type` field.
  value: "bar"
  # type of the `value` field. Available types: string | number | boolean
  type: string
  # whether the variable is required or not. If required is true, then mold will ask to fill the value.
  required: false
```

The following command line options are available for mold.
```
% ./bin/mold -h
Usage of mold (v0.1.0):
  -debug
    	Enables debug logging
  -output string
    	Where environment variables will be written. File path or stdout (default "stdout")
  -template string
    	Path to the mold environment template file (default "mold.yaml")
```

Example using a file (`mold.yaml`) relative to the mold binary.
```
./bin/mold
```
# Development

Install Go for [your system](https://go.dev/dl/) and verify the installation.
```bash
go version
```

Install golangci-lint for [your system](https://golangci-lint.run/usage/install/) and verify the installation.
```bash
golangci-lint version
```

Next, verify that the project builds on your system by running `make lint && make test`. If this step fails, [please
create an issue](https://github.com/0verbyte/envmold/issues/new/choose).

When making a PR, a few GitHub actions will run to verify that the suggested changes meet the project code guidelines. Each
check must pass before the changes can be merged into the `main` branch.
