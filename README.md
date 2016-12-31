
Blink
=====

This package lets you make your keyboard lights blink. You can use blinking
keyboard lights for all kinds of fun things, such as [notifications of when
someone visits your
website!](http://lelandbatey.com/posts/2016/12/Making-lights-blink-for-each-HTTP-request/)

Install
-------

If you just run "go install" on this, you'll end up with a binary that's only
usable by running `sudo`, and that's no fun!

Instead, since this binary is quite harmless, we want anyone to be able to run
it as root! So, to do that we can do this:

    # Assuming you're in this directory, will put an executable named "blink"
    # in current directory
    go build -o blink ./cmd/blink
    # Set the owning user to be root
    sudo chown root blink
    # Set the sticky bit on the executable owned by root, so no matter who
    # launches it it's launched with root priveleges.
    sudo chmod u+s blink

And now that you've got your magical "always runs as root" executable, you can
place it wherever you want (such as in your `$PATH`).

