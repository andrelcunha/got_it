[![Go](https://github.com/andrelcunha/got_it/actions/workflows/go.yml/badge.svg)](https://github.com/andrelcunha/got_it/actions/workflows/go.yml)
# *Got it!* Version Control System

Got_it is a simple, Git-like version control system written in Go. It allows you to initialize repositories, add files to the staging area, commit changes, and view the commit history.

## Features

- Initialize a new repository
- Add files to the staging area
- Commit changes
- View commit history

## Installation

To install Got_it, clone this repository and build the project:

```sh
git clone https://github.com/yourusername/got_it.git
cd got_it
make build

## Usage
### Initialize a new repository
```sh
./got init
```
### Add files to the Staging Area
```sh
./got add <file1> <file2> ...
```
### Commit Changes
```sh
./got commit -m "Your commit message"
```
### View Commit History
```sh
./got log
```
## Contributing
If you'd like to contribute to Got_it, please fork the repository and create a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.
## License
This project is licensed under the MIT [License](https://github.com/andrelcunha/got_it/blob/main/LICENSE).
