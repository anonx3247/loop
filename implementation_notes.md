# Notes

Currently I am moving to using regex for lexing.
Under the current system for skipping previous match groups, I am first regex matching,
and then eliminating those which intersect with &skip,
however it might be better (for performance) to later on change this so as to
instead separate the "program" and use `find_at` with `inc_skip` on a while loop
instead of the current `find_iter`.
