# 📦 metapkg

A lightweight, local package management tool for managing packages across different package managers.

## 🤔 What is metapkg?

`metapkg` is a straightforward command line tool designed to help users manage the installation of packages on their local machines across multiple package managers such as `apt`, `dnf`, and `snap`. It provides a simple way to install and track packages with a simple configuration file without the complexity of full-scale package managers such as Ansible or Nix.

## ✨ Features

- 💻 Local package installation
- 🤝 Compatibility with existing system package managers
- 📘 Simple configuration file
- 🔄 Install package updates (use your system package manager for this)
- 🧪 Handle scripted installations using system shell commands

## ⚠️ Limitations

To keep things simple, metapkg does not:

- 🚫 Automatically install package dependencies
- 🌐 Manage packages on remote machines
- 🔄 Install package updates directly (use your system package manager for this)
- 🧪 Handle complex installation recipes

## 🚀 Getting Started

To get started with metapkg, follow these steps to download the Go source and build it:

1. Ensure you have Go installed on your system (version 1.16 or later is recommended).

2. Clone the repository:
   ```
   git clone https://github.com/kenliu/metapkg.git
   cd metapkg
   ```

3. Build the project:
   ```
   go build -o bin ./...
   ```

   This command will compile all packages in the project and its subdirectories, creating an executable named `metapkg` in the `bin` directory.

4. (Optional) Add the bin directory to your PATH for easier access:
   ```
   export PATH=$PATH:$(pwd)/bin
   ```

   You can add this line to your shell configuration file (e.g., `.bashrc`, `.zshrc`) to make it permanent.

Now you can run metapkg by executing `./bin/metapkg` from the project directory, or simply `metapkg` if you've added it to your PATH.

## 📘 Usage

[Basic usage instructions and examples will be added here]

## 🤝 Contributing

Contributions are welcome! If you have ideas for improvements or have found a bug, please feel free to open an issue or submit a pull request.

## 📜 License

`metapkg` is open source software [licensed as MIT](LICENSE). Copyright (c) 2024 Ken Liu.

## 🆘 Support

If you encounter any problems or have questions, please open an issue in this repository.   

## 👏 Acknowledgements

[Recognize contributors and any third-party resources here]

Thank you for checking out `metapkg`. I hope it simplifies your local package management!