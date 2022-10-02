# Description
This command line program scans files of various programming languages, counts the number of lines of code with and without comments and automatically skips other file types like dotfiles (.git, ...), text files and so on.

List of currently supported languages:
- C/C++
- Go
- Java
- JavaScript
- Typescript
- Python
- Rust
- Ruby
- PHP
- R

Any other language can easily be added in the file `parse/parse_file.go` as long it has the same comment style as any of the languages above.

# Install
Just clone the repository and then use `go install loc` inside it to add the program to your `GOPATH`. After that you can use it in your terminal.

# Usage
You can query all the possible flags as well as usage instructions by entering `loc -h` into the terminal.