namespacer
====

For all C# codefiles under a given path, ensures they have a specific namespace.

Is this just a one-off tool?
====

Yeah. I wouldn't probably recommend using it if you're not me. It'll work just fine, but it's not meant to be some full-fledged GNU-quality tool. It does an absurdly specific job that I needed done one morning; no more no less.

What are the caveats?
====

* Uses string manipulation to determine this, rather than parsing and AST. This can lead to inconsistent results if your comments look like class/enum/struct signatures.

* Only does one namespace per file. Seriously, why would you do anything else?

* Depends upon having an explicit access modifier on your type signatures. Because seriously, why would you _not_ do that?
