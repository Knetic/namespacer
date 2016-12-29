namespacer
====

For all C# codefiles under a given path, ensures they have a specific namespace.

What are the caveats?
====

* Uses string manipulation to determine this, rather than parsing and AST. This can lead to inconsistent results if your comments look like class/enum/struct signatures.

* Only does one namespace per file. Seriously, why would you do anything else?

* Depends upon having an explicit access modifier on your type signatures. Because seriously, why would you _not_ do that?
