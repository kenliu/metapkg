# Design Principles for metapkg

- prioritize simplicity over power and feature-richness
    - lean toward simplicity in the design of the tool and user experience

- easy to use
    - simple configuration file format
    - does common things well
    - should be intuitive to do things without referring to the documentation
    - should be easy to learn and understand
- sensible defaults
    - some conventions to make common things easier

- scalable complexity
    - simple by default, easy to specify more complex things if needed

- it's a helpful layer on top of package managers but it's not a replacement for them or an abstraction layer

- doesn't try to do too much
    - let package managers do what they do best
    - doesn't handle uninstalling packages

- assume that users will be able to use advanced functionality of package managers if they need to

- built for linux and mac: unix-like systems with package managers

- easy to extend
    - design the structure of the code to make it easy to add support for new package managers on different distros
    - but don't plan to keep making the tool more abstract to support more and more distros (especially the niche ones)
    - prioritize the most common use cases and distros

- prioritize modern tools over outdated / legacy ones
    - for example support dnf instead of yum
    - apt instead of dpkg

- when in doubt -- it's meant to solve my own problems and I'm the main target user. Esoteric use cases won't be solved.

# Non-goals

- doesn't cater to advanced users
    - advanced users can use other tools like nix and ansible instead of this tool

- not designed to work with windows or powershell
    - but maybe it could work with WSL

# Doesn't do
- dependency management
- support obscure distros / package managers
- complex workflows
- try to manage complete system state like terraform

