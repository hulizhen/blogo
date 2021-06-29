#!/bin/sh

# When directly invoking 'make', which is specified by the 'command' field in the 
# tasks.json file, VSCode can't normally kill the subprocess, e.g. sass, run by the
# task in Makefile after endding a debug session. Here we wrap it up to workaround it.
make $@