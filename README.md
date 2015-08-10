# Homemaker #

Homemaker is a tool for straightforward and efficient management of *nix configuration files found in the user's home
directory, commonly known as [dot-files](https://en.wikipedia.org/wiki/Dot-file). It can also be readily used for
general purpose system bootstrapping (installing packages, cloning repositories, etc). This tool is written in
[Golang](https://golang.org/), requires no installation, has no dependencies and makes use of simple configuration file
format inspired by [make](https://en.wikipedia.org/wiki/Make_%28software%29) to generate symlinks and execute system
commands to aid in configuring a new system for use.

![](http://foosoft.net/projects/homemaker/img/demo.gif)

## Motivation ##

Ever since switching to using Linux as my daily driver operating system, I have been searching for a way to effectively
manage settings between different computers (and system reinstalls on the same machine) while avoiding the accumulation
of cruft that plagues our home directories.

Specifically, I required a solution that had the following characteristics:

*   **Do not require software installation**

    The user may not have root privileges on all the computers that she accesses, making it difficult to make use
    software that has additional dependencies. Even if the user has root on the machine in question, it is sub-optimal
    to require her to install software packages just to bootstrap the system.

*   **Support configuration variants**

    The user should be able to use cherry-pick which settings they wish to use. Different system installs have different
    configuration requirements; it must be possible to share common settings while keeping the unique bits unique to
    each machine. Furthermore, it should be possible to store all configuration settings in one location.

*   **Make no assumptions about synchronization**

    While Git often works well for managing dot-files, it is not ideal for every situation. I've seen some applications
    add and remove seemingly randomly-named files within their configuration directories, while others store their
    settings in large, opaque binary blobs. Managing configuration settings for such poorly-behaved applications may be
    easier through, [Dropbox](https://www.dropbox.com/), [rsync](https://en.wikipedia.org/wiki/Rsync), or some other
    utility outside of version control.

*   **Be easy to read and modify the source**

    As a programmer, I always consider the underlying technology of a tool when deciding whether or not to use it. The
    user should feel empowered to change and test the application without having to deal with the archaic incantations
    otherwise known as shell scripts. Scripting languages can work well but are closely tied to the environment in which
    they are executed.

It became apparent to me that utility which I was looking for did not exist. After making due with a simple Python
script that I hacked together, I decided that this problem deserved a clean, formal solution. I decided to create this
new utility in Go; in addition to the language itself being clear and easy to understand, executables built with the Go
compiler are statically linked, making them highly portable. The result of my work is Homemaker; I hope that you find it
suitable for your needs.

## Installation ##

If you already have the Go environment and toolchain set up, you can get the latest version by running:

```
$ go get github.com/FooSoft/homemaker
```

Otherwise, you can use the pre-built binaries for the platforms below:

 * [homemaker_darwin_386.zip](http://dl.foosoft.net/homemaker/homemaker_darwin_386.zip)
 * [homemaker_darwin_amd64.zip](http://dl.foosoft.net/homemaker/homemaker_darwin_amd64.zip)
 * [homemaker_freebsd_386.zip](http://dl.foosoft.net/homemaker/homemaker_freebsd_386.zip)
 * [homemaker_freebsd_amd64.zip](http://dl.foosoft.net/homemaker/homemaker_freebsd_amd64.zip)
 * [homemaker_freebsd_arm.zip](http://dl.foosoft.net/homemaker/homemaker_freebsd_arm.zip)
 * [homemaker_linux_386.tar.gz](http://dl.foosoft.net/homemaker/homemaker_linux_386.tar.gz)
 * [homemaker_linux_amd64.tar.gz](http://dl.foosoft.net/homemaker/homemaker_linux_amd64.tar.gz)
 * [homemaker_linux_arm.tar.gz](http://dl.foosoft.net/homemaker/homemaker_linux_arm.tar.gz)
 * [homemaker_netbsd_386.zip](http://dl.foosoft.net/homemaker/homemaker_netbsd_386.zip)
 * [homemaker_netbsd_amd64.zip](http://dl.foosoft.net/homemaker/homemaker_netbsd_amd64.zip)
 * [homemaker_netbsd_arm.zip](http://dl.foosoft.net/homemaker/homemaker_netbsd_arm.zip)
 * [homemaker_openbsd_386.zip](http://dl.foosoft.net/homemaker/homemaker_openbsd_386.zip)
 * [homemaker_openbsd_amd64.zip](http://dl.foosoft.net/homemaker/homemaker_openbsd_amd64.zip)
 * [homemaker_plan9_386.zip](http://dl.foosoft.net/homemaker/homemaker_plan9_386.zip)
 * [homemaker_snapshot_amd64.deb](http://dl.foosoft.net/homemaker/homemaker_snapshot_amd64.deb)
 * [homemaker_snapshot_armhf.deb](http://dl.foosoft.net/homemaker/homemaker_snapshot_armhf.deb)
 * [homemaker_snapshot_i386.deb](http://dl.foosoft.net/homemaker/homemaker_snapshot_i386.deb)
 * [homemaker_windows_386.zip](http://dl.foosoft.net/homemaker/homemaker_windows_386.zip)
 * [homemaker_windows_amd64.zip](http://dl.foosoft.net/homemaker/homemaker_windows_amd64.zip)

## Configuration Tutorial ##

Configuration files for Homemaker can be authored in your choice of [TOML](https://github.com/toml-lang/toml),
[YAML](http://yaml.org/) or [JSON](http://json.org/) markup languages. Being the easiest to read out of the three, TOML
will be used for the example configuration files. Worry not if you are unfamiliar with this format; everything you need
to know about it will be shown below.

Let's start by looking at a basic example configuration file, `example.toml`. Notice that Homemaker determines which
markdown language processor to use based on the extension of your configuration file. Use `.toml/.tml` for TOML,
`.yaml/.yml` for YAML, and `.json` for JSON. Having a wrong file extension will prevent your configuration file from
being parsed correctly.

```
[tasks.default]
    links = [
        [".config/fish"],
        [".config/keepassx"],
        [".config/terminal"],
        [".config/vlc"],
        [".gitconfig"],
        [".xinputrc"],
    ]
```

We could have just as easily written this configuration in JSON (or YAML for that matter), but it's subjectively uglier:

```
{
    "tasks": {
        "default": {
            "links": [
                [".config/fish"],
                [".config/keepassx"],
                [".config/terminal"],
                [".config/vlc"],
                [".gitconfig"],
                [".xinputrc"]
            ]
        }
    }
}
```

To create symlinks based on this configuration we invoke the `homemaker` utility:

```
$ homemaker example.json /mnt/data/config
```

To get a better idea of what `/mnt/data/config` is, let's look at the in-program documentation:

```
Usage: homemaker [options] conf src

Parameters:
  -clobber=false: delete files and directories at target
  -dest="/home/alex": target directory for tasks
  -force=true: create parent directories to target
  -nocmd=false: don't execute commands
  -nolink=false: don't create links
  -task="default": name of task to execute
  -verbose=false: verbose output
```

For the purpose of our illustration, `src` is defined on the command line to be `/mnt/data/config`; namely the source
directory where your dot-files live (this will be your Git repository, Dropbox folder, rsync root, etc.). The symlinks
that Homemaker creates will point to the configuration files in this directory. You may have noticed that you can also
provide a destination directory via the `-dest` command line argument; this is where the symlinks should be created (and
it defaults to your home directory).

Another useful parameter is `task`; it will be initialized to the value `default` unless you override it on the command
line. In practice, this means that Homemaker will try to find a task called `default` and execute it. You can create as
many unique tasks as necessary to correspond to your configuration requirements, and then choose which one will execute
by specifying it on the command line in the format `-task=taskname`. Good candidates for tasks are computer names, as
shown in the configuration file below:

```
[tasks.flatline]
    links = [
        [".config/syncthing", ".config/syncthing_flatline"],
        [".s3cfg"],
        [".sabnzbd"],
        [".ssh", ".ssh_flatline"],
    ]

[tasks.wintermute]
    links = [
        [".config/syncthing", ".config/syncthing_wintermute"],
        [".ssh", ".ssh_wintermute"],
    ]
```

Here we see two tasks, named after the computers that will be using them, `flatline` and `wintermute`. Certain
configuration data like key pairs and other per-machine settings should only be linked on the computer that is using
them. That is to say if `flatline` and `wintermute` both try to manage the `.ssh` directory, a conflict will occur at
both the source and destination directories. We can easily resolve the source directory conflict by giving the `.ssh`
directories unique names, such as `.ssh_flatline` and `.ssh_wintermute`. The conflict at the destination directory can
be fixed as shown above; we will create per-machine tasks that will symlink only the needed directory.

You may have noticed that each entry in the `links` collection is an array, which up until now has contained only one
item. A second item can be added if the source file or directory name is different from that in the destination. If the
paths provided are relative they will be assumed to be relative to the destination and source directories respectively.

Now that we have machine specific tasks defined in our configuration file, it would be nice to still be able to share
configuration settings that are common to the two computers. We can do this by adding a `dep` array to our tasks as
shown below:

```
[tasks.common]
    links = [
        [".config/fish"],
        [".config/keepassx"],
        [".config/terminal"],
        [".config/vlc"],
        [".gitconfig"],
        [".xinputrc"],
    ]

[tasks.flatline]
    deps = ["common"]
    links = [
        [".config/syncthing", ".config/syncthing_flatline"],
        [".s3cfg"],
        [".sabnzbd"],
        [".ssh", ".ssh_flatline"],
    ]

[tasks.wintermute]
    deps = ["common"]
    links = [
        [".config/syncthing", ".config/syncthing_wintermute"],
        [".ssh", ".ssh_wintermute"],
    ]

```

Homemaker will process the dependency tasks before processing the task itself.

In addition to creating links, Homemaker is capable of executing commands on a per-task basis. Commands should be
defined in an array called `cmds`, split into an item per each command line argument. All of the commands are executed
with `dest` as the working directory (as mentioned previously, this defaults to your home directory). If any command
returns a nonzero exit code, Homemaker will display an error message and stop.

The example task below will clone and install configuration files for [Vim](http://www.vim.org/) into the `.config`
directory, and create links to it from the home directory. You may notice that this task references an environment
variable (set by Homemaker itself) in the `links` block; you can read more about how to use environment variables in the
following section.

```
[tasks.vim]
    cmds = [
        ["rm", "-rf", ".config/vim"],
        ["git", "clone", "https://github.com/FooSoft/dotvim", ".config/vim"],
    ]
    links = [
        [".vimrc", "${HM_DEST}/.config/vim/.vimrc"],
        [".vim", "${HM_DEST}/.config/vim/.vim"],
    ]
```

## Environment Variables ##

Homemaker supports the expansion of environment variables for both command and link blocks. This is a good way of
avoiding having to hard code absolute paths into your configuration file. To reference an environment variable simply
use `${ENVVAR}`, where `ENVVAR` is the variable name. In addition to being able to use all of the environment variables
defined on your system, Homemaker defines a couple of extra ones for ease of use:

*   `HM_CONFIG`

    Path to the homemaker configuration file.

*   `HM_TASK`

    Task name invoked from the command line.

*   `HM_SRC`

    Source directory for link creation.

*   `HM_DEST`

    Destination directory for link creation.

Environment variables can also be set within tasks block by assigning them to the `envs` variable. The example below
demonstrates the setting and clearing of environment variables:

```
[tasks.default]
    envs = [
        ["MYENV1", "foo"],        # set MYENV1 to foo
        ["MYENV2", "foo", "bar"], # set MYENV2 to foo,bar
        ["MYENV3"],               # clear MYENV3
    ]
```

It should be pointed out that it is possible to reference other environment variables using the syntax shown in the
first part of this section. This makes it possible to expand variables like `PATH` without overwriting their existing
value.

## Command Macros ##

It can often be convenient to execute certain commands repeatedly within task blocks to install packages, clone git
repositories, etc. Homemaker provides macro blocks for this purpose; you can specify a command *prefix* and *suffix*
that is used to wrap the parameters you provide. For example, you can declare a macro for `apt-get install` and with the
declaration shown below (much like tasks, macro declarations are global).

```
[macros.apt-install]
    prefix = ["sudo", "apt-get", "install", "-y"]
```

Macros can be referenced from commands by prefixing the macro name with the `@` symbol (it must be the first character
of the first item of a command). For example, the task below installs several python packages using the macro above:

```
[tasks.python]
    cmds = [[
        "@apt-install",
        "python-dev",
        "python-pip",
        "python3-pip",
    ]]
```

This feature makes it possible to reduce the clutter that comes from the repeated commands which must be executed to
bootstrap a new system. When executed with the `verbose` option, Homemaker will echo the expanded macro commands before
executing them.

## Parameters ##

Executing Homemaker with the `-help` command line argument will trigger online help to be displayed. The list below
provides a more detailed description of what the parameters do.

*   `clobber`

    By default, Homemaker will only remove identically-named symlinks at the destination directory. Using this parameter
    will cause Homemaker to be more aggressive and delete clashing files and entire directories as well. This can be
    useful for getting rid of the default configuration settings some applications write when you run them for the first
    time, but should obviously be used with caution.

*  `dest`

    This parameter specifies destination where Homemaker is to create symlinks. This will default to the home directory
    for the current user, and as long as you are just using this application to manage dot-files, will probably never
    need to be changed.

*   `force`

    Sometimes dot-files for an application are nested within parent directories that must exist in order to allow the
    symlink to be successfully created (for example the `.config` directory in `.config/vlc`). As this is the expected
    behavior, this parameter defaults to `true`; however you can explicitly disable it if required. You can specify the
    access permissions of directories created by force by providing a third array item in the link descriptor. For
    example, if you wanted the `.ssh` directory to be created with mode `700`, you could write the following:
    `[".ssh/id_rsa.pub", ".ssh_flatline/id_rsa.pub", "0700"]`. Notice that you can specify permissions in octal notation
    by adding a leading zero value (the `0x` prefix signifies hexadecimal).

*   `nocmd`

    Do not execute commands for the `cmds` blocks inside of tasks.

*   `nolink`

    Do not create links for the `links` blocks inside of tasks.

*   `task`

    This parameter is used to specify which task Homemaker will process when executed. It defaults to the `default`
    task, which should be used when creating a configuration file that does not have system-specific tasks specified.

*   `verbose`

    When something isn't going the way you expect, you can use this parameter to make Homemaker to log everything it is
    doing to console.
