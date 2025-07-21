# imgz

A little tool for displaying images. It's terrible, but it fits my needs *perfectly*.

## Why?

I just needed a little something to display images for me to browse. Specifically, I was looking at sprites for a game project, so I just wanted something quick and dirty.

All the tools I found for this had 9000 features I didn't need, or dependencies on some random file browser, or whatever. I said to myself, "Self, this is trivial, just make one!"

## Trivial?! HAH!

Yeah, trivial. All the heavy lifting is done by [Ebitengine](https://ebitengine.org/), so all I had to do was to pick some horrible key bindings and smush it all together.

Since I was already poking at a simple game project at the time, ebitengine was already in my brain, so it was a natural fit for drawing images.

## What keys do what?

| Key        | Function                              |
|------------|---------------------------------------|
| F          | Toggle Fullscreen                     |
| Tab        | Switch auto mode                      |
| Ctrl+Up    | Decrease Y offset                     |
| Ctrl+Down  | Increase Y offset                     |
| Ctrl+Left  | Decrease X offset                     |
| Ctrl+Right | Increase X offset                     |
| PgUp       | Decrease scale                        |
| PgDown     | Increase scale                        |
| Right      | Next image                            |
| Left       | Previous image                        |
| Up         | (If in auto mode) Decrease auto delay |
| Down       | (If in auto mode) Increase auto delay |

## Auto mode?!

Yeah, automagically go to a different image after waiting a bit. The default is to not do that, because 99% of the time, I want to do that manually. Presing Tab puts it in auto mode, meaning it automagically progresses to the next picture. Prssing tab again puts it in the silly mde called auto mode random, where it waits a bit between each image, and then switches to a random one.

Shut up, it was a fun feature to implement, and I use it for ... uh ... LOOK! SQUIRREL!

## There is a file called "fuck.go"?!

Yeah, sorry, please don't tell my mom. It's just the dirtiest error checker on the block. Don't read too much into it.

## THIS LIT MY SOCK ON FIRE, AND MADE MY BOYFRIEND CRY!!!

Okay, I'm very sad you looked at the code. I apologize. I warned you, though! Please read [the license](LICENSE) for details on your lack of legal remidies.

## You closed my issue with the comment "lol ok"?!?

lol ok
