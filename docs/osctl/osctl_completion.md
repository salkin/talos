<!-- markdownlint-disable -->
## osctl completion

Output shell completion code for the specified shell (bash or zsh)

### Synopsis

Output shell completion code for the specified shell (bash or zsh).
The shell code must be evaluated to provide interactive
completion of osctl commands.  This can be done by sourcing it from
the .bash_profile.

Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2

```
osctl completion SHELL [flags]
```

### Examples

```
# Installing bash completion on macOS using homebrew
## If running Bash 3.2 included with macOS
	brew install bash-completion
## or, if running Bash 4.1+
	brew install bash-completion@2
## If osctl is installed via homebrew, this should start working immediately.
## If you've installed via other means, you may need add the completion to your completion directory
	osctl completion bash > $(brew --prefix)/etc/bash_completion.d/osctl

# Installing bash completion on Linux
## If bash-completion is not installed on Linux, please install the 'bash-completion' package
## via your distribution's package manager.
## Load the osctl completion code for bash into the current shell
	source <(osctl completion bash)
## Write bash completion code to a file and source if from .bash_profile
	osctl completion bash > ~/.talos/completion.bash.inc
	printf "
		# osctl shell completion
		source '$HOME/.talos/completion.bash.inc'
		" >> $HOME/.bash_profile
	source $HOME/.bash_profile
# Load the osctl completion code for zsh[1] into the current shell
	source <(osctl completion zsh)
# Set the osctl completion code for zsh[1] to autoload on startup
osctl completion zsh > "${fpath[1]}/_osctl"
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --context string       Context to be used in command
  -e, --endpoints strings    override default endpoints in Talos configuration
  -n, --nodes strings        target the specified nodes
      --talosconfig string   The path to the Talos configuration file (default "/home/user/.talos/config")
```

### SEE ALSO

* [osctl](osctl.md)	 - A CLI for out-of-band management of Kubernetes nodes created by Talos

