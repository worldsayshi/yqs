# Copilot Instructions

## Terminal output

Whenever you run a command in the terminal, pipe the output to a file, output.txt, that you can read from. You can use `tee <project-root>/output.txt`, so that I also may read the output. Make sure to overwrite each time so that it doesn't grow too big. There is a bug in the current version of Copilot that causes it to not read the output of commands correctly. This workaround allows you to read the output from the temporary file instead.

## Tasks

TODO's should be maintained in a file called `TODO.md` in the root of the repository. Each task should be a single line, and should be prefixed with `- [ ]` when it is not done, and `- [x]` when it is done. To keep keep the `TODO.md` clean, done tasks can be periodically moved to a file called `DONE.md` in the root of the repository. Try to keep the tasks few, to not bloat the TODO file, and small, so that they can be completed in a single commit. If you need to break a task into multiple commits, you may do so by adding an indented line below the task, similarly prefixed with `  - [ ]`.

## Requirements

Track requirements in a file called `REQUIREMENTS.md` in the root of the repository. Each requirement should be a single line, and should be prefixed with `- [ ]` when it is not done or malfunctioning, and `- [x]` when it is done. This file should be used to track the requirements of the project, and should be updated as the project evolves.

## Code style notes

- Always try to keep the code minimalistic.