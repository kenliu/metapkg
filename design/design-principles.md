# Design Principles for metapkg

1. prioritize simplicity over power and feature-richness
    - lean toward simplicity in the design of the tool and user experience
    - it's a helpful layer on top of package managers but it's not a replacement for them or an abstraction layer

2. easy to use
    - simple configuration file format
    - prioritize human-readability over machine-readability
    - does common things well
    - should be intuitive to do things without referring to the documentation
    - should be easy to learn and understand

3. sensible defaults
    - some conventions are used to make common things easier
    - some conventions will be implicit

4. scalable complexity
    - simple by default, easy to specify more complex things if needed
    - doesn't try to do too much
    - let package managers do what they do best
    - doesn't handle uninstalling packages

5. assume that users will be able to use advanced functionality of package managers if they need to

6. built for linux and mac: unix-like systems with package managers

7. easy to extend
    - design the structure of the code to make it easy to add support for new package managers on different distros
    - but don't plan to keep making the tool more abstract to support more and more distros (especially the niche ones)
    - prioritize the most common distros and use cases

8. prioritize modern tools over outdated / legacy ones
    - for example: support dnf instead of yum
    - apt instead of dpkg

9. when in doubt -- it's meant to solve my own problems and I'm the main target user. Unusual use cases won't be solved.

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

