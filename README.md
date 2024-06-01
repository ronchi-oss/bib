# bib

bib is a note-taking CLI utility written in Go.

For users comfortable working in a command line, it is both minimalistic and powerful. Though minimal shell scripting knowledge will allow a user to master bib in a few hours, this document is intended for anyone using a Linux or macOS computer on a regular basis. Windows isn't officially supported yet. Mobile operating systems such as iPadOS and Android probably never will be.

bib operates on *target directories*. A target directory is a directory in your computer where all your notes are stored. Though bib makes it pretty convenient to manage notes spread across dozens of target directories, the author suggests, for most people, to use a single target directory for all their notes; for working professionals, it's recommend to use one target for their private notes and one extra target per active or past employer.

## Installation

### Building from source

If building from source, you'll need `go` installed. The following program will build and install the `bib` binary to `$GOPATH/bin`.

```sh
git clone https://github.com/ronchi-oss/bib.git
cd bib
go install
```

### Shell completion

If installed using Homebrew, shell completion will be automatically installed.

If building from source, the setup depends on your shell of choice and operating system. For enabling bash completion on a Linux distribution, you could add the following to `~/.bashrc`.

```sh
if command -v bib >/dev/null; then
    eval "$(bib completion bash)"
fi
```

### Global config file

bib reads a global YAML config file in order to perform certain operations, however bib won't create this file on its own: before you start using bib you should create it yourself with the default contents.

bib will look for the global config file in two different locations, in this order:

1. `$BIB_GLOBAL_CONFIG` (the override)
1. `$XDG_CONFIG_HOME/bib/bib.yml` (the default)

Regardless of your choice, you can conveniently create this file with the default YAML definition by running either:

```sh
bib generate default-global-config > $BIB_GLOBAL_CONFIG
```

or:

```sh
bib generate default-global-config > $XDG_CONFIG_HOME/bib/bib.yml
```

## Getting started

`bib init` is a handy command that, given an empty directory, will:

