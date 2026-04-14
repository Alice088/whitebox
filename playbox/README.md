# Whitebox Context & Playbox

## Why there is no default behavior

This system intentionally does not include any “out-of-the-box” behavior.
There are no preloaded tools, skills, memories, or minds that automatically affect execution.

The reason is simple: control over the context belongs entirely to the user.

Any predefined behavior is a decision made in advance on behalf of the user. Even if convenient, it reduces transparency and limits control. Whitebox avoids that.

## The idea of Playbox

Instead of defaults, the system uses the concept of a **playbox**.

Playbox is just a folder containing optional components:

* tools
* skills
* memories
* minds

Nothing from it is loaded automatically.

The user decides:

* whether to use anything from the playbox
* what exactly to take
* how to combine it
* or to ignore it completely

## Why this exists

This approach provides several properties.

**Full transparency**
The context is built only from what the user explicitly includes.

**No hidden logic**
There is no implicit behavior influencing results.

**Flexibility**
Any context can be assembled for a specific task without fighting defaults.

**Reproducibility**
The same set of files always produces the same result.

## How it works

The context is constructed from directories:

* `minds/`
* `memories/`
* `skills/`
* `tools/`
* `sessions/`

Order matters. It defines how the final prompt is built:

1. minds
2. memories
3. skills
4. tools
5. session history

Playbox is simply an external source for these files.

## How to use playbox

1. Open the playbox
2. Select the needed components
3. Copy them into the corresponding directories (`~/.whitebox/context/...`)
4. Run the system

Or do nothing.

## Principle

Whitebox does not impose behavior.
It provides a mechanism.

Playbox is not active by default.
It exists only as an option.

Context is the responsibility of the user.
