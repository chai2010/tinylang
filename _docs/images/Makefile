# Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# https://github.com/asciinema/asciinema
# https://github.com/asciinema/asciicast2gif

# python3 -m pip install asciinema
#
# python3 -m asciinema rec file.cast
# exit
#
# docker run --rm -v $PWD:/data asciinema/asciicast2gif file.cast file.gif
# convert -delay 50 -loop 0 file.gif file-loop.gif

CAST_FILES=$(wildcard ./*.cast)
CAST_GIF_FILES=$(patsubst %.cast,%.cast.gif,$(CAST_FILES))

default: ${CAST_GIF_FILES}
	@echo "ok"

rec:
	-@rm _tmp.cast
	@python3 -m asciinema rec _tmp.cast
	
	# mv _tmp.cast xxx.cast

clean:

%.cast.gif: %.cast
	docker run --rm -v `pwd`:/data asciinema/asciicast2gif -s 5 $< _tmp.cast.gif
	convert -delay 100 -loop 0 _tmp.cast.gif $@
	-rm _tmp.cast.gif
