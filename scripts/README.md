# Scripts

Scripts to manage the project. They can also be invoked from the Makefile
at the project's root.

All scripts require Zsh. [Variable modifiers](http://zsh.sourceforge.net/Doc/Release/Expansion.html#Modifiers)
are too useful to consciously choose to use `/bin/sh` like a caveman. A recent
(enough) version of zsh installed in Macs by default, and every distro has a zsh
package.

| Script                | Purpose                                                             |
| ---                   | ---                                                                 |
| `build.zsh`           | meta build script, excluding container builds and integration tests |
| `build-mitmproxy.zsh` | builds the mitmproxy container                                      |
| `build-mittens.zsh`   | builds the kubectl-mittens binary                                   |
| `test.zsh`            | unit tests                                                          |
| `ig-test.zsh`         | integration tests                                                   |
| `_pre.zsh`            | run by other scripts - handles output formatting and PWD storage    |
| `_post.zsh`           | run by other scripts - restores PWD state from `pre.zsh`            |