* Create a basic `bib.yml` with [two filters](#querying-notes)
* Create your first note
* Initialize it as a git repository and commit all created files

The following commands will achieve that:

```
mkdir ~/bib-private-notes
bib init ~/bib-private-notes
```

`bib init` will ask you to confirm that you'd like to create the bib structure inside that directory. If you type `yes` and hit `Enter`, bib will do its work and you should end up with the following directory structure inside `~/bib-private-notes`:

```
.
├── bib.yml
└── src
    └── 1
        ├── metadata.yml
        └── README.md
```

All bib commands that operate on a single note expect to be provided with a note code; our first note code is `1`, and our next will be `2`. A note is composed of two files:

* A `metadata.yml` file, which stores context about your note as structured data (e.g. its creation date and time)
* A `README.md` file, which is the note you will edit in a text editor

We can output our first note metadata with:

```sh
bib show note 1 -d ~/bib-private-notes
```

We can output our first note Markdown file with:

```sh
bib cat 1 -d ~/bib-private-notes
```

We can edit our first note Markdown file with (uses the `EDITOR` environment variable):

```sh
bib edit 1 -d ~/bib-private-notes
```

The following sections will build on these basics and explain the features that bib provides in order to create, edit, browse and filter notes more efficiently.

## Managing bib target profiles

A target *profile* is a name that saves you some typing whenever you must tell bib which directory it should operate on.

If you follow the installation instructions, after installing bib you then create your global config file. As of the time of writing that file holds the list of profiles. You can list all profiles in that file by running:

```sh
bib get profiles
```

On a fresh installation, there should be no output, i.e. no profiles have been defined.

Let's create a profile for the target directory we initialized in the "Getting started" section:

```sh
bib create profile private ~/bib-private-notes
```

Running `bib get profiles` will show it:

```
private ~/bib-private-notes
```

Here's why profiles are handy: whenever a bib command accepts a `-d` (long form: `--target-dir`) option, it also accepts a `-p` (long form: `--profile`) option. Which means, for the hypothetical setup we have for this tutorial, the following two commands are equivalent:

```sh
bib cat 1 -d ~/bib-private-notes
bib cat 1 -p private
```

While the `-p <profile-name>` form allows us to type a bit less, we are still typing quite a lot. The good news is, if you leave out both `-d` and `-p` options, bib will attempt to use the value of the `BIB_PROFILE` environment variable as the profile name:


```sh
export BIB_PROFILE=private
bib cat 1
```

## Taking notes

We can create a new note with:

```sh
bib create note -d <target-dir>
```

This will create both the metadata and Markdown files, and will start a sub-process using the value of `EDITOR` as the editor and the path to the Markdown file as the sole argument.

## Pinning notes

Each note metadata contains a `pinned` property with a boolean value (`false` or `true`). By default, new notes will not be pinned (`pinned` = `false`). A pin is a simple way of marking a note as "special" (i.e. to be followed up on).

You can toggle the `pinned` property of a note using the `toggle-pin` command. For instance:

```sh
bib toggle-pin 1 -d <target-dir>
```

## Querying notes

You can list all your notes (metadata plus title) in TSV format using:

```sh
bib get notes -d <target-dir>
```
Consider a bib target directory with 5 notes, where notes 1 and 3 are pinned. The command above would output something similar to:

```
1	*	2024-05-31 22:26:28	My first note
2	 	2024-05-31 22:27:06	Bookmark: Bell Labs original Unix Programmer's Manual, 1971
3	*	2024-05-31 22:28:56	Bookmark: How to Write Go Code (go.dev)
4	 	2024-05-31 22:30:43	A note
5	 	2024-05-31 22:31:29	Another note
```

While we can certainly use tools such as `awk`, `grep`, `sed` etc. to filter these notes, bib provides a build-in feature called `filter` that achieves the same results with less typing.

By default, `bib init` will add two filters (`all` and `pinned`) to the target directory `bib.yml` definition:

```yaml
filters:
  - name: all
    cmd: cat
    cmd_args: []
  - name: pinned
    cmd: bib-pin-filter
    cmd_args: []
```

A filter is a command that bib will pipe the output of `get notes -f <filter-name>` into. Considering the definition above, the following commands are equivalent:

```sh
# list all notes
bib get notes -d ~/bib-private-notes
bib get notes -d ~/bib-private-notes -f all

# filter notes where `pinned` = `true`
bib get notes -d ~/bib-private-notes | bib-pin-filter
bib get notes -d ~/bib-private-notes -f pinned
```

The `-f pinned` filter would output:

```
1	*	2024-05-31 22:26:28	My first note
3	*	2024-05-31 22:28:56	Bookmark: How to Write Go Code (go.dev)
```

What about filtering all notes whose title start with "Bookmark: "? Since filtering by a note title pattern allows the creation of ad-hoc workflows, bib ships with an accompanying script `bib-title-filter` which takes an `awk` pattern as its only argument. Considering that script is in your PATH, the here's the bib filter definition that satisfies the query:

```yaml
filters:
  # ...
  - name: bookmark
    cmd: bib-title-filter
    cmd_args: ["^Bookmark: "]
```

Running `bib get notes -d ~/bib-private-notes -f bookmark` will output:

```
2	 	2024-05-31 22:27:06	Bookmark: Bell Labs original Unix Programmer's Manual, 1971
3	*	2024-05-31 22:28:56	Bookmark: How to Write Go Code (go.dev)
```

You should consider adopting the `-f <filter-name>` style because:

* bib will auto-complete the filter name, so typing `-f p<TAB>` will automatically expand to `-f pinned`
* the `bib-filter-fzf` script relies on filters (see [browsing notes with fzf](#Browsing-notes-with-fzf))

## Browsing notes with fzf

While not providing its own pager, bib ships with two shell scripts that serve as guidance in how to use `fzf` in order to create a powerful note browser. The pager is initialized with multiple keybindings that integrate with bib:

* filter by title using fzf query language
* toggle pin on a note (`ctrl-]`)
* edit a note without leaving the pager (`ctrl-j`)
* refresh the list (`ctrl-r`)
* toggle note preview visibility (`ctrl-i`)
* toggle note preview position (`ctrl-space`)

The script that's most straight-forward of the two is `bib-filter-fzf`, which presents you with a list of filters to choose from, and once you choose it, loads an fzf pager applying that filter to your bib notes. This script requires `BIB_PROFILE` to be set.

## Hooks

Hooks are external commands that bib will invoke whenever matching events happen. As of the time of writing the only event that bib will notify hooks about is `note.created`. Here's how you'd configure your target directory `bib.yml` config file in order to notify a hypothetical `bib-hook` script every time bib creates a note:

```yaml
filters:
  # ...
hooks:
  - events: ["note.created"]
    cmd: bib-hook
```

For the `note.created` event, bib will invoke the hook command with the following positional arguments:

1. The note numerical identifier
1. The note title
